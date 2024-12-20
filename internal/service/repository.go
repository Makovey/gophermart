package service

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/Makovey/gophermart/internal/repository/model"
)

//go:generate mockgen -source=repository.go -destination=../repository/mocks/repository_mock.go -package=mocks
type GophermartRepository interface {
	UserRepository
	OrderRepository
	BalancesRepository
}

type UserRepository interface {
	RegisterNewUser(ctx context.Context, user model.RegisterUser) error
	LoginUser(ctx context.Context, login string) (model.RegisterUser, error)
}

type OrderRepository interface {
	GetOrderByID(ctx context.Context, orderID string) (model.Order, error)
	GetOrders(ctx context.Context, userID string) ([]model.Order, error)
	PostNewOrder(ctx context.Context, orderID, userID string) error
	FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- model.Order) error
	UpdateOrder(ctx context.Context, status model.OrderStatus) error
}

type BalancesRepository interface {
	UpdateUsersBalance(ctx context.Context, userID string, reward decimal.Decimal) error
	GetUsersBalance(ctx context.Context, userID string) (model.Balance, error)
}
