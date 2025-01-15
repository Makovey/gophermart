package model

import (
	"github.com/shopspring/decimal"
)

type WithdrawRequest struct {
	Order string          `json:"order" validate:"required"`
	Sum   decimal.Decimal `json:"sum" validate:"required"`
}
