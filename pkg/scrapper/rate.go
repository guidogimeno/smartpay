package scrapper

import (
	"time"
)

type Ratable interface {
	Rate(date time.Time) (*Rate, error)
}

type Rate struct {
	Source string
	Index  float64
}
