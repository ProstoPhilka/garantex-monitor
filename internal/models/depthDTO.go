package models

import "time"

type DepthDTO struct {
	Timestamp time.Time
	Ask       string
	Bid       string
}
