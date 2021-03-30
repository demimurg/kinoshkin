package mongodb

import (
	"kinoshkin/domain"
)

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

type cinemaSchedule struct {
	Movie    `bson:",inline"`
	Schedule schedule `bson:"schedule"`
}

func toMongoCinema(cin *domain.Cinema) Cinema {
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

func toDomainCinema(cin *Cinema) domain.Cinema {
	return domain.Cinema{
		ID:      cin.ID,
		Name:    cin.Name,
		Address: cin.Address,
		Metro:   cin.Metros,
		Long:    cin.Location.Coordinates.Longitude,
		Lat:     cin.Location.Coordinates.Latitude,
	}
}

func toSessions(sched schedule) []domain.Session {
	sessions := make([]domain.Session, len(sched.Showtimes))
	for i := range sched.Showtimes {
		sessions[i] = domain.Session{
			ID:    sched.Showtimes[i].ID,
			Start: sched.Showtimes[i].Time,
			Price: int32(sched.Showtimes[i].Price),
		}
	}
	return sessions
}
