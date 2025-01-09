package service

import (
	"context"
)

type GarantexMotitorIface interface {
	GetRates(ctx context.Context) (*GetRatesOut, error)
}

type GetRatesOut struct {
	AskPrice  string
	BidPrice  string
	Timestamp int64
}
