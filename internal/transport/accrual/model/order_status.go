package model

import (
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/shopspring/decimal"
)

type OrderStatus struct {
	OrderID string          `json:"order"`
	Status  model.Status    `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
}
