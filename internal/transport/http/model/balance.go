package model

import (
	"github.com/Makovey/gophermart/internal/types"
)

type BalanceResponse struct {
	Current   types.FloatDecimal `json:"current"`
	Withdrawn types.FloatDecimal `json:"withdrawn"`
}
