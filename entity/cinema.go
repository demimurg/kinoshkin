package entity

// Cinema represents movie theater
type Cinema struct {
	ID        string
	Name      string
	Address   string
	Metro     []string
	Lat, Long float64
	Distance  int
}
