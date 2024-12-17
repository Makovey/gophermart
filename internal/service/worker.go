package service

type Worker interface {
	ProcessNewOrders()
	DownProcess()
}
