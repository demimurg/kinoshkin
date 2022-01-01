package mongodb

import "kinoshkin/entity"

type Cinema struct {
	ID       string   `bson:"_id"`
	Name     string   `bson:"name"`
	Address  string   `bson:"address"`
	City     string   `bson:"city"`
	Timezone string   `bson:"timezone"`
	Distance float64  `bson:"distance"`
	Metros   []string `bson:"metros,omitempty"`
	Location loc      `bson:"location"`
}

type coords struct {
	Longitude float64 `bson:"longitude"`
	Latitude  float64 `bson:"latitude"`
}

type loc struct {
	Type        string `bson:"type"`
	Coordinates coords `bson:"coordinates"`
}

func toMongoCinema(cin *entity.Cinema) Cinema {
	return Cinema{
		ID:       cin.ID,
		Name:     cin.Name,
		Address:  cin.Address,
		City:     "saint-petersburg",
		Timezone: "Europe/Moscow",
		Metros:   cin.Metro,
		Location: loc{
			Type: "Point",
			Coordinates: coords{
				Longitude: cin.Long,
				Latitude:  cin.Lat,
			},
		},
	}
}

func toDomainCinema(cin *Cinema) entity.Cinema {
	return entity.Cinema{
		ID:      cin.ID,
		Name:    cin.Name,
		Address: cin.Address,
		Metro:   cin.Metros,
		Long:    cin.Location.Coordinates.Longitude,
		Lat:     cin.Location.Coordinates.Latitude,
	}
}
