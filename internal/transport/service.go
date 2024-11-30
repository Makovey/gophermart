package transport

import (
	"context"

	"github.com/Makovey/gophermart/internal/transport/http/model"
)

type GophermartService interface {
	RegisterUser(ctx context.Context, request model.AuthRequest) (string, error)
}
