package aggregator

import (
	"encoding/json"
	"fmt"
	"kinoshkin/domain"
	"net/http"
	"strconv"
	"strings"
)

const (
	movieDataURL = "https://kinopoiskapiunofficial.tech/api/v2.1/films/%s?append_to_response=RATING"
	staffURL     = "https://kinopoiskapiunofficial.tech/api/v1/staff?filmId=%s"
)

type unoffMovieJSON struct {
	Data struct {
		FilmLength string
	}
	Rating struct {
		Rating     float64
		RatingImdb float64
	}
}

type unoffStaffJSON []struct {
	NameRu        string
	ProfessionKey string
}

type kpUnoffAPI struct {
}

func (api *kpUnoffAPI) get(resourceUrl, id string, dest interface{}) error {
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(resourceUrl, id),
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-KEY", cfg.TokenKP)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(dest)
}

func (api *kpUnoffAPI) extendMovie(mov *domain.Movie) {
	var (
		movieJSON unoffMovieJSON
		staffJSON unoffStaffJSON
	)
	_ = api.get(movieDataURL, mov.KpID, &movieJSON)
	_ = api.get(staffURL, mov.KpID, &staffJSON)

	mov.FilmCrew = make(map[domain.Position][]string)
	var pos domain.Position
	for _, employee := range staffJSON {
		switch employee.ProfessionKey {
		case "DIRECTOR":
			pos = domain.Director
		case "WRITER":
			pos = domain.Screenwriter
		case "OPERATOR":
			pos = domain.Operator
		case "COMPOSITOR":
			pos = domain.Operator
		case "ACTOR":
			if len(mov.FilmCrew["actor"]) > 6 {
				continue
			}
			pos = domain.Actor
		default:
			continue
		}

		if employee.NameRu != "" {
			mov.FilmCrew[pos] = append(
				mov.FilmCrew[pos],
				employee.NameRu,
			)
		}
	}

	mov.Duration = convertToMinutes(movieJSON.Data.FilmLength)
	mov.Rating.IMDB = movieJSON.Rating.RatingImdb
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
