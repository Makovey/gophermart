package gophermart

import (
	"context"
	"fmt"
	"github.com/Makovey/gophermart/internal/transport/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/pkg/jwt"
)

//go:generate mockgen -source=user_service.go -destination=../../repository/mocks/user_mock.go -package=mocks
type UserServiceRepository interface {
	RegisterNewUser(ctx context.Context, user repoModel.RegisterUser) error
	LoginUser(ctx context.Context, login string) (repoModel.RegisterUser, error)
}

type userService struct {
	repo UserServiceRepository
	jwt  *jwt.JWT
}

func NewUserService(
	repo UserServiceRepository,
	jwt *jwt.JWT,
) http.UserService {
	return &userService{
		repo: repo,
		jwt:  jwt,
	}
}

func (u *userService) RegisterNewUser(ctx context.Context, request model.AuthRequest) (string, error) {
	fn := "gophermart.RegisterUser"

	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("[%s]: %w", fn, service.ErrGeneratePass)
	}

	user := repoModel.RegisterUser{
		UserID:       uuid.NewString(),
		Login:        request.Login,
		PasswordHash: string(pass),
	}

	if err = u.repo.RegisterNewUser(ctx, user); err != nil {
		return "", fmt.Errorf("[%s]: %w", fn, err)
	}

	jwtToken, err := u.jwt.BuildNewJWT(user.UserID)
	if err != nil {
		return "", fmt.Errorf("[%s]: %w", fn, err)
	}

	return jwtToken, nil
}

func (u *userService) LoginUser(ctx context.Context, request model.AuthRequest) (string, error) {
	fn := "gophermart.LoginUser"

	user, err := u.repo.LoginUser(ctx, request.Login)
	if err != nil {
		return "", fmt.Errorf("[%s]: %w", fn, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return "", fmt.Errorf("[%s]: %w", fn, service.ErrPasswordDoesntMatch)
	}

	jwtToken, err := u.jwt.BuildNewJWT(user.UserID)
	if err != nil {
		return "", fmt.Errorf("[%s]: %w", fn, err)
	}

	return jwtToken, nil
}
