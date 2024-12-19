package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

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
	conn *pgxpool.Pool

	userRepo     service.UserRepository
	orderRepo    service.OrderRepository
	balancesRepo service.BalancesRepository
}

func NewPostgresRepo(log logger.Logger, cfg config.Config) service.GophermartRepository {
	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI())
	if err != nil {
		log.Error("unable to connect to database", "error", err.Error())
		panic(err)
	}

	return &repo{
		log:          log,
		conn:         conn,
		userRepo:     NewUserRepository(log, conn),
		orderRepo:    NewOrderRepository(log, conn),
		balancesRepo: nil,
	}
}

func (r *repo) RegisterNewUser(ctx context.Context, user model.RegisterUser) error {
	return r.userRepo.RegisterNewUser(ctx, user)
}

func (r *repo) LoginUser(ctx context.Context, login string) (model.RegisterUser, error) {
	return r.userRepo.LoginUser(ctx, login)
}

func (r *repo) GetOrderByID(ctx context.Context, orderID string) (model.Order, error) {
	return r.orderRepo.GetOrderByID(ctx, orderID)
}

func (r *repo) PostNewOrder(ctx context.Context, orderID, userID string) error {
	return r.orderRepo.PostNewOrder(ctx, orderID, userID)
}

func (r *repo) GetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	return r.orderRepo.GetOrders(ctx, userID)
}

func (r *repo) FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- model.Order) error {
	return r.orderRepo.FetchNewOrdersToChan(ctx, ordersCh)
}

func (r *repo) UpdateOrder(ctx context.Context, status model.OrderStatus) error {
	return r.orderRepo.UpdateOrder(ctx, status)
}

func (r *repo) Close() error {
	r.conn.Close()
	return nil
}
