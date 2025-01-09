package controller

import (
	"context"
	pb "garantex-monitor/gen/grpc"
	"garantex-monitor/internal/service"

	"google.golang.org/protobuf/types/known/emptypb"

	"go.uber.org/zap"
)

type Controller struct {
	gmService service.GarantexMotitorIface
	pb.UnimplementedGarantexMonitorServer
}

func NewController(gmService service.GarantexMotitorIface, logger *zap.Logger) *Controller {
	return &Controller{
		gmService: gmService,
	}
}

func (c *Controller) GetRates(ctx context.Context, empty *emptypb.Empty) (*pb.GetRatesResponse, error) {
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
