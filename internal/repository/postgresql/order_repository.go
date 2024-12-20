package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

type orderRepository struct {
	log  logger.Logger
	pool *pgxpool.Pool
}

func newOrderRepository(log logger.Logger, pool *pgxpool.Pool) service.OrderRepository {
	return &orderRepository{
		log:  log,
		pool: pool,
	}
}

func (o *orderRepository) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	fn := "postgresql.GetOrderByID"

	row := o.pool.QueryRow(
		ctx,
		`SELECT order_id, owner_user_id, status, accrual FROM gophermart_orders WHERE order_id = $1`,
		orderID,
	)
	var order model.Order
	err := row.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			o.log.Info(fmt.Sprintf("%s: user with order %s not found", fn, orderID))
			return model.Order{}, service.ErrNotFound
		default:
			o.log.Error(fmt.Sprintf("%s: failed to query user", fn), "error", err)
			return model.Order{}, service.ErrExecStmt
		}
	}
	return order, nil
}

func (o *orderRepository) PostNewOrder(ctx context.Context, orderID, userID string) error {
	fn := "postgresql.PostNewOrder"

	_, err := o.pool.Exec(
		ctx,
		`INSERT INTO gophermart_orders (order_id, owner_user_id, status) VALUES ($1, $2, 'NEW')`,
		orderID,
		userID,
	)
	if err != nil {
		o.log.Error(fmt.Sprintf("%s: failed to post new order", fn), "error", err)
		return service.ErrExecStmt
	}

	return nil
}

func (o *orderRepository) GetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	fn := "postgresql.GetOrders"

	rows, err := o.pool.Query(
		ctx,
		`SELECT * FROM gophermart_orders WHERE owner_user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		o.log.Error(fmt.Sprintf("%s: failed to query orders", fn), "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			o.log.Error(fmt.Sprintf("%s: failed to scan orders", fn), "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		o.log.Error(fmt.Sprintf("%s: failed to iterate orders", fn), "error", err)
		return nil, err
	}

	return orders, nil
}

func (o *orderRepository) FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- model.Order) error {
	fn := "postgresql.FetchNewOrdersToChan"

	rows, err := o.pool.Query(
		ctx,
		`SELECT * FROM gophermart_orders WHERE status = 'NEW' OR status = 'PROCESSING' ORDER BY created_at`,
	)
	if err != nil {
		o.log.Error(fmt.Sprintf("%s: failed to query orders", fn), "error", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			o.log.Error(fmt.Sprintf("%s: failed to scan orders", fn), "error", err)
			return err
		}

		select {
		case <-ctx.Done():
			o.log.Info(fmt.Sprintf("%s: context cancelled", fn))
			return ctx.Err()
		case ordersCh <- order:
		}
	}

	if err = rows.Err(); err != nil {
		o.log.Error(fmt.Sprintf("%s: failed to iterate orders", fn), "error", err)
		return err
	}

	return nil
}

func (o *orderRepository) UpdateOrder(ctx context.Context, status model.OrderStatus) error {
	fn := "postgresql.UpdateOrder"

	res, err := o.pool.Exec(
		ctx,
		`UPDATE gophermart_orders SET status = $1, accrual = $2 WHERE order_id = $3`,
		status.Status,
		status.Accrual,
		status.OrderID,
	)
	if err != nil {
		o.log.Error(fmt.Sprintf("%s: failed to update order", fn), "error", err)
		return service.ErrExecStmt
	}

	if res.RowsAffected() == 0 {
		o.log.Error(fmt.Sprintf("%s: didn't find order id, rows not affected", fn))
		return service.ErrNotFound
	}

	return nil
}
