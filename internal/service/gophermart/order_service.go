package gophermart

import (
	"context"
	"errors"
	"strconv"

	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/adapter"
	"github.com/Makovey/gophermart/internal/service/luhn"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

//go:generate mockgen -source=order_service.go -destination=../../repository/mocks/order_mock.go -package=mocks
type OrderServiceRepository interface {
	GetOrderByID(ctx context.Context, orderID string) (repoModel.Order, error)
	GetOrders(ctx context.Context, userID string) ([]repoModel.Order, error)
	PostNewOrder(ctx context.Context, orderID, userID string) error
}

type orderService struct {
	repo OrderServiceRepository
}

func NewOrderService(
	repo OrderServiceRepository,
) transport.OrderService {
	return &orderService{
		repo: repo,
	}
}

func (o *orderService) ValidateOrderID(orderID string) bool {
	orderInt, err := strconv.Atoi(orderID)
	if err != nil {
		return false
	}

	return luhn.IsValid(orderInt)
}

func (o *orderService) ProcessNewOrder(ctx context.Context, userID, orderID string) error {
	order, err := o.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return o.repo.PostNewOrder(ctx, orderID, userID)
		default:
			return err
		}
	}

	if order.OwnerUserID != userID {
		return service.ErrOrderConflict
	}

	return service.ErrOrderAlreadyPosted
}

func (o *orderService) GetOrders(ctx context.Context, userID string) ([]model.OrderResponse, error) {
	repoOrders, err := o.repo.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return adapter.FromRepoToOrders(repoOrders), nil
}
