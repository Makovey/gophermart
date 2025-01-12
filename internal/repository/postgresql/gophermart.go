package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/shopspring/decimal"
	"path/filepath"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
)

const (
	migrationPath = "internal/db/migrations"
)

type repo struct {
	log  logger.Logger
	conn *pgxpool.Pool

	userRepo          service.UserRepository
	orderRepo         service.OrderRepository
	balancesRepo      service.BalancesRepository
	historyRepository service.HistoryRepository
}

func NewPostgresRepo(log logger.Logger, cfg config.Config) (service.GophermartRepository, error) {
	fn := "postgresql.NewPostgresRepo"

	path, err := filepath.Abs(migrationPath)
	if err != nil {
		return nil, fmt.Errorf("[%s]: could not determine absolute path for migrations: %w", fn, err)
	}

	err = upMigrations(cfg.DatabaseURI(), path)
	if err != nil {
		return nil, fmt.Errorf("[%s]: could not up migrations: %w", fn, err)
	}

	conn, err := pgxpool.New(context.Background(), cfg.DatabaseURI())
	if err != nil {
		return nil, fmt.Errorf("[%s]: could not connect to database: %w", fn, err)
	}

	return &repo{
		log:               log,
		conn:              conn,
		userRepo:          newUserRepository(log, conn),
		orderRepo:         newOrderRepository(log, conn),
		balancesRepo:      newBalancesRepository(log, conn),
		historyRepository: newHistoryRepository(log, conn),
	}, nil
}

func upMigrations(databaseURI, migrationsDir string) error {
	db, err := sql.Open("pgx", databaseURI)
	if err != nil {
		return err
	}
	defer db.Close()

	if err = goose.Up(db, migrationsDir); err != nil {
		return err
	}

	return nil
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

func (r *repo) IncreaseUsersBalance(ctx context.Context, userID string, reward decimal.Decimal) error {
	return r.balancesRepo.IncreaseUsersBalance(ctx, userID, reward)
}

func (r *repo) DecreaseUsersBalance(ctx context.Context, userID string, withdraw decimal.Decimal) error {
	return r.balancesRepo.DecreaseUsersBalance(ctx, userID, withdraw)
}

func (r *repo) GetUsersBalance(ctx context.Context, userID string) (model.Balance, error) {
	return r.balancesRepo.GetUsersBalance(ctx, userID)
}

func (r *repo) RecordUsersWithdraw(ctx context.Context, userID, orderID string, amount decimal.Decimal) error {
	return r.historyRepository.RecordUsersWithdraw(ctx, userID, orderID, amount)
}

func (r *repo) GetUsersHistory(ctx context.Context, userID string) ([]model.Withdraw, error) {
	return r.historyRepository.GetUsersHistory(ctx, userID)
}

func (r *repo) Close() error {
	r.conn.Close()
	return nil
}
