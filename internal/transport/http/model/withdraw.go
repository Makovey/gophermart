package model

import (
	"github.com/shopspring/decimal"
)

type WithdrawRequest struct {
	Order string          `json:"order"`
	Sum   decimal.Decimal `json:"sum"`
}
