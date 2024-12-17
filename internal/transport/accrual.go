package transport

import "context"

type Accrual interface {
	RegisterNewGoods(ctx context.Context) error
	RegisterNewOrder(ctx context.Context, orderID string) error
	UpdateOrderStatus(ctx context.Context, orderID string) error
}
