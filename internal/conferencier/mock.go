package conferencier

import (
	"kinoshkin/domain"
	"math/rand"
	"sort"
	"syreclabs.com/go/faker"
	"time"
)

type Mock struct {
	domain.Conferencier
}

func (m Mock) FindMovies(userID int, pag domain.P) ([]*domain.Movie, error) {
	sort.Slice(mockMovies, func(i, j int) bool {
		return mockMovies[i].Rating.KP > mockMovies[j].Rating.KP
	})
	return mockMovies, nil
}

func (m Mock) FindCinemas(userID int, pag domain.P) ([]*domain.Cinema, error) {
	sort.Slice(mockCinemas, func(i, j int) bool {
		return mockCinemas[i].Distance < mockCinemas[j].Distance
	})
	return mockCinemas, nil
}

func (m Mock) GetMovie(movieID string) (*domain.Movie, error) {
	for _, mov := range mockMovies {
		if mov.ID == movieID {
			return mov, nil
		}
	}
	return nil, nil
}

func (m Mock) GetMovieSchedule(userID int, movieID string, pag domain.P) ([]domain.CinemaWithSessions, error) {
	var schedule []domain.CinemaWithSessions
	for _, cin := range mockCinemas {
		schedule = append(schedule, domain.CinemaWithSessions{
			Cinema:   cin,
			Sessions: genMockSessions(),
		})
	}
	sort.Slice(schedule, func(i, j int) bool {
		return schedule[i].Cinema.Distance < schedule[j].Cinema.Distance
	})
	return schedule, nil
}

func (m Mock) GetCinema(cinemaID string) (*domain.Cinema, error) {
	for _, cin := range mockCinemas {
		if cin.ID == cinemaID {
			return cin, nil
		}
	}
	return nil, nil
}

func (m Mock) GetCinemaSchedule(cinemaID string, pag domain.P) ([]domain.MovieWithSessions, error) {
	var schedule []domain.MovieWithSessions
	for _, mov := range mockMovies {
		schedule = append(schedule, domain.MovieWithSessions{
			Movie:    mov,
			Sessions: genMockSessions(),
		})
	}
	sort.Slice(schedule, func(i, j int) bool {
		return schedule[i].Movie.Rating.KP > schedule[j].Movie.Rating.KP
	})
	return schedule, nil
}

func (m Mock) UpdateUserLocation(userID int, lat, long float32) error {
	return nil
}
func (m Mock) RegisterUser(userID int, name string) error {
	return nil
}

func genMockSessions() []domain.Session {
	var sessions []domain.Session
	for i := 0; i < rand.Intn(5)+2; i++ {
		eveningTime, _ := time.Parse("15:04", "19:00")
		sessions = append(sessions, domain.Session{
			Start: faker.Time().Between(
				eveningTime.Truncate(7*time.Hour),
				eveningTime,
			),
			Price: int32(faker.RandomInt(100, 500)),
		})
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].Start.Before(sessions[j].Start)
	})

	return sessions
}

var mockMovies = []*domain.Movie{
	{
		ID:        "58972c13682d1a7c201f95b7c",
		Title:     "Патерсон",
		PosterURL: "https://avatars.mds.yandex.net/get-afishanew/29022/9c53aecc9c3491cd698ff36bdfb6273d/s744x446",
		Duration:  118,
		FilmCrew: map[domain.Position][]string{
			domain.Director: {"Джим Джармуш"},
			domain.Actor:    {"Адам Драйвер", "Голшифте Фарахани"},
		},
		Rating:      domain.Rating{KP: 7.3, IMDB: 7.4},
		Description: "Жизнь Патерсона – сплошная романтика: он работает водителем автобуса в городе Патерсон, штат Нью-Джерси, а в свободное время пишет стихи для любимой жены Лоры. Патерсон облачает красоту повседневности в стихи и встречает поэтов повсюду – такова магия города - родины поэтов Аллена Гинзберга и Уильяма Карлоса Уильямса. Патерсон пишет в стол и даже не мечтает публиковаться, однако одно маленькое событие меняет его планы.",
	},
	{
		ID:        "5e3d0adf47b50c9ff6b607c3",
		Title:     "Непосредственно Каха!",
		PosterURL: "https://avatars.mds.yandex.net/get-afishanew/36842/a3f951b2fd0dbf9b414abdf07f6d752c/s744x446",
		Duration:  110,
		FilmCrew: map[domain.Position][]string{
			domain.Director: {"Виктор Шамиров"},
			domain.Actor:    {"Артем Калайджян", "Артем Карокозян"},
		},
		Rating:      domain.Rating{KP: 5.5},
		Description: "Разборки в Сочи: герой на «копейке» запал на девушку обладателя BMW",
	},
	{
		ID:        "5dfb89a47a972aad9ae9f915",
		Title:     "Обратная связь",
		PosterURL: "https://avatars.mds.yandex.net/get-afishanew/21422/376997cbc9ebba604e8ac336544246b1/s744x446",
		Duration:  97,
		FilmCrew: map[domain.Position][]string{
			domain.Director: {"Алексей Нужный"},
			domain.Actor:    {"Ростислав Хаит", "Мария Миронова"},
		},
		Rating:      domain.Rating{KP: 6.7},
		Description: "Ирина Горбачёва и «Квартет И» в новогоднем сиквеле «Громкой связи»",
	},
	{
		ID:        "5575fad0cc1c725c1b9865f2",
		Title:     "Семейка Крудс: Новоселье",
		PosterURL: "https://avatars.mds.yandex.net/get-afishanew/23222/1a9348c9850df2484be152024d6ed70d/s744x446",
		Duration:  95,
		FilmCrew: map[domain.Position][]string{
			domain.Director: {"Джоэль Кроуфорд"},
			domain.Actor:    {"Райан Рейнольдс", "Эмма Стоун"},
		},
		Rating:      domain.Rating{KP: 7.1, IMDB: 7.0},
		Description: "Новая часть франшизы о жизни первобытных людей",
	},
}

var mockCinemas = []*domain.Cinema{
	{
		ID:       "551797ea1f7d154a12ddf058",
		Name:     "Формула Кино Академ Парк",
		Address:  "Гражданский просп., 41, ТРК «Академ-Парк»",
		Metro:    []string{"Академическая", "Политехническая"},
		Lat:      60.011567,
		Long:     30.397759,
		Distance: 1231,
	},
	{
		ID:       "57f03ccc682d1a1d756e2575",
		Name:     "Мираж Синема Озерки",
		Address:  "просп. Энгельса, 124",
		Metro:    []string{"Озерки", "Проспект Просвещения"},
		Lat:      60.040254,
		Long:     30.324294,
		Distance: 2015,
	},
	{
		ID:       "55424df91f6fd64c8d37933b",
		Name:     "Формула Кино Галерея",
		Address:  "Лиговский просп., 30а, ТРЦ «Галерея», 4 этаж",
		Metro:    []string{"Владимирская", "Площадь Восстания"},
		Lat:      59.92741,
		Long:     30.36064,
		Distance: 7317,
	},
	{
		ID:       "578c1748682d1ade3a28bea3",
		Name:     "Формула Кино Родео Драйв",
		Address:  "просп. Культуры, 1, ТРК «Родео Драйв», 3 этаж",
		Metro:    []string{"Академическая", "Политехническая"},
		Lat:      60.011567,
		Long:     30.397759,
		Distance: 1302,
	},
	{
		ID:       "57f03c78b4660194c141e900",
		Name:     "Мираж Синема Европолис",
		Address:  "Полюстровский просп., 84а, ТРК «Европолис»",
		Metro:    []string{"Лесная"},
		Lat:      59.987456,
		Long:     30.354913,
		Distance: 1504,
	},
}
