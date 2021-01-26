package conferencier

import "kinoshkin/domain"

type conf struct {
	cinemas domain.CinemasRepo
	movies  domain.MoviesRepo
	users   domain.UsersRepo
}

func (c conf) FindMovies(userID int, pag domain.P) ([]*domain.Movie, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.movies.FindByRating(user.City, pag)
}

func (c conf) FindCinemas(userID int, pag domain.P) ([]*domain.Cinema, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.cinemas.FindNearby(user.Lat, user.Long, pag)
}

func (c conf) GetMovie(movieID string) (*domain.Movie, error) {
	return c.movies.Get(movieID)
}

func (c conf) GetMovieSchedule(userID int, movieID string, pag domain.P) ([]domain.CinemaWithSessions, error) {
	user, err := c.users.Get(userID)
	if err != nil {
		return nil, err
	}

	return c.movies.GetSchedule(movieID, user.City, pag)
}

func (c conf) GetCinema(cinemaID string) (*domain.Cinema, error) {
	return c.cinemas.Get(cinemaID)
}

func (c conf) GetCinemaSchedule(cinemaID string, pag domain.P) ([]domain.MovieWithSessions, error) {
	return c.cinemas.GetSchedule(cinemaID, pag)
}

func (c conf) UpdateUserLocation(userID int, lat, long float32) error {
	return c.users.UpdateLoc(userID, lat, long)
}

func (c conf) RegisterUser(userID int, name string) error {
	return c.users.Create(userID, name)
}
