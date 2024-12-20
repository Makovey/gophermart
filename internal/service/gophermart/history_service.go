package gophermart

import (
	"context"

	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
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
	resp, err := h.repo.GetUsersHistory(ctx, userID)
	if err != nil {
		return nil, err
	}

	var models []model.WithdrawHistoryResponse
	for _, mod := range resp {
		withdraw := model.WithdrawHistoryResponse{
			Order:       mod.OrderID,
			Sum:         types.FloatDecimal(mod.Withdraw),
			ProcessedAt: mod.CreatedAt,
		}

		models = append(models, withdraw)
	}

	return models, nil
}
