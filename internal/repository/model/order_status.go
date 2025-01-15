package model

import "github.com/shopspring/decimal"

type OrderStatus struct {
	OrderID string
	Status  Status
	Accrual decimal.Decimal
}
