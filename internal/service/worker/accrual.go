package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Makovey/gophermart/internal/logger"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/accrual"
	"github.com/Makovey/gophermart/internal/transport/accrual/model"
)

type worker struct {
	orderRepo   service.OrderRepository
	balanceRepo service.BalancesRepository
	client      transport.Accrual
	ticker      *time.Ticker
	log         logger.Logger
}

func NewWorker(
	repo service.GophermartRepository,
	client transport.Accrual,
	log logger.Logger,
) service.Worker {
	return &worker{
		orderRepo:   repo,
		balanceRepo: repo,
		client:      client,
		ticker:      time.NewTicker(time.Second * 1),
		log:         log,
	}
}

func (w *worker) ProcessNewOrders() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orders := make(chan repoModel.Order, 5)
	w.registerNewGoods(ctx)
	w.runFetchingProcess(ctx, orders)
	w.runUpdatingProcess(ctx, orders)
}

func (w *worker) runFetchingProcess(
	ctx context.Context,
	orders chan<- repoModel.Order,
) {
	fn := "worker.runFetchingProcess"

	go func() {
		for {
			select {
			case <-w.ticker.C:
				err := w.orderRepo.FetchNewOrdersToChan(ctx, orders)
				if err != nil {
					w.log.Error(fmt.Sprintf("%s: failed to fetch new orders", fn))
				}
			case <-ctx.Done():
				close(orders)
				return
			}
		}
	}()
}

func (w *worker) registerNewGoods(ctx context.Context) {
	fn := "worker.registerNewGoods"

	if err := w.client.RegisterNewGoods(ctx); err != nil {
		w.log.Error(fmt.Sprintf("%s: failed to register new goods", fn))
	}
}

func (w *worker) runUpdatingProcess(ctx context.Context, orders <-chan repoModel.Order) {
	fn := "worker.runUpdatingProcess"

	for {
		select {
		case <-ctx.Done():
			w.log.Debug(fmt.Sprintf("%s: stopping updating process, context is closed", fn))
			return
		case order, ok := <-orders:
			if !ok {
				w.log.Debug(fmt.Sprintf("%s: stopping updating process, channel is closed", fn))
				return
			}
			w.client.RegisterNewOrder(ctx, order.OrderID)

			res, err := w.client.UpdateOrderStatus(ctx, order.OrderID)
			if err == nil {
				w.updateOrderInfo(ctx, res, order.OwnerUserID)
				continue
			}

			var manyReqErr *accrual.ManyRequestError
			switch {
			case errors.As(err, &manyReqErr):
				go func(order repoModel.Order, retryAfter time.Duration) {
					select {
					case <-ctx.Done():
						w.log.Info(fmt.Sprintf("%s: context cancelled before retrying order %s", fn, order.OrderID))
						return
					case <-time.After(retryAfter + time.Second):
						retryRes, retryErr := w.client.UpdateOrderStatus(ctx, order.OrderID)
						if retryErr != nil {
							w.log.Error(fmt.Sprintf("%s: retried methods is failed", fn), "error", retryErr)
						}
						w.updateOrderInfo(ctx, retryRes, order.OwnerUserID)
					}
				}(order, manyReqErr.RetryAfter)
			default:
				w.log.Error(fmt.Sprintf("%s: failed to update order status", fn), "error", err)
			}
		}
	}
}

func (w *worker) updateOrderInfo(ctx context.Context, status model.OrderStatus, userID string) {
	fn := "worker.updateOrderInfo"
	if ctx.Err() != nil {
		w.log.Error(fmt.Sprintf("%s: can't update order info", fn), "error", ctx.Err().Error())
		return
	}

	err := w.orderRepo.UpdateOrder(ctx, repoModel.OrderStatus{OrderID: status.OrderID, Status: status.Status, Accrual: status.Accrual})
	if err != nil {
		w.log.Error(fmt.Sprintf("%s: failed to update order info", fn), "error", err)
	}

	err = w.balanceRepo.IncreaseUsersBalance(ctx, userID, status.Accrual)
	if err != nil {
		w.log.Error(fmt.Sprintf("%s: failed to update users balance", fn), "error", err)
	}
}

func (w *worker) DownProcess() {
	w.ticker.Stop()
}
