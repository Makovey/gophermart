package model

import (
	"github.com/shopspring/decimal"

	"github.com/Makovey/gophermart/internal/repository/model"
)

type OrderStatus struct {
	OrderID string          `json:"order"`
	Status  model.Status    `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
}
