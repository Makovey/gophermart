package transport

import (
	"context"

	"github.com/Makovey/gophermart/internal/transport/http/model"
)

//go:generate mockgen -source=service.go -destination=../service/mocks/service_mock.go -package=mocks
type GophermartService interface {
	UserService
	OrderService
	BalanceService
	HistoryService
}

type UserService interface {
	RegisterNewUser(ctx context.Context, request model.AuthRequest) (string, error)
	LoginUser(ctx context.Context, request model.AuthRequest) (string, error)
}

type OrderService interface {
	ValidateOrderID(orderID string) bool
	ProcessNewOrder(ctx context.Context, userID, orderID string) error
	GetOrders(ctx context.Context, userID string) ([]model.OrderResponse, error)
}

type BalanceService interface {
	GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error)
	WithdrawUsersBalance(ctx context.Context, userID string, withdraw model.WithdrawRequest) error
}

type HistoryService interface {
	GetUsersWithdrawHistory(ctx context.Context, userID string) ([]model.WithdrawHistoryResponse, error)
}
