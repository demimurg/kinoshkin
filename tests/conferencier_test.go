package tests

import (
	"context"
	"kinoshkin/internal/adapters/mongodb"
	"kinoshkin/internal/entity"
	"kinoshkin/internal/usecase"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestConferencier(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Conferencier Suite")
}

var (
	ctx = context.Background()

	// we will test this service
	conferencier usecase.Conferencier
	// we need repos to test that conferencier works correct
	cinemas  usecase.CinemasRepository
	movies   usecase.MoviesRepository
	users    usecase.UsersRepository
	schedule usecase.SchedulesRepository

	rudolf = entity.User{
		ID:   rand.Intn(1e6),
		Name: "Rudolf Kinchikov",
		City: "saint-petersburg",
		// пр-кт Раевского, д. 9
		Lat:  60.01469421386719,
		Long: 30.363439559936523,
	}
	citymoll = entity.Cinema{
		ID:      uuid.NewString(),
		Name:    "Формула Кино Сити Молл",
		Address: "Коломяжский просп., 17, корп. 1, ТРК «Сити Молл», 4 этаж",
		Metro:   []string{"Пионерская"},
		Lat:     60.005425,
		Long:    30.298915,
	}
	europolis = entity.Cinema{
		ID:      uuid.NewString(),
		Name:    "Мираж Синема Европолис",
		Address: "Полюстровский просп., 84а, ТРК «Европолис»",
		Metro:   []string{"Лесная"},
		Lat:     59.987456,
		Long:    30.354913,
	}
	soul = entity.Movie{
		ID:             uuid.NewString(),
		KpID:           "775273",
		Title:          "Душа",
		Description:    "Пронзительная история от студии Pixar: джазмен гибнет, но душа остаётся",
		PosterURL:      "https://avatars.mds.yandex.net/get-afishanew/28638/9198c97a8c9722b7e617d8a1e92f8ca2/s744x446",
		Duration:       106,
		AgeRestriction: "6+",
		FilmCrew: map[string][]string{
			entity.Director: {"Пит Доктер"},
			entity.Actor:    {"Джейми Фокс", "Тина Фей"},
		},
		Rating: entity.Rating{
			IMDB: 8.1,
			KP:   8.3,
		},
		DateReleased: date("2021-01-21"),
	}
	laLaLend = entity.Movie{
		ID:             "5874ea2a685ae0b186614bb5",
		KpID:           "841081",
		Title:          "Ла-Ла Ленд",
		Description:    "«Золотой глобус» за лучшую музыку к фильму",
		PosterURL:      "https://avatars.mds.yandex.net/get-afishanew/21626/efca9618f9297c5342f3459eace99db5/s744x446",
		Duration:       128,
		AgeRestriction: "16+",
		FilmCrew: map[string][]string{
			entity.Director:     {"Дэмьен Шазелл"},
			entity.Screenwriter: {"Дэмьен Шазелл"},
			entity.Actor:        {"Райан Гослинг", "Эмма Стоун"},
		},
		Rating: entity.Rating{
			IMDB: 8,
			KP:   7.9,
		},
		DateReleased: time.Time{},
	}
)

var _ = BeforeSuite(func() {
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(
		"mongodb://localhost:27017",
	))
	Expect(err).Should(BeNil())
	DeferCleanup(mongoClient.Disconnect, ctx)

	db := mongoClient.Database("kinoshkin")
	Expect(db.Drop(ctx)).Should(Succeed())

	for _, col := range []string{"users", "cinemas"} {
		_, err = db.Collection(col).Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		})
		Expect(err).To(BeNil())
	}

	cinemas = mongodb.NewCinemasRepository(db)
	movies = mongodb.NewMoviesRepository(db)
	users = mongodb.NewUsersRepository(db)
	schedule = mongodb.NewSchedulesRepository(db)
	conferencier = usecase.NewConferencier(cinemas, movies, users, schedule)

	By("create test user with geo location")
	err = conferencier.RegisterUser(rudolf.ID, rudolf.Name)
	Expect(err).Should(BeNil())
	Expect(conferencier.UpdateUserLocation(rudolf.ID, float32(rudolf.Lat), float32(rudolf.Long))).Should(Succeed())

	user, err := users.Get(rudolf.ID)
	Expect(err).To(BeNil())
	Expect(*user).To(Equal(rudolf))
})

var _ = Describe("work with cinemas", Ordered, func() {
	BeforeAll(func() {
		Expect(cinemas.Create([]entity.Cinema{citymoll, europolis})).Should(Succeed())

		// later we will only retrieve cinemas, so we can mutate it now
		// we need this because cinemas returns with distance property from db
		// europolis.Distance = dist(rudolf.Lat, rudolf.Long, europolis.Lat, europolis.Long)
		// citymoll.Distance = dist(rudolf.Lat, rudolf.Long, citymoll.Lat, citymoll.Long)
	})

	It("can get cinema by id", func() {
		cinema, err := conferencier.GetCinema(europolis.ID)
		Expect(err).To(BeNil())
		Expect(*cinema).To(Equal(europolis))
	})

	It("can find cinemas near user", func() {
		found, err := conferencier.FindCinemas(rudolf.ID, usecase.P{Limit: 5})
		Expect(err).To(BeNil())
		// they returns sorted by distance
		Expect(found).To(BeEquivalentTo([]entity.Cinema{europolis, citymoll}))
	})
})

var _ = Describe("work with movies", Ordered, func() {
	BeforeAll(func() {
		Expect(movies.Create([]entity.Movie{soul, laLaLend})).Should(Succeed())
	})

	It("can get movie by id", func() {
		movie, err := conferencier.GetMovie(soul.ID)
		Expect(err).To(BeNil())
		Expect(movie).ToNot(Equal(soul))
	})

	It("can find movies near user", func() {
		found, err := conferencier.FindMovies(rudolf.ID, usecase.P{Limit: 5})
		Expect(err).To(BeNil())
		// they returns sorted by distance
		Expect(found).ToNot(BeNil())
	})
})

// dist calculates geospatial distance for two points in meters
// func dist(latA, longA, latB, longB float64) int {
// 	_, km := haversine.Distance(
// 		haversine.Coord{Lat: latA, Lon: longA},
// 		haversine.Coord{Lat: latB, Lon: longB},
// 	)
// 	return int(km * 1000)
// }

func date(dateStr string) time.Time {
	parsedDate, _ := time.Parse(time.RFC3339, dateStr)
	return parsedDate
}
