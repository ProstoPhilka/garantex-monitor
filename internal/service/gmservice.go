package service

import (
	"context"
	"encoding/json"
	"fmt"
	"garantex-monitor/internal/models"
	"garantex-monitor/internal/storage"
	"garantex-monitor/pkg/metrics"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	URL = "https://garantex.org/api/v2/depth?market=usdtrub"
)

type GMService struct {
	logger  *zap.Logger
	client  *http.Client
	storage storage.GMStorageIface
}

func NewGMService(storage storage.GMStorageIface, logger *zap.Logger) *GMService {
	return &GMService{
		logger:  logger,
		client:  &http.Client{},
		storage: storage,
	}
}

func (g *GMService) GetRates(ctx context.Context) (*GetRatesOut, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			g.logger.Error("Failed to close body response", zap.Error(err))
		}
	}()

	var res models.Depth
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(res.Asks) == 0 || len(res.Bids) == 0 {
		return nil, fmt.Errorf("no data in response")
	}

	// TODO: save to DB
	go func() {
		// Метрики
		timer := prometheus.NewTimer(metrics.DbDurationSeconds.WithLabelValues("AddRate"))
		defer timer.ObserveDuration()

		err := g.storage.AddRate(
			context.Background(), // TODO добавить контекст
			&storage.AddRateIn{
				Timestamp: time.Unix(res.Timestamp, 0),
				Ask:       res.Asks[0].Price,
				Bid:       res.Bids[0].Price,
			})
		if err != nil {
			g.logger.Error("failed to add rate", zap.Error(err))
			return
		}
	}()

	return &GetRatesOut{
		AskPrice:  res.Asks[0].Price,
		BidPrice:  res.Bids[0].Price,
		Timestamp: res.Timestamp,
	}, nil
}
