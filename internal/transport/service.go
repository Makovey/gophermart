package transport

import (
	"context"

	"github.com/Makovey/gophermart/internal/transport/http/model"
)

//go:generate mockgen -source=service.go -destination=../service/mocks/service_mock.go -package=mocks
type GophermartService interface {
	RegisterUser(ctx context.Context, request model.AuthRequest) (string, error)
}
