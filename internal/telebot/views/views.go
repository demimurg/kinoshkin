package views

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"kinoshkin/domain"
	"strings"
	"time"
)

var (
	CinemasCmd  = tb.ReplyButton{Text: "Кинотеатры🍿"}
	MoviesCmd   = tb.ReplyButton{Text: "Фильмы🎬"}
	LocationCmd = tb.ReplyButton{Text: "Обновить локацию📍", Location: true}
)

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

func CinemasList(cinemas []*domain.Cinema) [][]tb.InlineButton {
	var table [][]tb.InlineButton
	for _, cinema := range cinemas {
		table = append(table, []tb.InlineButton{
			{
				Text: fmt.Sprintf("%s ~ %.2fкм", cinema.Name, float32(cinema.Distance)/1000),
				Data: Encode(CinemaPrefix, cinema.ID),
			},
		})
	}

	return table
}

func MovieCard(mov *domain.Movie) (msg interface{}, opts []interface{}) {
	title := applyEscaping(fmt.Sprintf(
		"*%s* _(imdb: %.1f, kp: %.1f)_",
		mov.Title, mov.Rating.IMDB, mov.Rating.KP,
	))
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
		}, []interface{}{tb.ModeMarkdownV2, &tb.ReplyMarkup{
			InlineKeyboard: [][]tb.InlineButton{{{
				Text: "Где посмотреть?🙈",
				Data: Encode(MovieSchedulePrefix, mov.ID),
			}}},
		}}
}

func CinemaCard(cinema *domain.Cinema, schedule map[*domain.Movie][]domain.Session) (msg interface{}, opts []interface{}) {
	address := cinema.Address
	if len(cinema.Metro) > 0 {
		address = "🚇" + strings.Join(cinema.Metro, ", ") + "\n" + address
	}

	formattedSchedule := make(map[string][]domain.Session, len(schedule))
	for mov, sess := range schedule {
		movTitle := fmt.Sprintf("%s (%.1f)", mov.Title, mov.Rating.KP)
		formattedSchedule[movTitle] = sess
	}

	return &tb.Venue{
			Location:     tb.Location{Lat: cinema.Lat, Lng: cinema.Long},
			Title:        fmt.Sprintf("%s ~ %.2fкм", cinema.Name, float32(cinema.Distance)/1000),
			Address:      address,
			FoursquareID: "4bf58dd8d48988d17f941735",
		}, []interface{}{&tb.ReplyMarkup{
			InlineKeyboard: scheduleTable(formattedSchedule),
		}}
}

func MovieScheduleTable(schedule map[*domain.Cinema][]domain.Session) (interface{}, []interface{}) {
	formattedSchedule := make(map[string][]domain.Session, len(schedule))
	for cin, sess := range schedule {
		cinemaTitle := fmt.Sprintf("%s ~ %.2fкм", cin.Name, float32(cin.Distance)/1000)
		formattedSchedule[cinemaTitle] = sess
	}

	msg := "_Расписание на cегодня_ " + time.Now().Format("02-01")
	return applyEscaping(msg), []interface{}{tb.ModeMarkdownV2, &tb.ReplyMarkup{
		InlineKeyboard: scheduleTable(formattedSchedule),
	}}
}

func scheduleTable(schedule map[string][]domain.Session) (table [][]tb.InlineButton) {
	for title, sess := range schedule {
		table = append(table, []tb.InlineButton{{
			Text: title,
		}})

		var sessions []tb.InlineButton
		for i, ses := range sess {
			sessions = append(sessions, tb.InlineButton{
				Text: ses.Start.Format("15:04") + fmt.Sprintf(" %dр", ses.Price),
			})
			if i%2 != 0 {
				table = append(table, sessions)
				sessions = nil
			}
		}
		if len(sessions) > 0 {
			table = append(table, sessions)
		}
	}
	return
}
