package storage

import (
	"context"
	"garantex-monitor/internal/models"
)

type GMStorageIface interface {
	AddRate(context.Context, *models.DepthDTO) error
}
