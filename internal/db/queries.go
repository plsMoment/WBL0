package db

import (
	"WBL0/internal/model"
	"context"
	"github.com/jackc/pgx/v5"
)

const (
	getAllStmt = "SELECT * FROM orders"
	createStmt = "INSERT INTO orders (id, data) VALUES ($1, $2)"
)

type OrderStorage interface {
	GetAllOrders(ctx context.Context) ([]model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) error
}

// GetAllOrders returns all orders from database
func (s *Storage) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	rows, err := s.pool.Query(ctx, getAllStmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Order])
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// CreateOrder creates order in database
func (s *Storage) CreateOrder(ctx context.Context, order model.Order) error {
	_, err := s.pool.Exec(ctx, createStmt, order.Id, order.Data)
	if err != nil {
		return err
	}
	return err
}
