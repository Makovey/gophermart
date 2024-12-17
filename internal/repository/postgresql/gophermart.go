package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

const (
	errUniqueViolatesCode = "23505"
)

type repo struct {
	log  logger.Logger
	conn *pgx.Conn
}

func NewPostgresRepo(log logger.Logger, cfg config.Config) service.GophermartRepository {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURI())
	if err != nil {
		log.Error("unable to connect to database", "error", err.Error())
		panic(err)
	}

	return &repo{log: log, conn: conn}
}

func (r *repo) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	fn := "postgresql.RegisterNewUser"

	_, err := r.conn.Exec(
		ctx,
		`INSERT INTO gophermart_users (user_id, login, password_hash) VALUES ($1, $2, $3)`,
		user.UserID,
		user.Login,
		user.PasswordHash,
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to execute new user", fn), "error", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok && pgErr.Code == errUniqueViolatesCode {
			return service.ErrLoginIsAlreadyExist
		}

		return service.ErrExecStmt
	}

	return nil
}

func (r *repo) LoginUser(ctx context.Context, login string) (model.RegisterUser, error) {
	fn := "postgresql.LoginUser"

	row := r.conn.QueryRow(
		ctx,
		`SELECT user_id, login, password_hash FROM gophermart_users WHERE login = $1`,
		login,
	)

	var user model.RegisterUser
	err := row.Scan(&user.UserID, &user.Login, &user.PasswordHash)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			r.log.Info(fmt.Sprintf("%s: user with login %s not found", fn, login))
			return model.RegisterUser{}, service.ErrNotFound
		default:
			r.log.Error(fmt.Sprintf("%s: failed to query user", fn), "error", err)
			return model.RegisterUser{}, service.ErrExecStmt
		}
	}

	return user, nil
}

func (r *repo) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	fn := "postgresql.GetOrderByID"

	row := r.conn.QueryRow(
		ctx,
		`SELECT order_id, owner_user_id, status, accrual FROM gophermart_orders WHERE order_id = $1`,
		orderID,
	)
	var order model.Order
	err := row.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			r.log.Info(fmt.Sprintf("%s: user with order %s not found", fn, orderID))
			return model.Order{}, service.ErrNotFound
		default:
			r.log.Error(fmt.Sprintf("%s: failed to query user", fn), "error", err)
			return model.Order{}, service.ErrExecStmt
		}
	}
	return order, nil
}

func (r *repo) PostNewOrder(ctx context.Context, orderID, userID string) error {
	fn := "postgresql.PostNewOrder"

	_, err := r.conn.Exec(
		ctx,
		`INSERT INTO gophermart_orders (order_id, owner_user_id, status) VALUES ($1, $2, 'NEW')`,
		orderID,
		userID,
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to post new order", fn), "error", err)
		return service.ErrExecStmt
	}

	return nil
}

func (r *repo) GetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	fn := "postgresql.GetOrders"

	rows, err := r.conn.Query(
		ctx,
		`SELECT * FROM gophermart_orders WHERE owner_user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to query orders", fn), "error", err)
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			r.log.Error(fmt.Sprintf("%s: failed to scan orders", fn), "error", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to iterate orders", fn), "error", err)
		return nil, err
	}

	return orders, nil
}

func (r *repo) FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- model.Order) error {
	fn := "postgresql.FetchNewOrdersToChan"

	rows, err := r.conn.Query(
		ctx,
		`SELECT * FROM gophermart_orders WHERE status = 'NEW' ORDER BY created_at`,
	)
	if err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to query orders", fn), "error", err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.OrderID, &order.OwnerUserID, &order.Status, &order.Accrual, &order.CreatedAt)
		if err != nil {
			r.log.Error(fmt.Sprintf("%s: failed to scan orders", fn), "error", err)
			return err
		}
		ordersCh <- order
	}

	if err = rows.Err(); err != nil {
		r.log.Error(fmt.Sprintf("%s: failed to iterate orders", fn), "error", err)
		return err
	}

	return nil
}

func (r *repo) Close() error {
	return r.conn.Close(context.Background())
}
