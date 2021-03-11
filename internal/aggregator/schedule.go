package aggregator

import (
	"context"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"github.com/kr/pretty"
	"github.com/schollz/progressbar/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type schedule struct {
	City         string     `bson:"city"`
	CinemaId     string     `bson:"cinema_id"`
	MovieId      string     `bson:"movie_id"`
	LastShowtime time.Time  `bson:"last"`
	Showtimes    []showtime `bson:"showtimes"`
}

type showtime struct {
	ID    string    `bson:"_id"`
	Time  time.Time `bson:"time"`
	Price int       `bson:"price"`
}

type movie struct {
	ID             string              `bson:"_id"`
	KpId           string              `bson:"kp_id"`
	Title          string              `bson:"title"`
	TitleOriginal  string              `bson:"title_original,omitempty"`
	Rating         rating              `bson:"rating,omitempty"`
	AgeRestriction string              `bson:"age_restriction,omitempty"`
	Duration       int                 `bson:"duration,omitempty"`
	Description    string              `bson:"description"`
	Staff          map[string][]string `bson:"staff"`
	LandscapeImg   string              `bson:"landscape_img"`
	Trailer        trailer             `bson:"trailer,omitempty"`
	DateReleased   time.Time           `bson:"date_released,omitempty"`
}

type rating struct {
	KP   float64 `bson:"kp,omitempty"`
	IMDb float64 `bson:"imdb,omitempty"`
}

type trailer struct {
	Name string `bson:"name,omitempty"`
	Url  string `bson:"url,omitempty"`
}

type scheduleAgg struct {
	db          *mongo.Database
	schedules   []interface{}
	movies      map[string]*movie
	emptyMovies map[string]*movie
}

func (sa *scheduleAgg) Aggregate() error {
	cinemasId, err := sa.collectCinemasID()
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(cinemasId)),
		"Cinemas schedule aggregating...")
	for _, id := range cinemasId {
		bar.Add(1)
		err := sa.aggregateSchedule(id)
		if err != nil {
			log.Printf(
				"Cinema (%s): fetching schedule error - %s\n",
				id, err,
			)
		}
	}
	_ = bar.Clear()

	// todo: handle empty movies, don't discard them
	if len(sa.emptyMovies) != 0 {
		pretty.Printf("Have %d half-empty movies:\n%# v", len(sa.emptyMovies), sa.emptyMovies)
	}
	if err := sa.extendMovies(); err != nil {
		log.Println("Error occurred while extending movies: ", err)
	}

	movies := make([]interface{}, 0, len(sa.movies))
	for _, mov := range sa.movies {
		movies = append(movies, *mov)
	}

	ctx := context.TODO()
	opts := &options.UpdateOptions{}
	opts = opts.SetUpsert(true)

	for _, mov := range sa.movies {
		_, err = sa.db.Collection("movies").UpdateOne(ctx, bson.M{"_id": mov.ID}, bson.M{
			"$set": interface{}(*mov),
		}, opts)
		if err != nil {
			return err
		}
	}
	_, err = sa.db.Collection("schedule").InsertMany(ctx, sa.schedules)

	return err
}

func (sa *scheduleAgg) collectCinemasID() ([]string, error) {
	cursor, err := sa.db.Collection("cinemas").Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var ids []string
	for cursor.Next(context.TODO()) {
		var cinema = make(map[string]interface{})
		err := cursor.Decode(cinema)
		if err != nil {
			return nil, err
		}
		ids = append(ids, cinema["_id"].(string))
	}

	return ids, nil
}

func (sa *scheduleAgg) aggregateSchedule(cinemaId string) error {
	bytes, err := getScheduleJSON(cinemaId)
	if err != nil {
		return err
	}
	page, err := gabs.ParseJSON(bytes)
	if err != nil {
		return err
	}

	for _, item := range page.S("schedule", "items").Children() {
		event := item.S("event")
		movie := movie{}

		movie.ID, _ = event.S("id").Data().(string)
		movie.Title, _ = event.S("title").Data().(string)
		movie.TitleOriginal, _ = event.S("originalTitle").Data().(string)
		movie.LandscapeImg, _ = event.S("image", "eventCoverL2x", "url").Data().(string)
		movie.Description, _ = event.S("argument").Data().(string)
		movie.AgeRestriction, _ = event.S("contentRating").Data().(string)

		kpURL, ok := event.S("kinopoisk", "url").Data().(string)
		if !ok {
			kpURL, ok = event.S("image", "source", "url").Data().(string)
			if !ok || !strings.Contains(kpURL, "kinopoisk") {
				sa.emptyMovies[movie.ID] = &movie
				continue
			}
		}
		url := strings.Split(kpURL, "/")
		movie.KpId = url[len(url)-1]

		dateRaw, _ := event.S("dateReleased").Data().(string)
		movie.DateReleased, _ = time.Parse("2006-01-02", dateRaw)
		movie.Rating.KP, _ = event.
			S("kinopoisk", "value").
			Data().(float64)

		sa.movies[movie.KpId] = &movie

		// ignore different formats (like "2D", "3D") for sessions
		var jsonSessions []*gabs.Container
		for _, format := range item.S("schedule").Children() {
			jsonSessions = append(jsonSessions, format.S("sessions").Children()...)
		}

		var showtimes []showtime
		for _, session := range jsonSessions {
			ticketID, _ := session.S("ticket", "id").Data().(string)
			startAt, err := time.Parse(
				"2006-01-02T15:04:05",
				strings.Trim(session.S("datetime").Data().(string), "\""),
			)
			if err != nil {
				log.Printf(
					"showtime time (%q) parsing error\n",
					session.S("datetime").Data().(string),
				)
				continue
			}
			price, ok := session.
				S("ticket", "price", "min").
				Data().(float64)
			if !ok {
				// todo: log no price tickets
				continue
			}

			showtimes = append(showtimes, showtime{
				ID:    ticketID,
				Time:  startAt,
				Price: int(price / 100),
			})
		}

		sort.Slice(showtimes, func(i, j int) bool {
			return showtimes[i].Time.Before(
				showtimes[j].Time,
			)
		})

		if len(showtimes) > 0 {
			sa.schedules = append(sa.schedules, schedule{
				City:         "saint-petersburg",
				CinemaId:     cinemaId,
				MovieId:      movie.ID,
				LastShowtime: showtimes[len(showtimes)-1].Time,
				Showtimes:    showtimes,
			})
		}
	}

	return nil
}

func getScheduleJSON(cinemaId string) ([]byte, error) {
	// todo: remove add 24 * time.Hour (debug only)
	resp, err := http.Get(fmt.Sprintf(
		cfg.ScheduleURL+"/%s/schedule_cinema?date=%s&city=saint-petersburg&limit=200",
		cinemaId, time.Now().Format("2006-01-02"),
	))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (sa *scheduleAgg) extendMovies() error {
	bar := progressbar.Default(int64(len(sa.movies)),
		"Extending movies data")

	for id, movie := range sa.movies {
		movData, err := getFromKpApi(movieDataUri, id)
		if err != nil {
			return err
		}
		trailers, err := getFromKpApi(trailersUri, id)
		if err != nil {
			return err
		}
		staff, err := getFromKpApi(staffUri, id)
		if err != nil {
			return err
		}

		movie.Staff = make(map[string][]string)
		for _, employee := range staff.Children() {
			empKey, _ := employee.S("professionKey").Data().(string)
			switch empKey {
			case "DIRECTOR", "WRITER", "OPERATOR", "COMPOSITOR":
			case "ACTOR":
				if len(movie.Staff["actor"]) > 6 {
					continue
				}
			default:
				continue
			}

			key := strings.ToLower(empKey)
			empName, _ := employee.S("nameRu").Data().(string)
			movie.Staff[key] = append(movie.Staff[key], empName)
		}

		duration, ok := movData.S("data", "filmLength").Data().(string)
		if ok {
			movie.Duration = convertToMinutes(duration)
		}

		movie.Rating.IMDb, _ = movData.S("rating", "ratingImdb").Data().(float64)
		movie.Trailer.Url, _ = trailers.S("trailers", "0", "url").Data().(string)
		movie.Trailer.Name, _ = trailers.S("trailers", "0", "name").Data().(string)
		bar.Add(1)
	}
	_ = bar.Clear()

	return nil
}

func convertToMinutes(dur string) int {
	t := strings.Split(dur, ":")
	if len(t) != 2 {
		return 0
	}

	h, _ := strconv.Atoi(t[0])
	m, _ := strconv.Atoi(t[1])
	return 60*h + m
}

const (
	movieDataUri = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s?append_to_response=RATING"
	trailersUri  = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s/videos"
	staffUri     = "https://kinopoiskapiunofficial.tech/api/v1/staff?filmId=%s"
)

func getFromKpApi(uri, id string) (*gabs.Container, error) {
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(uri, id),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", cfg.TokenKP)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return gabs.ParseJSONBuffer(resp.Body)
}
