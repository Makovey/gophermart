package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/accrual"
	"github.com/Makovey/gophermart/internal/transport/accrual/model"
)

//go:generate mockgen -source=accrual.go -destination=../../repository/mocks/worker_mock.go -package=mocks
type WorkerRepository interface {
	FetchNewOrdersToChan(ctx context.Context, ordersCh chan<- repoModel.Order) error
	UpdateOrder(ctx context.Context, status repoModel.OrderStatus) error
	IncreaseUsersBalance(ctx context.Context, userID string, reward decimal.Decimal) error
}

type worker struct {
	repo   WorkerRepository
	client transport.Accrual
	ticker *time.Ticker
	log    logger.Logger
	cfg    config.Config
	wg     *sync.WaitGroup
}

func NewWorker(
	repo WorkerRepository,
	client transport.Accrual,
	cfg config.Config,
	log logger.Logger,
) service.Worker {
	return &worker{
		repo:   repo,
		client: client,
		ticker: time.NewTicker(cfg.TickerTimer()),
		log:    log,
		wg:     &sync.WaitGroup{},
	}
}

func (w *worker) ProcessNewOrders(ctx context.Context) {
	orders := make(chan repoModel.Order, 5)
	w.registerNewGoods(ctx)
	go w.runFetchingProcess(ctx, orders)
	go w.runUpdatingProcess(ctx, orders)
	w.wg.Wait()
}

func (w *worker) runFetchingProcess(
	ctx context.Context,
	orders chan<- repoModel.Order,
) {
	fn := "worker.runFetchingProcess"
	w.wg.Add(1)
	defer w.wg.Done()

	for {
		select {
		case <-w.ticker.C:
			err := w.repo.FetchNewOrdersToChan(ctx, orders)
			if err != nil {
				w.log.Error(fmt.Sprintf("[%s] failed to fetch new orders", fn))
			}
		case <-ctx.Done():
			w.log.Debug(fmt.Sprintf("[%s] stopping fetching process, context is closed", fn))
			close(orders)
			return
		}
	}
}

func (w *worker) registerNewGoods(ctx context.Context) {
	fn := "worker.registerNewGoods"

	if err := w.client.RegisterNewGoods(ctx); err != nil {
		w.log.Error(fmt.Sprintf("[%s] failed to register new goods", fn))
	}
}

func (w *worker) runUpdatingProcess(ctx context.Context, orders <-chan repoModel.Order) {
	fn := "worker.runUpdatingProcess"
	w.wg.Add(1)
	defer w.wg.Done()

	for order := range orders {
		if err := w.client.RegisterNewOrder(ctx, order.OrderID); err != nil {
			if !errors.Is(err, service.ErrOrderAlreadyRegistered) {
				w.log.Error(fmt.Sprintf("[%s] failed to register new order", fn), "error", err.Error())
				continue
			}
		}

		status, err := w.client.UpdateOrderStatus(ctx, order.OrderID)
		if err == nil {
			w.updateOrderInfo(ctx, status, order.OwnerUserID)
			continue
		}

		var manyReqErr *accrual.ManyRequestError
		switch {
		case errors.As(err, &manyReqErr):
			go func(order repoModel.Order, retryAfter time.Duration) {
				select {
				case <-ctx.Done():
					w.log.Info(fmt.Sprintf("[%s] context cancelled before retrying order %s", fn, order.OrderID))
					return
				case <-time.After(retryAfter + time.Second):
					retryRes, retryErr := w.client.UpdateOrderStatus(ctx, order.OrderID)
					if retryErr != nil {
						w.log.Error(fmt.Sprintf("[%s] retried methods is failed", fn), "error", retryErr)
					}
					w.updateOrderInfo(ctx, retryRes, order.OwnerUserID)
				}
			}(order, manyReqErr.RetryAfter)
		default:
			w.log.Error(fmt.Sprintf("[%s] failed to update order status", fn), "error", err)
		}
	}
}

func (w *worker) updateOrderInfo(ctx context.Context, status model.OrderStatus, userID string) {
	fn := "worker.updateOrderInfo"
	if ctx.Err() != nil {
		w.log.Error(fmt.Sprintf("[%s] can't update order info", fn), "error", ctx.Err().Error())
		return
	}

	err := w.repo.UpdateOrder(ctx, repoModel.OrderStatus{OrderID: status.OrderID, Status: status.Status, Accrual: status.Accrual})
	if err != nil {
		w.log.Error(fmt.Sprintf("[%s] failed to update order info", fn), "error", err)
	}

	err = w.repo.IncreaseUsersBalance(ctx, userID, status.Accrual)
	if err != nil {
		w.log.Error(fmt.Sprintf("[%s] failed to update users balance", fn), "error", err)
	}
}

func (w *worker) DownProcess() {
	w.ticker.Stop()
}
