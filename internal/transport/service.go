package transport

import (
	"context"

	"github.com/Makovey/gophermart/internal/transport/http/model"
)

//go:generate mockgen -source=service.go -destination=../service/mocks/service_mock.go -package=mocks
type GophermartService interface {
	RegisterNewUser(ctx context.Context, request model.AuthRequest) (string, error)
	LoginUser(ctx context.Context, request model.AuthRequest) (string, error)

	ValidateOrderID(orderID string) bool
	ProcessNewOrder(ctx context.Context, userID, orderID string) error
	GetOrders(ctx context.Context, userID string) ([]model.Order, error)
}
