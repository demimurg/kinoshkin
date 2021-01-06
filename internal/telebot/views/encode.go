package views

import "strings"

const (
	// MoviePrefix used for identify views
	MoviePrefix = "mov"
	// CinemaPrefix used for identify views
	CinemaPrefix = "cin"
	// MoviesPrefix used for identify views
	MoviesPrefix = "movs"
	// CinemasPrefix used for identify views
	CinemasPrefix = "cins"
	// MovieSchedulePrefix appears when somebody requests schedule from MovieCard
	MovieSchedulePrefix = "msc"
)

func Encode(prefix string, id string) string {
	return prefix + "|" + id
}

func Decode(data string) (string, string) {
	// todo: error handling
	components := strings.Split(data, "|")
	return components[0], components[1]
}
