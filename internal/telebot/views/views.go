package views

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"kinoshkin/domain"
)

func MoviesTable(movies []*domain.Movie) [][]tb.InlineButton {
	var table [][]tb.InlineButton
	for _, mov := range movies {
		table = append(table, []tb.InlineButton{
			{
				Text: fmt.Sprintf("%s %.2f", mov.Title, mov.Rating.KP),
				Data: fmt.Sprintf("mov%s", mov.ID),
			},
		})
	}

	return table
}

func CinemasTable(cinemas []*domain.Cinema, distances []int) [][]tb.InlineButton {
	var table [][]tb.InlineButton
	for i, cinema := range cinemas {
		table = append(table, []tb.InlineButton{
			{
				Text: fmt.Sprintf("%s (%dm)", cinema.Name, distances[i]),
				Data: fmt.Sprintf("cin%s", cinema.ID),
			},
		})
	}

	return table
}

func MovieCard(*domain.Movie) (string, [][]tb.InlineButton) {
	return "", nil
}

func CinemaCard(*domain.Cinema) (string, [][]tb.InlineButton) {
	return "", nil
}

func Commands() [][]tb.ReplyButton {
	return nil
}
