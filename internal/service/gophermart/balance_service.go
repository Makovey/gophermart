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
	repo service.BalancesRepository
	log  logger.Logger
	jwt  *jwt.JWT
}

func newBalanceService(
	repo service.BalancesRepository,
	log logger.Logger,
	jwt *jwt.JWT,
) transport.BalanceService {
	return &balanceService{
		repo: repo,
		log:  log,
		jwt:  jwt,
	}
}

func (b *balanceService) GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error) {
	balance, err := b.repo.GetUsersBalance(ctx, userID)
	if err != nil {
		return model.BalanceResponse{}, err
	}

	return model.BalanceResponse{
		Current:   types.FloatDecimal(balance.Accrual),
		Withdrawn: types.FloatDecimal(balance.Withdrawn),
	}, nil
}
