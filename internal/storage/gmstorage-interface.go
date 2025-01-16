package storage

import (
	"context"
	"time"
)

type GMStorageIface interface {
	AddRate(context.Context, *AddRateIn) error
}

type AddRateIn struct {
	Timestamp time.Time
	Ask       string
	Bid       string
}
