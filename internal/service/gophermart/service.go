package gophermart

import (
	"context"
	"github.com/Makovey/gophermart/internal/transport/http"

	"github.com/Makovey/gophermart/internal/transport/http/model"
)

type serv struct {
	userServ    http.UserService
	orderServ   http.OrderService
	balanceServ http.BalanceService
	historyServ http.HistoryService
}

func NewGophermartService(
	userServ http.UserService,
	orderServ http.OrderService,
	balanceServ http.BalanceService,
	historyServ http.HistoryService,
) http.GophermartService {
	return &serv{
		userServ:    userServ,
		orderServ:   orderServ,
		balanceServ: balanceServ,
		historyServ: historyServ,
	}
}

func (s *serv) RegisterNewUser(ctx context.Context, request model.AuthRequest) (string, error) {
	return s.userServ.RegisterNewUser(ctx, request)
}

func (s *serv) LoginUser(ctx context.Context, request model.AuthRequest) (string, error) {
	return s.userServ.LoginUser(ctx, request)
}

func (s *serv) ValidateOrderID(orderID string) bool {
	return s.orderServ.ValidateOrderID(orderID)
}

func (s *serv) ProcessNewOrder(ctx context.Context, userID, orderID string) error {
	return s.orderServ.ProcessNewOrder(ctx, userID, orderID)
}

func (s *serv) GetOrders(ctx context.Context, userID string) ([]model.OrderResponse, error) {
	return s.orderServ.GetOrders(ctx, userID)
}

func (s *serv) GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error) {
	return s.balanceServ.GetUsersBalance(ctx, userID)
}

func (s *serv) WithdrawUsersBalance(ctx context.Context, userID string, withdrawModel model.WithdrawRequest) error {
	return s.balanceServ.WithdrawUsersBalance(ctx, userID, withdrawModel)
}

func (s *serv) GetUsersWithdrawHistory(ctx context.Context, userID string) ([]model.WithdrawHistoryResponse, error) {
	return s.historyServ.GetUsersWithdrawHistory(ctx, userID)
}
