package adapter

import (
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
