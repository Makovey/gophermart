package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Withdraw struct {
	OrderID   string
	Withdraw  decimal.Decimal
	CreatedAt time.Time
}
