package storage

import (
	"context"
	"fmt"
	"garantex-monitor/internal/models"

	"github.com/jackc/pgx/v5"
)

type GMStorage struct {
	conn *pgx.Conn
}

func NewGMStorage(conn *pgx.Conn) *GMStorage {
	return &GMStorage{
		conn: conn,
	}
}

func (g *GMStorage) AddRate(ctx context.Context, depth *models.DepthDTO) error {
	query := "INSERT INTO usdt_rates (timestamp, ask, bid) VALUES ($1, $2, $3)"
	_, err := g.conn.Exec(ctx, query, depth.Timestamp, depth.Ask, depth.Bid)
	if err != nil {
		return fmt.Errorf("failed to insert rate: %w", err)
	}

	return nil
}
