package aggregator

import (
	"kinoshkin/pkg/set"
	"testing"

	"github.com/kr/pretty"
)

func TestKpApi(t *testing.T) {
	kp := &kpAPI{seenMovies: set.New()}
	kp.aggregateCinemaData("57f03c78b4660194c141e900")
	pretty.Println(kp.result())
}
