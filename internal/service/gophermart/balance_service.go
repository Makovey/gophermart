package gophermart

import (
	"context"

	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/adapter"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

type balanceService struct {
	balanceRepo service.BalancesRepository
	orderRepo   service.OrderRepository
	historyRepo service.HistoryRepository
}

func newBalanceService(
	repo service.GophermartRepository,
) transport.BalanceService {
	return &balanceService{
		balanceRepo: repo,
		orderRepo:   repo,
		historyRepo: repo,
	}
}

func (b *balanceService) GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error) {
	balance, err := b.balanceRepo.GetUsersBalance(ctx, userID)
	if err != nil {
		return model.BalanceResponse{}, err
	}

	return adapter.FromRepoToBalance(balance), nil
}

func (b *balanceService) WithdrawUsersBalance(ctx context.Context, userID string, withdrawModel model.WithdrawRequest) error {
	// выключено для тестов
	// но, прежде чем списывать баллы, нужно убедиться, что заказ существует и пренадлежит именно этому пользователю
	//order, err := b.orderRepo.GetOrderByID(ctx, withdrawModel.Order)
	//if err != nil {
	//	return err
	//}
	//
	//if order.OwnerUserID != userID {
	//	return service.ErrOrderConflict
	//}

	balance, err := b.balanceRepo.GetUsersBalance(ctx, userID)
	if err != nil {
		return err
	}

	if balance.Accrual.LessThan(withdrawModel.Sum) {
		return service.ErrNotEnoughFounds
	}

	err = b.balanceRepo.DecreaseUsersBalance(ctx, userID, withdrawModel.Sum)
	if err != nil {
		return err
	}

	err = b.historyRepo.RecordUsersWithdraw(ctx, userID, withdrawModel.Order, withdrawModel.Sum)
	if err != nil {
		return err
	}

	return nil
}
