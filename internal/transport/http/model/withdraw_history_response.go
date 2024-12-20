package model

import (
	"time"

	"github.com/Makovey/gophermart/internal/types"
)

type WithdrawHistoryResponse struct {
	Order       string             `json:"order"`
	Sum         types.FloatDecimal `json:"sum"`
	ProcessedAt time.Time          `json:"processed_at"`
}
