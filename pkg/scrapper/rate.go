package scrapper

import (
	"context"
	"time"
)

type Ratable interface {
	Rate(ctx context.Context, date time.Time) (*Rate, error)
}

type Rate struct {
	Source string
	Index  float64
}
