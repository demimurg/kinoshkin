package mongodb

import (
	"kinoshkin/domain"
	"time"
)

type showtime struct {
	ID    string    `bson:"_id"`
	Time  time.Time `bson:"time"`
	Price int       `bson:"price"`
}

type schedule struct {
	City         string     `bson:"city"`
	CinemaId     string     `bson:"cinema_id"`
	MovieId      string     `bson:"movie_id"`
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

func toDomainSessions(sched schedule) []domain.Session {
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

func toMongoSchedule(sched domain.Schedule) schedule {
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
			Price: int(ses.Price),
		}
	}

	return schedule{
		City:         "saint-petersburg",
		CinemaId:     sched.CinemaID,
		MovieId:      sched.MovieID,
		LastShowtime: last,
		Showtimes:    showtimes,
	}
}
