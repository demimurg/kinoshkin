package aggregator

import (
	"encoding/json"
	"fmt"
	"kinoshkin/entity"
	"kinoshkin/pkg/set"
	"log"
	"net/http"
	"strings"
	"time"
)

type scheduleJSON struct {
	Schedule struct {
		Items []struct {
			Event struct {
				ID            string
				Title         string
				OriginalTitle string
				ContentRating string
				DateReleased  string
				Argument      string
				Kinopoisk     struct {
					URL   string
					Value float64
				}
				Image struct {
					Source struct {
						URL string
					}
					EventCoverL2x struct {
						URL string
					}
				}
			}
			Schedule []struct {
				Sessions []struct {
					Datetime string
					Ticket   struct {
						ID    string
						Price struct {
							Min float64
						}
					}
				}
			}
		}
	}
}

type kpAPI struct {
	movies     []entity.Movie
	schedules  []entity.Schedule
	seenMovies set.Strings
}

func (api *kpAPI) get(cinemaID string, dest interface{}) error {
	resp, err := http.Get(fmt.Sprintf(
		cfg.ScheduleURL+"/%s/schedule_cinema?date=%s&city=saint-petersburg&limit=200",
		cinemaID, time.Now().Format("2006-01-02"),
	))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dest)
}

func (api *kpAPI) aggregateCinemaData(cinemaID string) {
	var data scheduleJSON
	if err := api.get(cinemaID, &data); err != nil {
		log.Printf("Cinema (%s): fetching schedule error - %s\n", cinemaID, err)
		return
	}

	for _, item := range data.Schedule.Items {
		var sessions []entity.Session
		for _, format := range item.Schedule {
			for _, ses := range format.Sessions {
				start, _ := time.Parse("2006-01-02T15:04:05", ses.Datetime)
				sessions = append(sessions, entity.Session{
					ID:    ses.Ticket.ID,
					Start: start,
					Price: int(ses.Ticket.Price.Min / 100),
				})
			}
		}

		api.schedules = append(api.schedules, entity.Schedule{
			MovieID:  item.Event.ID,
			CinemaID: cinemaID,
			Sessions: sessions,
		})

		if api.seenMovies.Have(item.Event.ID) {
			continue
		}
		api.seenMovies.Add(item.Event.ID)

		extractID := func(url string) string {
			url = strings.TrimRight(url, "/")
			if len(url) == 0 {
				return ""
			}
			return url[strings.LastIndex(url, "/")+1:]
		}
		kpID := extractID(item.Event.Kinopoisk.URL)
		if kpID == "" {
			kpID = extractID(item.Event.Image.Source.URL)
		}

		dateReleased, _ := time.Parse("2006-01-02", item.Event.DateReleased)

		api.movies = append(api.movies, entity.Movie{
			KpID:           kpID,
			DateReleased:   dateReleased,
			ID:             item.Event.ID,
			Title:          item.Event.Title,
			Description:    item.Event.Argument,
			PosterURL:      item.Event.Image.EventCoverL2x.URL,
			AgeRestriction: item.Event.ContentRating,
			Rating:         entity.Rating{KP: item.Event.Kinopoisk.Value},
		})
	}
}

func (api *kpAPI) result() ([]entity.Movie, []entity.Schedule) {
	return api.movies, api.schedules
}
