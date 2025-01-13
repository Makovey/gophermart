package postgresql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

func (r *Repo) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	fn := "postgresql.GetOrderByID"

	row := r.pool.QueryRow(
		ctx,
		`SELECT order_id, owner_user_id, status, accrual FROM gophermart_orders WHERE order_id = $1`,
		orderID,
	)
	var order model.Order
	err := row.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return model.Order{}, fmt.Errorf("[%s] user with order %s not found: %w", fn, orderID, service.ErrNotFound)
		default:
			return model.Order{}, fmt.Errorf("[%s] failed to query user : %w", fn, service.ErrExecStmt)
		}
	}
	return order, nil
}

func (r *Repo) PostNewOrder(ctx context.Context, orderID, status, userID string) error {
	fn := "postgresql.PostNewOrder"

	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO gophermart_orders (order_id, owner_user_id, status, created_at) VALUES ($1, $2, $3, $4)`,
		orderID,
		userID,
		status,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("[%s] failed to insert user : %w", fn, service.ErrExecStmt)
	}

	return nil
}

func (r *Repo) GetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	fn := "postgresql.GetOrders"

	rows, err := r.pool.Query(
		ctx,
		`SELECT * FROM gophermart_orders WHERE owner_user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed to query orders: %w", fn, err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("[%s] failed to scan orders: %w", fn, err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("[%s] failed to iterate orders: %w", fn, err)
	}

	return orders, nil
}

func (r *Repo) FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- model.Order, newStatus, inProgressStatus string) error {
	fn := "postgresql.FetchNewOrdersToChan"

	rows, err := r.pool.Query(
		ctx,
		`SELECT * FROM gophermart_orders WHERE status = $1 OR status = $2 ORDER BY created_at`,
		newStatus,
		inProgressStatus,
	)
	if err != nil {
		return fmt.Errorf("[%s] failed to query orders: %w", fn, err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			return fmt.Errorf("[%s] failed to scan orders: %w", fn, err)
		}

		ordersCh <- order
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("[%s] failed to iterate orders: %w", fn, err)
	}

	return nil
}

func (r *Repo) UpdateOrder(ctx context.Context, status model.OrderStatus) error {
	fn := "postgresql.UpdateOrder"

	res, err := r.pool.Exec(
		ctx,
		`UPDATE gophermart_orders SET status = $1, accrual = $2 WHERE order_id = $3`,
		status.Status,
		status.Accrual,
		status.OrderID,
	)
	if err != nil {
		return fmt.Errorf("[%s]: failed to update order: %w", fn, service.ErrExecStmt)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("[%s] didn't find order id, rows not affected: %w", fn, service.ErrNotFound)
	}

	return nil
}
