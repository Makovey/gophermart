package model

import "github.com/shopspring/decimal"

type Balance struct {
	Accrual   decimal.Decimal
	Withdrawn decimal.Decimal
}
