package storage

import (
	"context"
	"garantex-monitor/internal/models"
)

type GMStorageInterface interface {
	AddRate(context.Context, *models.DepthDTO) error
}
