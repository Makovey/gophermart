package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/accrual"
)

type worker struct {
	repo   service.GophermartRepository
	client transport.Accrual
	ticker *time.Ticker
	log    logger.Logger
	quit   chan bool
}

func NewWorker(
	repo service.GophermartRepository,
	client transport.Accrual,
	log logger.Logger,
) service.Worker {
	return &worker{
		repo:   repo,
		client: client,
		ticker: time.NewTicker(time.Second),
		log:    log,
		quit:   make(chan bool),
	}
}

func (w *worker) ProcessNewOrders() {
	fn := "worker.ProcessNewOrders"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := w.client.RegisterNewGoods(ctx); err != nil {
		w.log.Error(fmt.Sprintf("%s: failed to register new goods", fn))
	}

	if err := w.client.RegisterNewOrder(ctx, "0018"); err != nil {
		w.log.Error(fmt.Sprintf("%s: failed to register new order", fn))
	}

	res, err := w.client.UpdateOrderStatus(ctx, "0018")
	if err != nil {
		var manyReqErr *accrual.ManyRequestError
		switch {
		case errors.As(err, &manyReqErr):
			time.AfterFunc(manyReqErr.RetryAfter+time.Second, func() {
				retrtyRes, retryErr := w.client.UpdateOrderStatus(ctx, "0018")
				if retryErr != nil {
					w.log.Error(fmt.Sprintf("%s: retried methods is failed", fn), "error", retryErr)
				}
				fmt.Println(retrtyRes)
				w.repo.UpdateOrder(ctx, model.OrderStatus{OrderID: res.OrderID, Status: res.Status, Accrual: res.Accrual})
			})
		default:
			//
		}
	}
	w.repo.UpdateOrder(ctx, model.OrderStatus{OrderID: res.OrderID, Status: res.Status, Accrual: res.Accrual})
	//orders := make(chan model.Order, 5)
	//ctx, _ := context.WithCancel(context.Background())
	//err := w.repo.FetchNewOrdersToChan(ctx, orders)
	//if err != nil {
	//	return
	//}
	//
	//w.client.SendOrder()
	//
	////go func() {
	////	for {
	////		select {
	////		case <-w.ticker.C:
	////			fmt.Println("processing new orders")
	////		case <-w.quit:
	////			stop()
	////			close(orders)
	////			return
	////		}
	////	}
	////}()
}

func (w *worker) DownProcess() {
	w.ticker.Stop()
	w.quit <- true
	close(w.quit)
}
