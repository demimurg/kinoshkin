package mongodb

import (
	"kinoshkin/domain"
	"time"
)

type Movie struct {
	ID             string              `bson:"_id"`
	KpID           string              `bson:"kp_id"`
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
	URL  string `bson:"url,omitempty"`
}

func toDomainMovie(mov *Movie) domain.Movie {
	return domain.Movie{
		ID:             mov.ID,
		KpID:           mov.KpID,
		Title:          mov.Title,
		Description:    mov.Description,
		PosterURL:      mov.LandscapeImg,
		Duration:       mov.Duration,
		AgeRestriction: mov.AgeRestriction,
		FilmCrew:       mov.Staff,
		DateReleased:   mov.DateReleased,
		Rating: domain.Rating{
			IMDB: mov.Rating.IMDb,
			KP:   mov.Rating.KP,
		},
	}
}

func toMongoMovie(mov *domain.Movie) *Movie {
	return &Movie{
		ID:             mov.ID,
		KpID:           mov.KpID,
		Title:          mov.Title,
		Staff:          mov.FilmCrew,
		Duration:       mov.Duration,
		Description:    mov.Description,
		LandscapeImg:   mov.PosterURL,
		AgeRestriction: mov.AgeRestriction,
		DateReleased:   mov.DateReleased,
		Rating: rating{
			IMDb: mov.Rating.IMDB,
			KP:   mov.Rating.KP,
		},
	}
}
