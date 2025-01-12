package gophermart

import (
	"context"

	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/adapter"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

type historyService struct {
	repo service.HistoryRepository
}

func newHistoryService(
	repo service.HistoryRepository,
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
