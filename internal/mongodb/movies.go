package mongodb

import (
	"kinoshkin/domain"
	"time"
)

type Movie struct {
	ID             string              `bson:"_id"`
	KpId           string              `bson:"kp_id"`
	Title          string              `bson:"title"`
	TitleOriginal  string              `bson:"title_original,omitempty"`
	Rating         rating              `bson:"rating,omitempty"`
	AgeRestriction string              `bson:"age_restriction,omitempty"`
	Duration       int                 `bson:"duration,omitempty"`
	Description    string              `bson:"description"`
	Staff          map[string][]string `bson:"staff"`
	LandscapeImg   string              `bson:"landscape_img"`
	Trailer        trailer             `bson:"trailer,omitempty"`
	DateReleased   time.Time           `bson:"date_released,omitempty"`
}

type rating struct {
	KP   float64 `bson:"kp,omitempty"`
	IMDb float64 `bson:"imdb,omitempty"`
}

type trailer struct {
	Name string `bson:"name,omitempty"`
	Url  string `bson:"url,omitempty"`
}

type schedule struct {
	City         string     `bson:"city"`
	CinemaId     string     `bson:"cinema_id"`
	MovieId      string     `bson:"movie_id"`
	LastShowtime time.Time  `bson:"last"`
	Showtimes    []showtime `bson:"showtimes"`
}

type showtime struct {
	ID    string    `bson:"_id"`
	Time  time.Time `bson:"time"`
	Price int       `bson:"price"`
}

type movieSchedule struct {
	Cinema   `bson:",inline"`
	Schedule schedule `bson:"schedule"`
}

func toDomainMovie(mov *Movie) domain.Movie {
	return domain.Movie{
		ID:             mov.ID,
		Title:          mov.Title,
		Description:    mov.Description,
		PosterURL:      mov.LandscapeImg,
		Duration:       int32(mov.Duration),
		AgeRestriction: mov.AgeRestriction,
		FilmCrew:       mov.Staff,
		Rating: domain.Rating{
			IMDB: mov.Rating.IMDb,
			KP:   mov.Rating.KP,
		},
	}
}

func toMongoMovie(mov *domain.Movie) *Movie {
	return &Movie{
		ID:             mov.ID,
		Title:          mov.Title,
		Staff:          mov.FilmCrew,
		Duration:       int(mov.Duration),
		Description:    mov.Description,
		LandscapeImg:   mov.PosterURL,
		AgeRestriction: mov.AgeRestriction,
		Rating: rating{
			IMDb: mov.Rating.IMDB,
			KP:   mov.Rating.KP,
		},
	}
}
