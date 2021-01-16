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
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ticket struct {
	ID       string    `bson:"_id"`
	MovieId  string    `bson:"movie_id"`
	CinemaId string    `bson:"cinema_id"`
	City     string    `bson:"city"`
	Time     time.Time `bson:"time"`
	Price    int       `bson:"price"`
}

type movie struct {
	ID             string              `bson:"_id"`
	KpId           string              `bson:"kp_id"`
	Title          string              `bson:"title"`
	TitleOriginal  string              `bson:"title_original,omitempty"`
	Staff          map[string][]string `bson:"staff"`
	Duration       int                 `bson:"duration,omitempty"`
	DateReleased   time.Time           `bson:"date_released,omitempty"`
	LandscapeImg   string              `bson:"landscape_img"`
	Description    string              `bson:"description"`
	RatingKP       float64             `bson:"rating_kp,omitempty"`
	RatingIMDB     float64             `bson:"rating_imdb,omitempty"`
	AgeRestriction string              `bson:"age_restriction,omitempty"`
	Trailer        string              `bson:"trailer,omitempty"`
	TrailerName    string              `bson:"trailer_name,omitempty"`
}

type scheduleAgg struct {
	db          *mongo.Database
	tickets     []ticket
	movies      map[string]*movie
	emptyMovies map[string]*movie
}

func (sa *scheduleAgg) Aggregate() error {
	cinemasId, err := sa.collectCinemasID()
	if err != nil {
		return err
	}

	bar := progressbar.Default(int64(len(cinemasId)),
		"Schedule downloading...")

	for _, id := range cinemasId {
		bar.Add(1)
		err := sa.aggregateMoviesAndTickets(id)
		if err != nil {
			log.Printf(
				"Cinema (%s): fetching schedule error - %s\n",
				id, err,
			)
		}
	}
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
	tickets := make([]interface{}, 0, len(sa.tickets))
	for _, t := range sa.tickets {
		tickets = append(tickets, t)
	}

	ctx := context.TODO()
	opts := &options.InsertManyOptions{}
	opts = opts.SetOrdered(false)

	_, err = sa.db.Collection("movies").InsertMany(ctx, movies, opts)
	if err != nil {
		return err
	}
	_, err = sa.db.Collection("tickets").InsertMany(ctx, tickets)
	if err != nil {
		return err
	}

	return nil
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

func (sa *scheduleAgg) aggregateMoviesAndTickets(cinemaId string) error {
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
		movie.RatingKP, _ = event.
			S("kinopoisk", "value").
			Data().(float64)

		sa.movies[movie.KpId] = &movie

		// ignore different formats for sessions
		var sessions []*gabs.Container
		for _, format := range item.S("schedule").Children() {
			sessions = append(sessions, format.S("sessions").Children()...)
		}

		for _, session := range sessions {
			t := ticket{
				MovieId:  movie.ID,
				CinemaId: cinemaId,
				City:     "saint-petersburg",
			}

			t.ID, _ = session.S("ticket", "id").Data().(string)
			startAt, err := time.Parse(
				"2006-01-02T15:04:05",
				strings.Trim(session.S("datetime").Data().(string), "\""),
			)
			if err != nil {
				log.Printf(
					"ticket time (%q) parsing error\n",
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

			t.Time = startAt
			t.Price = int(price / 100)

			sa.tickets = append(sa.tickets, t)
		}
	}

	return nil
}

func getScheduleJSON(cinemaId string) ([]byte, error) {
	// todo: remove add 24 * time.Hour (debug only)
	resp, err := http.Get(fmt.Sprintf(
		cfg.ScheduleURL+"/%s/schedule_cinema?date=%s&city=saint-petersburg&limit=200",
		cinemaId, time.Now().Add(24*time.Hour).Format("2006-01-02"),
	))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (sa *scheduleAgg) extendMovies() error {
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
			movie.Duration = convertToMins(duration)
		}

		movie.RatingIMDB, _ = movData.S("rating", "ratingImdb").Data().(float64)
		movie.Trailer, _ = trailers.S("trailers", "0", "url").Data().(string)
		movie.TrailerName, _ = trailers.S("trailers", "0", "name").Data().(string)
	}
	return nil
}

func convertToMins(dur string) int {
	t := strings.Split(dur, ":")
	if len(t) != 2 {
		return 100
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
