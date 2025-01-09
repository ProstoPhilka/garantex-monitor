package service

import (
	"context"
	"encoding/json"
	"fmt"
	"garantex-monitor/internal/models"
	"net/http"

	"go.uber.org/zap"
)

const (
	URL = "https://garantex.org/api/v2/depth?market=usdtrub"
)

type GarantexMonitor struct {
	logger *zap.Logger
	client *http.Client
}

func NewGarantexMonitor(logger *zap.Logger) *GarantexMonitor {
	return &GarantexMonitor{
		logger: logger,
		client: &http.Client{},
	}
}

func (g *GarantexMonitor) GetRates(ctx context.Context) (*GetRatesOut, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		g.logger.Error("Failed to create request", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	r, err := g.client.Do(req)
	if err != nil {
		g.logger.Error("Failed to execute request", zap.Error(err))
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			g.logger.Error("Failed to close body response", zap.Error(err))
		}
	}()

	var res models.Depth
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		g.logger.Error("Failed to unmarshal response", zap.Error(err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(res.Asks) == 0 || len(res.Bids) == 0 {
		g.logger.Error("No data in response")
		return nil, fmt.Errorf("no data in response")
	}

	// TODO: save to DB

	return &GetRatesOut{
		AskPrice:  res.Asks[0].Price,
		BidPrice:  res.Bids[0].Price,
		Timestamp: res.Timestamp,
	}, nil
}
