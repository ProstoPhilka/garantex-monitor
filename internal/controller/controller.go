package controller

import (
	"context"
	"garantex-monitor/gen/pb"
	"garantex-monitor/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"

	"go.uber.org/zap"
)

type GMController struct {
	gmService service.GMServiceIface
	pb.UnimplementedGarantexMonitorServer
}

func NewGMController(gmService service.GMServiceIface, logger *zap.Logger) *GMController {
	return &GMController{
		gmService: gmService,
	}
}

func (c *GMController) GetRates(ctx context.Context, empty *emptypb.Empty) (*pb.GetRatesResponse, error) {
	depthOut, err := c.gmService.GetRates(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetRatesResponse{
		AskPrice:  depthOut.AskPrice,
		BidPrice:  depthOut.BidPrice,
		Timestamp: depthOut.Timestamp,
	}, nil
}
