package gophermart

import (
	"context"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
	"github.com/Makovey/gophermart/pkg/jwt"
)

type balanceService struct {
	balanceRepo service.BalancesRepository
	orderRepo   service.OrderRepository
	historyRepo service.HistoryRepository
	log         logger.Logger
	jwt         *jwt.JWT
}

func newBalanceService(
	repo service.GophermartRepository,
	log logger.Logger,
	jwt *jwt.JWT,
) transport.BalanceService {
	return &balanceService{
		balanceRepo: repo,
		orderRepo:   repo,
		historyRepo: repo,
		log:         log,
		jwt:         jwt,
	}
}

func (b *balanceService) GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error) {
	balance, err := b.balanceRepo.GetUsersBalance(ctx, userID)
	if err != nil {
		return model.BalanceResponse{}, err
	}

	return model.BalanceResponse{
		Current:   types.FloatDecimal(balance.Accrual),
		Withdrawn: types.FloatDecimal(balance.Withdrawn),
	}, nil
}

func (b *balanceService) WithdrawBalance(ctx context.Context, userID string, withdrawModel model.WithdrawRequest) error {
	order, err := b.orderRepo.GetOrderByID(ctx, withdrawModel.Order)
	if err != nil {
		return err
	}

	if order.OwnerUserID != userID {
		return service.ErrOrderConflict
	}

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
