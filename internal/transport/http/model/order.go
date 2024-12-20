package model

import (
	"github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/types"
)

type OrderResponse struct {
	Number     string              `json:"number"`
	Status     model.Status        `json:"status"`
	Accrual    *types.FloatDecimal `json:"accrual,omitempty"`
	UploadedAt string              `json:"uploaded_at"`
}
