package mongodb

import (
	"kinoshkin/entity"
	"time"
)

type showtime struct {
	ID    string    `bson:"_id"`
	Time  time.Time `bson:"time"`
	Price int       `bson:"price"`
}

type schedule struct {
	City         string     `bson:"city"`
	CinemaID     string     `bson:"cinema_id"`
	MovieID      string     `bson:"movie_id"`
	LastShowtime time.Time  `bson:"last"`
	Showtimes    []showtime `bson:"showtimes"`
}

type cinemaSchedule struct {
	Movie    `bson:",inline"`
	Schedule schedule `bson:"schedule"`
}

type movieSchedule struct {
	Cinema   `bson:",inline"`
	Schedule schedule `bson:"schedule"`
}

func toDomainSessions(sched schedule) []entity.Session {
	sessions := make([]entity.Session, len(sched.Showtimes))
	for i := range sched.Showtimes {
		sessions[i] = entity.Session{
			ID:    sched.Showtimes[i].ID,
			Start: sched.Showtimes[i].Time,
			Price: sched.Showtimes[i].Price,
		}
	}
	return sessions
}

func toMongoSchedule(sched entity.Schedule) schedule {
	var (
		last      time.Time
		showtimes = make([]showtime, len(sched.Sessions))
	)
	for i, ses := range sched.Sessions {
		if ses.Start.After(last) {
			last = ses.Start
		}
		showtimes[i] = showtime{
			ID:    ses.ID,
			Time:  ses.Start,
			Price: ses.Price,
		}
	}

	return schedule{
		City:         "saint-petersburg",
		CinemaID:     sched.CinemaID,
		MovieID:      sched.MovieID,
		LastShowtime: last,
		Showtimes:    showtimes,
	}
}
