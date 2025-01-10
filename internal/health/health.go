package health

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type GMHealhController struct {
	grpc_health_v1.UnimplementedHealthServer
	db     *pgx.Conn
	logger *zap.Logger
}

func NewGMHealhController(db *pgx.Conn, logger *zap.Logger) *GMHealhController {
	return &GMHealhController{
		db:     db,
		logger: logger,
	}
}

func (hc *GMHealhController) Check(
	ctx context.Context,
	in *grpc_health_v1.HealthCheckRequest,
) (*grpc_health_v1.HealthCheckResponse, error) {
	if err := hc.db.Ping(ctx); err != nil {
		hc.logger.Error("no database connection")
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	hc.logger.Debug("successful server healthcheck")
	return &grpc_health_v1.HealthCheckResponse{

		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}