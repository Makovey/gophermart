package transport

import (
	"context"

	"github.com/Makovey/gophermart/internal/transport/accrual/model"
)

type Accrual interface {
	RegisterNewGoods(ctx context.Context) error
	RegisterNewOrder(ctx context.Context, orderID string) error
	UpdateOrderStatus(ctx context.Context, orderID string) (model.OrderStatus, error)
}
