package service

import (
	"context"
)

type GMServiceIface interface {
	GetRates(context.Context) (*GetRatesOut, error)
}

type GetRatesOut struct {
	AskPrice  string
	BidPrice  string
	Timestamp int64
}
