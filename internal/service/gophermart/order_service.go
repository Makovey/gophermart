package gophermart

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/adapter"
	"github.com/Makovey/gophermart/internal/service/luhn"
	"github.com/Makovey/gophermart/internal/transport/http"
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
) http.OrderService {
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
	fn := "gophermart.ProcessNewOrder"

	order, err := o.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return o.repo.PostNewOrder(ctx, orderID, userID)
		default:
			return fmt.Errorf("[%s]: %w", fn, err)
		}
	}

	if order.OwnerUserID != userID {
		return fmt.Errorf("[%s]: %w", fn, service.ErrOrderConflict)
	}

	return fmt.Errorf("[%s]: %w", fn, service.ErrOrderAlreadyPosted)
}

func (o *orderService) GetOrders(ctx context.Context, userID string) ([]model.OrderResponse, error) {
	fn := "gophermart.GetOrders"

	repoOrders, err := o.repo.GetOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", fn, err)
	}

	return adapter.FromRepoToOrders(repoOrders), nil
}
