package aggregator

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type coords struct {
	Longitude float64 `bson:"longitude"`
	Latitude  float64 `bson:"latitude"`
}

type loc struct {
	Type        string `bson:"type"`
	Coordinates coords `bson:"coordinates"`
}

type cinema struct {
	ID       string   `bson:"_id"`
	Name     string   `bson:"name"`
	Address  string   `bson:"address"`
	City     string   `bson:"city"`
	Timezone string   `bson:"timezone"`
	Metros   []string `bson:"metros,omitempty"`
	Location loc      `bson:"location"`
}

type cinemaAgg struct {
	db *mongo.Database
}

func (c cinemaAgg) Aggregate() error {
	resp, err := http.Get(cfg.CinemasURL + "?city=saint-petersburg&limit=200")
	if err != nil {
		return errors.Wrap(err, "can't get cinemas external from api")
	}

	var cinemasRaw struct{ Items []map[string]interface{} }
	err = json.NewDecoder(resp.Body).Decode(&cinemasRaw)
	if err != nil {
		return errors.Wrap(err, "decode body err")
	}

	lenItems := int64(len(cinemasRaw.Items))
	bar := progressbar.Default(lenItems,
		"Cinemas aggregation...")

	cinemas := make([]interface{}, lenItems)
	for i, raw := range cinemasRaw.Items {
		bar.Add(1)
		city := raw["city"].(map[string]interface{})
		coordinates := raw["coordinates"].(map[string]interface{})

		cinemas[i] = cinema{
			ID:       raw["id"].(string),
			Name:     raw["title"].(string),
			Address:  raw["address"].(string),
			City:     city["id"].(string),
			Timezone: city["timezone"].(string),
			Metros:   extractMetros(raw["metro"]),
			Location: loc{
				Type: "Point",
				Coordinates: coords{
					Longitude: coordinates["longitude"].(float64),
					Latitude:  coordinates["latitude"].(float64),
				},
			},
		}
	}

	_, err = c.db.Collection("cinemas").InsertMany(context.TODO(), cinemas)

	return err
}

func extractMetros(m interface{}) []string {
	metrosI, ok := m.([]interface{})
	if !ok {
		return nil
	}

	var metros = make([]string, len(metrosI))
	for i, metroI := range metrosI {
		metro := metroI.(map[string]interface{})
		metros[i] = metro["name"].(string)
	}

	return metros
}
