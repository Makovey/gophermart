package worker

import (
	"context"
	"time"

	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
)

type worker struct {
	repo   service.GophermartRepository
	client transport.Accrual
	ticker *time.Ticker
	quit   chan bool
}

func NewWorker(
	repo service.GophermartRepository,
	client transport.Accrual,
) service.Worker {
	return &worker{
		repo:   repo,
		client: client,
		ticker: time.NewTicker(time.Second),
		quit:   make(chan bool),
	}
}

func (w *worker) ProcessNewOrders() {
	w.client.RegisterNewGoods(context.Background())
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
