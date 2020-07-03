package bookng

import (
	"time"
)

const Version = 0.1

var (
	Location = time.Now().Location()
	MinTime = time.Date(1, time.January, 1, 0, 0, 0, 0, Location)
	MaxTime = time.Date(9999, time.December, 31, 23, 59, 0, 0, Location)
)

func Today() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
