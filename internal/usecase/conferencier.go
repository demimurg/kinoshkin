package usecase

import "kinoshkin/internal/entity"

// P is a pagination info
type P struct {
	Offset, Limit int
}

// Conferencier is a main service, it manage all data about user, movie and cinema
type Conferencier interface {
	FindMovies(userID int, pag P) ([]entity.Movie, error)
	FindCinemas(userID int, pag P) ([]entity.Cinema, error)

	GetMovie(movieID string) (*entity.Movie, error)
	GetMovieSchedule(userID int, movieID string, pag P) ([]entity.CinemaWithSessions, error)
	GetCinema(cinemaID string) (*entity.Cinema, error)
	GetCinemaSchedule(cinemaID string, pag P) ([]entity.MovieWithSessions, error)

	UpdateUserLocation(userID int, lat, long float32) error
	RegisterUser(userID int, name string) error
}

func NewConferencier(
	cinemasRepo CinemasRepository,
	moviesRepo MoviesRepository,
	usersRepo UsersRepository,
	schedulesRepo SchedulesRepository,
) Conferencier {
	return &conf{
		cinemas:   cinemasRepo,
		movies:    moviesRepo,
		users:     usersRepo,
		schedules: schedulesRepo,
	}
}

type conf struct {
	cinemas   CinemasRepository
	movies    MoviesRepository
	users     UsersRepository
	schedules SchedulesRepository
}

func (c *conf) FindMovies(userID int, pag P) ([]entity.Movie, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.movies.FindByRating(user.City, pag)
}

func (c *conf) FindCinemas(userID int, pag P) ([]entity.Cinema, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.cinemas.FindNearby(user, pag)
}

func (c *conf) GetMovie(movieID string) (*entity.Movie, error) {
	return c.movies.Get(movieID)
}

func (c *conf) GetMovieSchedule(userID int, movieID string, pag P) ([]entity.CinemaWithSessions, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.schedules.GetForMovie(movieID, user, pag)
}

func (c *conf) GetCinema(cinemaID string) (*entity.Cinema, error) {
	return c.cinemas.Get(cinemaID)
}

func (c *conf) GetCinemaSchedule(cinemaID string, pag P) ([]entity.MovieWithSessions, error) {
	return c.schedules.GetForCinema(cinemaID, pag)
}

func (c *conf) UpdateUserLocation(userID int, lat, long float32) error {
	return c.users.UpdateLoc(userID, lat, long)
}

func (c *conf) RegisterUser(userID int, name string) error {
	return c.users.Create(userID, name)
}
