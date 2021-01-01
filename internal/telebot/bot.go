package telebot

import (
	"kinoshkin/domain"
	"kinoshkin/internal/telebot/views"
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

type bot struct {
	*tb.Bot
	domain.Conferencier
}

// New initialize handlers and start telegram bot server
func New() bot {
	b, err := tb.NewBot(tb.Settings{
		Token:   cfg.Token,
		Verbose: cfg.LogTrace,
		Poller:  &tb.LongPoller{Timeout: cfg.UpdateInterval},
	})
	if err != nil {
		log.Fatal("Bot initialization error: ", err)
	}

	return bot{b, nil}
}

func (b bot) Start() {
	var (
		cinemasCmd = tb.Command{Text: "cinemas", Description: "movie theaters near by you"}
		moviesCmd  = tb.Command{Text: "movies", Description: "actual movies related by rating"}
	)

	err := b.SetCommands([]tb.Command{cinemasCmd, moviesCmd})
	if err != nil {
		log.Fatal("Set commands error: ", err)
	}

	b.Handle(tb.OnText, func(m *tb.Message) {
		_, err := b.Send(m.Sender, "👌")
		log.Print(err)
	})

	b.Handle("/"+moviesCmd.Text, func(m *tb.Message) {
		movies, _ := b.FindMovies(m.Sender.ID, domain.P{})
		_, err := b.Send(m.Sender, "movies", &tb.ReplyMarkup{
			InlineKeyboard: views.MoviesTable(movies),
		})
		if err != nil {
			log.Print(err)
		}
	})
	b.Handle("/"+cinemasCmd.Text, func(m *tb.Message) {
		cinemas, distances, _ := b.FindCinemasNearby(m.Sender.ID, domain.P{})
		_, err := b.Send(m.Sender, "cinemas", &tb.ReplyMarkup{
			InlineKeyboard: views.CinemasTable(cinemas, distances),
		})
		if err != nil {
			log.Print(err)
		}
	})

	b.Handle(tb.OnCallback, func(cb *tb.Callback) {
		var (
			msg     string
			buttons [][]tb.InlineButton
		)

		switch prefix, id := views.Decode(cb.Data); prefix {
		case views.MoviePrefix:
			movie, _ := b.GetMovie(id)
			msg, buttons = views.MovieCard(movie)
		case views.CinemaPrefix:
			cinema, _, _ := b.GetCinema(id)
			msg, buttons = views.CinemaCard(cinema)
		case views.MoviesPrefix:
			movies, _ := b.FindMovies(cb.Sender.ID, domain.P{})
			msg = "movies"
			buttons = views.MoviesTable(movies)
		case views.CinemasPrefix:
			cinemas, dists, _ := b.FindCinemasNearby(cb.Sender.ID, domain.P{})
			msg = "cinemas"
			buttons = views.CinemasTable(cinemas, dists)
		}

		_, err := b.Send(cb.Sender, msg, tb.ModeMarkdownV2, &tb.ReplyMarkup{
			InlineKeyboard: buttons,
		})
		if err != nil {
			log.Print(err)
		}
	})

	b.Handle("/start", func(m *tb.Message) {
		errR := b.RegisterUser(
			m.Sender.ID, strings.Join([]string{
				m.Sender.FirstName, m.Sender.LastName,
			}, " "),
		)
		_, errS := b.Send(m.Sender, "Hello my friend!", &tb.ReplyMarkup{
			ReplyKeyboard: views.Commands(),
		})
		if errR != nil || errS != nil {
			log.Print(errR, errS)
		}
	})

	b.Handle(tb.OnLocation, func(m *tb.Message) {
		err := b.UpdateUserLocation(m.Sender.ID, m.Location.Lat, m.Location.Lng)
		if err != nil {
			log.Print(err)
		}
	})

	b.Start()
}
