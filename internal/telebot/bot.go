package telebot

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"kinoshkin/domain"
	"kinoshkin/internal/telebot/views"
	"log"
	"strings"
)

type BotServer interface {
	Start()
}

// New initialize handlers
func New(svc domain.Conferencier) BotServer {
	b, err := tb.NewBot(tb.Settings{
		Token:   cfg.Token,
		Verbose: cfg.LogTrace,
		Poller:  &tb.LongPoller{Timeout: cfg.UpdateInterval},
	})
	if err != nil {
		log.Fatal("Bot initialization error: ", err)
	}

	b.Handle("/start", func(m *tb.Message) {
		errR := svc.RegisterUser(
			m.Sender.ID, strings.Join([]string{
				m.Sender.FirstName, m.Sender.LastName,
			}, " "),
		)
		_, errS := b.Send(m.Sender, "Hello my friend!", &tb.ReplyMarkup{
			ReplyKeyboard: [][]tb.ReplyButton{
				{views.CinemasCmd, views.MoviesCmd},
				{views.LocationCmd},
			},
			ResizeReplyKeyboard: true,
		})
		if errR != nil || errS != nil {
			log.Print(errR, errS)
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		var (
			msg     string
			buttons [][]tb.InlineButton
		)

		switch m.Text {
		case views.CinemasCmd.Text:
			cinemas, _ := svc.FindCinemas(m.Sender.ID, domain.P{Limit: 6})
			msg = "cinemas"
			buttons = views.CinemasList(cinemas)
		case views.MoviesCmd.Text:
			movies, _ := svc.FindMovies(m.Sender.ID, domain.P{Limit: 6})
			msg = "movies"
			buttons = views.MoviesList(movies)
		default:
			msg = "👌"
		}

		_, err := b.Send(m.Sender, msg, tb.ModeMarkdownV2, &tb.ReplyMarkup{
			InlineKeyboard: buttons,
		})
		if err != nil {
			log.Print(err)
		}
	})

	b.Handle(tb.OnCallback, func(cb *tb.Callback) {
		var (
			msg  interface{}
			opts []interface{}
		)

		switch prefix, id := views.Decode(cb.Data); prefix {
		case views.MoviePrefix:
			movie, _ := svc.GetMovie(id)
			msg, opts = views.MovieCard(movie)
		case views.CinemaPrefix:
			cinema, _ := svc.GetCinema(id)
			schedule, _ := svc.GetCinemaSchedule(id)
			msg, opts = views.CinemaCard(cinema, schedule)
		case views.MovieSchedulePrefix:
			schedule, _ := svc.GetMovieSchedule(cb.Sender.ID, id)
			msg, opts = views.MovieScheduleTable(schedule)
		default:
			return
		}

		_, err := b.Send(cb.Sender, msg, opts...)
		if err != nil {
			log.Print(err)
		}
	})

	b.Handle(tb.OnLocation, func(m *tb.Message) {
		err := svc.UpdateUserLocation(m.Sender.ID, m.Location.Lat, m.Location.Lng)
		if err != nil {
			log.Print(err)
		}
	})

	return b
}
