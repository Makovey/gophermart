package gophermart

import (
	"time"

	repo "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
)

func FromRepoToBalance(balance repo.Balance) model.BalanceResponse {
	return model.BalanceResponse{
		Current:   types.FloatDecimal(balance.Accrual),
		Withdrawn: types.FloatDecimal(balance.Withdrawn),
	}
}

func FromRepoToOrders(repoOrders []repo.Order) []model.OrderResponse {
	var orders []model.OrderResponse
	for _, repoOrder := range repoOrders {
		orders = append(orders, fromRepoToOrder(repoOrder))
	}

	return orders
}

func fromRepoToOrder(repoOrder repo.Order) model.OrderResponse {
	var accrual *types.FloatDecimal
	if repoOrder.Accrual != nil {
		val := types.FloatDecimal(*repoOrder.Accrual)
		accrual = &val
	}

	return model.OrderResponse{
		Number:     repoOrder.OrderID,
		Status:     repoOrder.Status,
		Accrual:    accrual,
		UploadedAt: repoOrder.CreatedAt.Format(time.RFC3339),
	}
}

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
