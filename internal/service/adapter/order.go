package adapter

import (
	"time"

	repo "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/internal/types"
)

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
