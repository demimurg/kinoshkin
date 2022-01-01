package usecase

import "kinoshkin/internal/entity"

type UsersRepository interface {
	Get(id int) (*entity.User, error)
	// UpdateLoc used for updating user geo location
	UpdateLoc(id int, lat, long float32) error
	Create(id int, name string) error
}

type SchedulesRepository interface {
	Create(schedules []entity.Schedule) error
	GetForMovie(movieID string, user *entity.User, pag P) ([]entity.CinemaWithSessions, error)
	GetForCinema(cinemaID string, pag P) ([]entity.MovieWithSessions, error)
}

// MoviesRepository work with movies collection
type MoviesRepository interface {
	Create(movs []entity.Movie) error
	Get(movID string) (*entity.Movie, error)
	FindByRating(city string, pag P) ([]entity.Movie, error)
}

// CinemasRepository provides basic db methods for cinemas collection
type CinemasRepository interface {
	Create(cinemas []entity.Cinema) error
	Get(cinemaID string) (*entity.Cinema, error)
	GetAll(city string) ([]entity.Cinema, error)
	// FindNearby search cinemas near user location
	FindNearby(user *entity.User, pag P) ([]entity.Cinema, error)
}

type CitiesRepository interface {
	Create(cities []entity.City) error
}
