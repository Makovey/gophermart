package gophermart

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/luhn"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
	"github.com/Makovey/gophermart/pkg/jwt"
)

type orderService struct {
	repo service.OrderRepository
	log  logger.Logger
	jwt  *jwt.JWT
}

func newOrderService(
	repo service.OrderRepository,
	log logger.Logger,
	jwt *jwt.JWT,
) transport.OrderService {
	return &orderService{
		repo: repo,
		log:  log,
		jwt:  jwt,
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

	var models []model.OrderResponse
	for _, repOrder := range repoOrders {
		var accrual *types.FloatDecimal
		if repOrder.Accrual != nil {
			val := types.FloatDecimal(*repOrder.Accrual)
			accrual = &val
		}

		order := model.OrderResponse{
			Number:     repOrder.OrderID,
			Status:     repOrder.Status,
			Accrual:    accrual,
			UploadedAt: repOrder.CreatedAt.Format(time.RFC3339),
		}

		models = append(models, order)
	}

	return models, nil
}
