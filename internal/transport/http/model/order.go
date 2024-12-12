package model

import (
	"github.com/Makovey/gophermart/internal/repository/model"
)

type Order struct {
	Number     string       `json:"number"`
	Status     model.Status `json:"status"`
	Accrual    *float64     `json:"accrual,omitempty"`
	UploadedAt string       `json:"uploaded_at"`
}
