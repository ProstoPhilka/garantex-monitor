package controller

import (
	"context"
	"garantex-monitor/gen/pb"
	"garantex-monitor/internal/service"
	"garantex-monitor/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GMController struct {
	gmService service.GMServiceIface
	logger    *zap.Logger
	pb.UnimplementedGarantexMonitorServer
}

func NewGMController(gmService service.GMServiceIface, logger *zap.Logger, tracerName string) *GMController {
	return &GMController{
		gmService: gmService,
		logger:    logger,
	}
}

func (c *GMController) GetRates(ctx context.Context, empty *emptypb.Empty) (*pb.GetRatesResponse, error) {
	// Метрики
	timer := prometheus.NewTimer(metrics.ApiDurationSeconds.WithLabelValues("GetRates"))
	defer timer.ObserveDuration()

	depthOut, err := c.gmService.GetRates(ctx)
	if err != nil {
		c.logger.Error("failed get rates", zap.Error(err))
		return nil, status.Error(codes.Internal, "internal service error")
	}

	return &pb.GetRatesResponse{
		AskPrice:  depthOut.AskPrice,
		BidPrice:  depthOut.BidPrice,
		Timestamp: depthOut.Timestamp,
	}, nil
}
