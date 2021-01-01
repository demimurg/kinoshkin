package aggregator

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type city struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

type cityAgg struct {
	db     *mongo.Database
	cities []city
}

func (cityAgg) Name() string {
	return "Cities Aggregator"
}

func (c cityAgg) Aggregate() error {
	resp, err := http.Get(cfg.CitiesURL + "?city=saint-petersburg")
	if err != nil {
		return errors.Wrap(err, "get api request")
	}

	var cities struct{ Data []city }
	err = json.NewDecoder(resp.Body).Decode(&cities)
	if err != nil {
		return errors.Wrap(err, "decoding response body")
	}

	lenCities := int64(len(cities.Data))
	bar := progressbar.Default(lenCities,
		"Cities aggregation...")

	// convert to db argument type
	citiesI := make([]interface{}, lenCities)
	for i, city := range cities.Data {
		bar.Add(1)
		citiesI[i] = city
	}
	_, err = c.db.Collection("cities").InsertMany(context.TODO(), citiesI)

	return err
}
