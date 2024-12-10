package service

import (
	"context"

	"github.com/Makovey/gophermart/internal/repository/model"
)

//go:generate mockgen -source=repository.go -destination=../repository/mocks/repository_mock.go -package=mocks
type GophermartRepository interface {
	RegisterNewUser(ctx context.Context, user model.RegisterUser) error
	LoginUser(ctx context.Context, login string) (model.RegisterUser, error)

	GetOrderByID(ctx context.Context, orderID string) (model.Order, error)
	PostNewOrder(ctx context.Context, orderID, userID string) error
}
