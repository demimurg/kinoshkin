package conferencier

import "kinoshkin/domain"

type Mock struct {
	domain.Conferencier
}

func (m Mock) FindMovies(userID int, pag domain.P) ([]*domain.Movie, error) {
	return mockMovies, nil
}

func (m Mock) FindCinemas(userID int, pag domain.P) ([]*domain.Cinema, []int, error) {
	return mockCinemas, []int{1231, 2015, 7317, 1302, 1504}, nil
}

func (m Mock) GetMovie(movieID string) (*domain.Movie, error) {
	for _, mov := range mockMovies {
		if mov.ID == movieID {
			return mov, nil
		}
	}
	return nil, nil
}

func (m Mock) GetCinema(cinemaID string) (*domain.Cinema, error) {
	for _, cin := range mockCinemas {
		if cin.ID == cinemaID {
			return cin, nil
		}
	}
	return nil, nil
}

func (m Mock) UpdateUserLocation(userID int, lat, long float32) error {
	return nil
}
func (m Mock) RegisterUser(userID int, name string) error {
	return nil
}

var mockMovies = []*domain.Movie{
	{
		ID:             "58972c13682d1a7c201f95b7c",
		Title:          "Патерсон",
		PosterURL:      "https://avatars.mds.yandex.net/get-afishanew/29022/9c53aecc9c3491cd698ff36bdfb6273d/s744x446",
		Duration:       118,
		AgeRestriction: true,
		FilmCrew: map[domain.Position]domain.Persons{
			domain.Director: {"Джим Джармуш"},
			domain.Actor:    {"Адам Драйвер", "Голшифте Фарахани", "Нелли"},
		},
		Rating: domain.Rating{KP: 7.3, IMDB: 7.4},
	},
	{
		ID:        "5e3d0adf47b50c9ff6b607c3",
		Title:     "Непосредственно Каха!",
		PosterURL: "https://avatars.mds.yandex.net/get-afishanew/36842/a3f951b2fd0dbf9b414abdf07f6d752c/s744x446",
		Duration:  110,
		FilmCrew: map[domain.Position]domain.Persons{
			domain.Director: {"Виктор Шамиров"},
			domain.Actor:    {"Артем Калайджян", "Артем Карокозян", "Виктор Шамиров"},
		},
		Rating: domain.Rating{KP: 5.5},
	},
	{
		ID:        "5dfb89a47a972aad9ae9f915",
		Title:     "Обратная связь",
		PosterURL: "https://avatars.mds.yandex.net/get-afishanew/21422/376997cbc9ebba604e8ac336544246b1/s744x446",
		Duration:  97,
		FilmCrew: map[domain.Position]domain.Persons{
			domain.Director: {"Алексей Нужный"},
			domain.Actor:    {"Ростислав Хаит", "Мария Миронова", "Леонид Барац"},
		},
		Rating: domain.Rating{KP: 6.7},
	},
	{
		ID:        "5575fad0cc1c725c1b9865f2",
		Title:     "Семейка Крудс: Новоселье",
		PosterURL: "https://avatars.mds.yandex.net/get-afishanew/23222/1a9348c9850df2484be152024d6ed70d/s744x446",
		Duration:  95,
		FilmCrew: map[domain.Position]domain.Persons{
			domain.Director: {"Джоэль Кроуфорд"},
			domain.Actor:    {"Райан Рейнольдс", "Эмма Стоун", "Николас Кейдж"},
		},
		Rating: domain.Rating{KP: 7.1, IMDB: 7.0},
	},
}

var mockCinemas = []*domain.Cinema{
	{
		ID:      "551797ea1f7d154a12ddf058",
		Name:    "Формула Кино Академ Парк",
		Address: "Гражданский просп., 41, ТРК «Академ-Парк»",
		Metro:   []string{"Академическая", "Политехническая"},
		Lat:     60.011567,
		Long:    30.397759,
	},
	{
		ID:      "57f03ccc682d1a1d756e2575",
		Name:    "Мираж Синема Озерки",
		Address: "просп. Энгельса, 124",
		Metro:   []string{"Озерки", "Проспект Просвещения"},
		Lat:     60.040254,
		Long:    30.324294,
	},
	{
		ID:      "55424df91f6fd64c8d37933b",
		Name:    "Формула Кино Галерея",
		Address: "Лиговский просп., 30а, ТРЦ «Галерея», 4 этаж",
		Metro:   []string{"Владимирская", "Площадь Восстания"},
		Lat:     59.92741,
		Long:    30.36064,
	},
	{
		ID:      "578c1748682d1ade3a28bea3",
		Name:    "Формула Кино Родео Драйв",
		Address: "просп. Культуры, 1, ТРК «Родео Драйв», 3 этаж",
		Metro:   []string{"Академическая", "Политехническая"},
		Lat:     60.011567,
		Long:    30.397759,
	},
	{
		ID:      "57f03c78b4660194c141e900",
		Name:    "Мираж Синема Европолис",
		Address: "Полюстровский просп., 84а, ТРК «Европолис»",
		Metro:   []string{"Лесная"},
		Lat:     59.987456,
		Long:    30.354913,
	},
}
