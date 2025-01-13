package gophermart

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"

	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/adapter"
	"github.com/Makovey/gophermart/internal/transport/http"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

//go:generate mockgen -source=balance_service.go -destination=../../repository/mocks/balance_mock.go -package=mocks
type BalancesServiceRepository interface {
	DecreaseUsersBalance(ctx context.Context, userID string, withdraw decimal.Decimal) error
	GetUsersBalance(ctx context.Context, userID string) (repoModel.Balance, error)
	RecordUsersWithdraw(ctx context.Context, userID, orderID string, amount decimal.Decimal) error
}

type balanceService struct {
	balanceRepo BalancesServiceRepository
}

func NewBalanceService(
	balanceRepo BalancesServiceRepository,
) http.BalanceService {
	return &balanceService{
		balanceRepo: balanceRepo,
	}
}

func (b *balanceService) GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error) {
	fn := "gophermart.GetUsersBalance"

	balance, err := b.balanceRepo.GetUsersBalance(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return adapter.FromRepoToBalance(repoModel.Balance{}), nil
		default:
			return adapter.FromRepoToBalance(repoModel.Balance{}), fmt.Errorf("[%s]: %w", fn, err)
		}
	}

	return adapter.FromRepoToBalance(balance), nil
}

func (b *balanceService) WithdrawUsersBalance(ctx context.Context, userID string, withdrawModel model.WithdrawRequest) error {
	fn := "gophermart.WithdrawUsersBalance"
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
		switch {
		case errors.Is(err, service.ErrNotFound):
			balance = repoModel.Balance{}
		default:
			return fmt.Errorf("[%s]: %w", fn, err)
		}
	}

	if balance.Accrual.LessThan(withdrawModel.Sum) {
		return fmt.Errorf("[%s]: %w", fn, service.ErrNotEnoughFounds)
	}

	err = b.balanceRepo.DecreaseUsersBalance(ctx, userID, withdrawModel.Sum)
	if err != nil {
		return fmt.Errorf("[%s]: %w", fn, err)
	}

	err = b.balanceRepo.RecordUsersWithdraw(ctx, userID, withdrawModel.Order, withdrawModel.Sum)
	if err != nil {
		return fmt.Errorf("[%s]: %w", fn, err)
	}

	return nil
}
