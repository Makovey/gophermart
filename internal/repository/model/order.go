package model

import "github.com/shopspring/decimal"

type Status string

const (
	New        Status = "NEW"
	Processing Status = "PROCESSING"
	Invalid    Status = "INVALID"
	Processed  Status = "PROCESSED"
)

type Order struct {
	OrderID     string
	OwnerUserID string
	Status      Status
	Accrual     decimal.Decimal
}
