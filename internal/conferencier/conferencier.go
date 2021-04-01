package conferencier

import "kinoshkin/domain"

type conf struct {
	cinemas   domain.CinemasRepository
	movies    domain.MoviesRepository
	users     domain.UsersRepository
	schedules domain.SchedulesRepository
}

func (c *conf) FindMovies(userID int, pag domain.P) ([]domain.Movie, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.movies.FindByRating(user.City, pag)
}

func (c *conf) FindCinemas(userID int, pag domain.P) ([]domain.Cinema, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.cinemas.FindNearby(user, pag)
}

func (c *conf) GetMovie(movieID string) (*domain.Movie, error) {
	return c.movies.Get(movieID)
}

func (c *conf) GetMovieSchedule(userID int, movieID string, pag domain.P) ([]domain.CinemaWithSessions, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.schedules.GetForMovie(movieID, user, pag)
}

func (c *conf) GetCinema(cinemaID string) (*domain.Cinema, error) {
	return c.cinemas.Get(cinemaID)
}

func (c *conf) GetCinemaSchedule(cinemaID string, pag domain.P) ([]domain.MovieWithSessions, error) {
	return c.schedules.GetForCinema(cinemaID, pag)
}

func (c *conf) UpdateUserLocation(userID int, lat, long float32) error {
	return c.users.UpdateLoc(userID, lat, long)
}

func (c *conf) RegisterUser(userID int, name string) error {
	return c.users.Create(userID, name)
}

func New(
	cinemasRepo domain.CinemasRepository,
	moviesRepo domain.MoviesRepository,
	usersRepo domain.UsersRepository,
	schedulesRepo domain.SchedulesRepository,
) *conf {
	return &conf{
		cinemas:   cinemasRepo,
		movies:    moviesRepo,
		users:     usersRepo,
		schedules: schedulesRepo,
	}
}
