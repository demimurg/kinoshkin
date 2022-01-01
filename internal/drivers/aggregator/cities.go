package aggregator

import (
	"encoding/json"
	"kinoshkin/internal/entity"
	"kinoshkin/internal/usecase"
	"net/http"

	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
)

type citiesJSON struct {
	Data []struct {
		ID   string
		Name string
	}
}

type cityAgg struct {
	repo usecase.CitiesRepository
}

func (c cityAgg) Aggregate() error {
	resp, err := http.Get(cfg.CitiesURL + "?raw=saint-petersburg")
	if err != nil {
		return errors.Wrap(err, "get api request")
	}

	var citiesRaw citiesJSON
	err = json.NewDecoder(resp.Body).Decode(&citiesRaw)
	if err != nil {
		return errors.Wrap(err, "decoding response body")
	}
	_ = resp.Body.Close()

	lenCities := int64(len(citiesRaw.Data))
	bar := progressbar.Default(lenCities,
		"Cities aggregation...")

	cities := make([]entity.City, lenCities)
	for i, raw := range citiesRaw.Data {
		bar.Add(1)
		cities[i] = entity.City{
			ID:   raw.ID,
			Name: raw.Name,
		}
	}

	return c.repo.Create(cities)
}
