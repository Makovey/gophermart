package gophermart

import (
	"context"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/pkg/jwt"
)

type serv struct {
	userServ    transport.UserService
	orderServ   transport.OrderService
	balanceServ transport.BalanceService
	historyServ transport.HistoryService
}

func NewGophermartService(
	repo service.GophermartRepository,
	logger logger.Logger,
	jwt *jwt.JWT,
) transport.GophermartService {
	return &serv{
		userServ:    newUserService(repo, logger, jwt),
		orderServ:   newOrderService(repo),
		balanceServ: newBalanceService(repo),
		historyServ: newHistoryService(repo),
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

func (s *serv) WithdrawBalance(ctx context.Context, userID string, withdrawModel model.WithdrawRequest) error {
	return s.balanceServ.WithdrawBalance(ctx, userID, withdrawModel)
}

func (s *serv) GetUsersWithdrawHistory(ctx context.Context, userID string) ([]model.WithdrawHistoryResponse, error) {
	return s.historyServ.GetUsersWithdrawHistory(ctx, userID)
}
