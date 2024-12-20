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
	quit        chan bool
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
		ticker:      time.NewTicker(time.Second * 30),
		log:         log,
		quit:        make(chan bool),
	}
}

func (w *worker) ProcessNewOrders() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	orders := make(chan repoModel.Order, 5)
	w.registerNewGoods(ctx)
	w.runFetchingProcess(ctx, cancel, orders)
	w.runUpdatingProcess(ctx, orders)
}

func (w *worker) runFetchingProcess(
	ctx context.Context,
	cancel context.CancelFunc,
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
			case <-w.quit:
				cancel()
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

	for order := range orders {
		if err := w.client.RegisterNewOrder(ctx, order.OrderID); err != nil {
			w.log.Error(fmt.Sprintf("%s: failed to register new order", fn), "error", err.Error())
		}

		res, err := w.client.UpdateOrderStatus(ctx, order.OrderID)
		if err == nil {
			w.updateOrderInfo(ctx, res, order.OwnerUserID)
			continue
		}

		var manyReqErr *accrual.ManyRequestError
		switch {
		case errors.As(err, &manyReqErr):
			time.AfterFunc(manyReqErr.RetryAfter+time.Second, func() {
				retryRes, retryErr := w.client.UpdateOrderStatus(ctx, order.OrderID)
				if retryErr != nil {
					w.log.Error(fmt.Sprintf("%s: retried methods is failed", fn), "error", retryErr)
				}
				w.updateOrderInfo(ctx, retryRes, order.OwnerUserID)
			})
		default:
			w.log.Error(fmt.Sprintf("%s: failed to update order status", fn), "error", err)
		}
	}
}

func (w *worker) updateOrderInfo(ctx context.Context, status model.OrderStatus, userID string) {
	fn := "worker.updateOrderInfo"

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
	w.quit <- true
	close(w.quit)
}
