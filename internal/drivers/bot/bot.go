package bot

import (
	"kinoshkin/internal/drivers/bot/views"
	"kinoshkin/internal/usecase"
	logger "log"
	"strings"
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
)

func log(err error, msg ...string) bool {
	if err != nil {
		err = errors.Wrap(err, strings.Join(msg, ""))
		logger.Print(err)
		return true
	}
	return false
}

func Start(
	svc usecase.Conferencier,
	token string, verboseLogs bool,
	telegramPollInterval time.Duration,
) {
	b, err := tb.NewBot(tb.Settings{
		Token:   token,
		Verbose: verboseLogs,
		Poller:  &tb.LongPoller{Timeout: telegramPollInterval},
	})
	log(err, "Bot initialization error")

	var limit = usecase.P{Limit: 6}

	b.Handle("/start", func(m *tb.Message) {
		err := svc.RegisterUser(
			int(m.Sender.ID), strings.Join([]string{
				m.Sender.FirstName, m.Sender.LastName,
			}, " "),
		)
		log(err, "User registration")

		_, err = b.Send(m.Sender, "Категорически приветствую!", &tb.ReplyMarkup{
			ReplyKeyboard: [][]tb.ReplyButton{
				{views.CinemasCmd, views.MoviesCmd},
				{views.LocationCmd},
			},
			ResizeReplyKeyboard: true,
		})
		log(err)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		var (
			msg     string
			buttons [][]tb.InlineButton
		)

		switch m.Text {
		case views.CinemasCmd.Text:
			cinemas, err := svc.FindCinemas(int(m.Sender.ID), limit)
			log(err, "Find cinemas")
			msg = "cinemas"
			buttons = views.CinemasList(cinemas)
		case views.MoviesCmd.Text:
			movies, err := svc.FindMovies(int(m.Sender.ID), limit)
			log(err, "Find movies")
			msg = "movies"
			buttons = views.MoviesList(movies)
		default:
			msg = "👌"
		}

		_, err := b.Send(m.Sender, msg, tb.ModeMarkdownV2, &tb.ReplyMarkup{
			InlineKeyboard: buttons,
		})
		log(err)
	})

	b.Handle(tb.OnCallback, func(cb *tb.Callback) {
		var (
			msg  interface{}
			opts []interface{}
		)

		switch prefix, id := views.Decode(cb.Data); prefix {
		case views.MoviePrefix:
			movie, err := svc.GetMovie(id)
			log(err, "Get movie ", id)
			msg, opts = views.MovieCard(movie)
		case views.CinemaPrefix:
			cinema, err := svc.GetCinema(id)
			log(err, "Get cinema ", id)

			schedule, err := svc.GetCinemaSchedule(id, limit)
			log(err, "Get cinema schedule")

			msg, opts = views.CinemaCard(cinema, schedule)
		case views.MovieSchedulePrefix:
			schedule, err := svc.GetMovieSchedule(int(cb.Sender.ID), id, limit)
			log(err, "Get movie schedule")

			msg, opts = views.MovieScheduleTable(schedule)
		default:
			return
		}

		_, err := b.Send(cb.Sender, msg, opts...)
		log(err)
	})

	b.Handle(tb.OnLocation, func(m *tb.Message) {
		err := svc.UpdateUserLocation(int(m.Sender.ID), m.Location.Lat, m.Location.Lng)
		log(err, "Update user location")

		_, err = b.Send(m.Sender, "Местоположение обновлено")
		log(err)
	})

	b.Start()
}
