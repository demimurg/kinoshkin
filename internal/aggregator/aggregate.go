package aggregator

import (
	"kinoshkin/domain"
)

type Aggregator interface {
	Aggregate() error
}

func Cinemas(repo domain.CinemasRepository) Aggregator {
	return cinemaAgg{repo}
}

func Cities(repo domain.CitiesRepository) Aggregator {
	return cityAgg{repo}
}

func Schedule(
	movies domain.MoviesRepository,
	cinemas domain.CinemasRepository,
	schedules domain.SchedulesRepository,
) Aggregator {
	return &scheduleAgg{
		movies:    movies,
		cinemas:   cinemas,
		schedules: schedules,
	}
}
