package views

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"kinoshkin/domain"
	"strings"
)

var CinemasCmd = tb.ReplyButton{Text: "Кинотеатры🍿"}
var MoviesCmd = tb.ReplyButton{Text: "Фильмы🎬"}

func MoviesList(movies []*domain.Movie) [][]tb.InlineButton {
	var table [][]tb.InlineButton
	for _, mov := range movies {
		table = append(table, []tb.InlineButton{
			{
				Text: fmt.Sprintf("%s (%.2f)", mov.Title, mov.Rating.KP),
				Data: Encode(MoviePrefix, mov.ID),
			},
		})
	}

	return table
}

func CinemasList(cinemas []*domain.Cinema, distances []int) [][]tb.InlineButton {
	var table [][]tb.InlineButton
	for i, cinema := range cinemas {
		table = append(table, []tb.InlineButton{
			{
				Text: fmt.Sprintf("%s (%dm)", cinema.Name, distances[i]),
				Data: Encode(CinemaPrefix, cinema.ID),
			},
		})
	}

	return table
}

func MovieCard(mov *domain.Movie) (msg interface{}, opts []interface{}) {
	title := fmt.Sprintf("*%s* _(imdb: %.1f, kp: %.1f)_", mov.Title, mov.Rating.IMDB, mov.Rating.KP)
	// todo: generalize
	title = strings.ReplaceAll(title, ".", "\\.")
	title = strings.ReplaceAll(title, "!", "\\!")
	title = strings.ReplaceAll(title, "(", "\\(")
	title = strings.ReplaceAll(title, ")", "\\)")

	duration := fmt.Sprintf("Продолжительность: `%d мин`", mov.Duration)
	// todo: remove duplicates
	creators := "Создатели: `" + strings.Join(append(
		mov.FilmCrew[domain.Director],
		mov.FilmCrew[domain.Screenwriter]...,
	), ", ") + "`"
	actors := "Актеры: `" + strings.Join(mov.FilmCrew[domain.Actor], ", ") + "`"

	return &tb.Photo{
		File: tb.File{FileURL: mov.PosterURL},
		Caption: strings.Join(
			[]string{title, duration, creators, actors},
			"\n",
		),
	}, []interface{}{tb.ModeMarkdownV2}
}

func CinemaCard(cinema *domain.Cinema) (msg interface{}, opts []interface{}) {
	address := cinema.Address
	if len(cinema.Metro) > 0 {
		address = "🚇" + strings.Join(cinema.Metro, ", ") + "\n" + address
	}

	return &tb.Venue{
		Location:     tb.Location{Lat: cinema.Lat, Lng: cinema.Long},
		Title:        cinema.Name,
		Address:      address,
		FoursquareID: "4bf58dd8d48988d17f941735",
	}, nil
}
