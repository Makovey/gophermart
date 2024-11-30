package service

import (
	"context"

	"github.com/Makovey/gophermart/internal/repository/model"
)

type GophermartRepository interface {
	RegisterNewUser(ctx context.Context, user model.RegisterUser) error
}
