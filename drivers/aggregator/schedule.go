package aggregator

import (
	"kinoshkin/pkg/set"
	"kinoshkin/usecase"

	"github.com/schollz/progressbar/v3"
)

type scheduleAgg struct {
	movies    usecase.MoviesRepository
	cinemas   usecase.CinemasRepository
	schedules usecase.SchedulesRepository
}

func (s *scheduleAgg) Aggregate() error {
	cinemas, err := s.cinemas.GetAll("saint-petersburg")
	if err != nil {
		return err
	}

	kp := &kpAPI{seenMovies: set.New()}

	bar := progressbar.Default(int64(len(cinemas)),
		"Cinemas schedule aggregating...")
	for _, cinema := range cinemas {
		kp.aggregateCinemaData(cinema.ID)
		bar.Add(1)
	}
	_ = bar.Clear()

	movies, schedules := kp.result()

	kpUnoff := kpUnoffAPI{}

	bar = progressbar.Default(int64(len(movies)),
		"Extending movies data")
	for i, mov := range movies {
		kpUnoff.extendMovie(&mov)
		movies[i] = mov
		bar.Add(1)
	}
	_ = bar.Clear()

	if err := s.movies.Create(movies); err != nil {
		return err
	}
	return s.schedules.Create(schedules)
}
