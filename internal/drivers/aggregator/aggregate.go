package aggregator

import (
	"kinoshkin/internal/usecase"
)

type Aggregator interface {
	Aggregate() error
}

func Cinemas(repo usecase.CinemasRepository) Aggregator {
	return cinemaAgg{repo}
}

func Cities(repo usecase.CitiesRepository) Aggregator {
	return cityAgg{repo}
}

func Schedule(
	movies usecase.MoviesRepository,
	cinemas usecase.CinemasRepository,
	schedules usecase.SchedulesRepository,
) Aggregator {
	return &scheduleAgg{
		movies:    movies,
		cinemas:   cinemas,
		schedules: schedules,
	}
}
