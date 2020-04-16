package core

import (
	"time"
)

const Version = 0.1

var (
	MinTime = time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	MaxTime = time.Date(9999, time.December, 31, 23, 59, 0, 0, time.UTC)
)
