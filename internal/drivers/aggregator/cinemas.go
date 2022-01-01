package aggregator

import (
	"encoding/json"
	"kinoshkin/internal/entity"
	"kinoshkin/internal/usecase"
	"net/http"

	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
)

type cinemasJSON struct {
	Items []struct {
		ID      string
		Title   string
		Address string
		City    struct {
			ID       string
			Timezone string
		}
		Coordinates struct {
			Longitude float64
			Latitude  float64
		}
		Metro []struct {
			Name string
		}
	}
}

type cinemaAgg struct {
	repo usecase.CinemasRepository
}

func (c cinemaAgg) Aggregate() error {
	resp, err := http.Get(cfg.CinemasURL + "?city=saint-petersburg&limit=200")
	if err != nil {
		return errors.Wrap(err, "can't get cinemas external from api")
	}

	var cinemasRaw cinemasJSON
	err = json.NewDecoder(resp.Body).Decode(&cinemasRaw)
	if err != nil {
		return errors.Wrap(err, "decode body err")
	}
	_ = resp.Body.Close()

	lenItems := int64(len(cinemasRaw.Items))
	bar := progressbar.Default(lenItems,
		"Cinemas aggregation...")

	cinemas := make([]entity.Cinema, lenItems)
	for i, raw := range cinemasRaw.Items {
		bar.Add(1)
		metros := make([]string, len(raw.Metro))
		for i, metro := range raw.Metro {
			metros[i] = metro.Name
		}

		cinemas[i] = entity.Cinema{
			ID:      raw.ID,
			Name:    raw.Title,
			Address: raw.Address,
			Metro:   metros,
			Lat:     raw.Coordinates.Latitude,
			Long:    raw.Coordinates.Longitude,
		}
	}

	return c.repo.Create(cinemas)
}
