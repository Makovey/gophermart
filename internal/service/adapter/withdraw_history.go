package adapter

import (
	repo "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
)

func FromRepoToHistoryWithdraws(repoWithdraws []repo.Withdraw) []model.WithdrawHistoryResponse {
	var withdraws []model.WithdrawHistoryResponse
	for _, repoWithdraw := range repoWithdraws {
		withdraws = append(withdraws, fromRepoToHistoryWithdraw(repoWithdraw))
	}

	return withdraws
}

func fromRepoToHistoryWithdraw(repoWithdraw repo.Withdraw) model.WithdrawHistoryResponse {
	return model.WithdrawHistoryResponse{
		Order:       repoWithdraw.OrderID,
		Sum:         types.FloatDecimal(repoWithdraw.Withdraw),
		ProcessedAt: repoWithdraw.CreatedAt,
	}
}
