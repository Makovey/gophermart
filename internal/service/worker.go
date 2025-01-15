package service

import "context"

type Worker interface {
	ProcessNewOrders(ctx context.Context)
	DownProcess()
}
