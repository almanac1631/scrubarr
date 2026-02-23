package domain

import (
	"time"
)

type Tracker struct {
	Name     string
	MinRatio float64
	MinAge   time.Duration
}
