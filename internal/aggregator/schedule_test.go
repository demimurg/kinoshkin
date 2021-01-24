package aggregator

import "testing"

func Test_convertToMins(t *testing.T) {
	tests := []struct {
		dur  string
		want int
	}{
		{"1:37", 97},
		{"0:24", 24},
		{"abc", 0},
	}
	for _, tt := range tests {
		if got := convertToMins(tt.dur); got != tt.want {
			t.Errorf("convertToMins() = %v, want %v", got, tt.want)
		}
	}
}
