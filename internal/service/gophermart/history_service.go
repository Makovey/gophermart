package gophermart

import (
	"context"

	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service/adapter"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

//go:generate mockgen -source=history_service.go -destination=../../repository/mocks/history_mock.go -package=mocks
type HistoryServiceRepository interface {
	GetUsersHistory(ctx context.Context, userID string) ([]repoModel.Withdraw, error)
}

type historyService struct {
	repo HistoryServiceRepository
}

func NewHistoryService(
	repo HistoryServiceRepository,
) transport.HistoryService {
	return &historyService{
		repo: repo,
	}
}

func (h *historyService) GetUsersWithdrawHistory(ctx context.Context, userID string) ([]model.WithdrawHistoryResponse, error) {
	history, err := h.repo.GetUsersHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	return adapter.FromRepoToHistoryWithdraws(history), nil
}
