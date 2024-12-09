package gophermart

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"

	"github.com/Makovey/gophermart/internal/logger"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/pkg/jwt"
)

type serv struct {
	repo   service.GophermartRepository
	logger logger.Logger
	jwt    *jwt.JWT
}

func NewGophermartService(
	repo service.GophermartRepository,
	logger logger.Logger,
	jwt *jwt.JWT,
) transport.GophermartService {
	return &serv{repo: repo, logger: logger, jwt: jwt}
}

func (s serv) RegisterNewUser(ctx context.Context, request model.AuthRequest) (string, error) {
	fn := "gophermart.RegisterUser"

	pass, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: failed to generate password hash", fn), "error", err.Error())
		return "", service.ErrGeneratePass
	}

	user := repoModel.RegisterUser{
		UserID:       uuid.NewString()[:8],
		Login:        request.Login,
		PasswordHash: string(pass),
	}

	if err = s.repo.RegisterNewUser(ctx, user); err != nil {
		return "", err
	}

	jwtToken, err := s.jwt.BuildNewJWT(user.UserID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s serv) LoginUser(ctx context.Context, request model.AuthRequest) (string, error) {
	user, err := s.repo.LoginUser(ctx, request.Login)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return "", service.ErrPasswordDoesntMatch
	}

	jwtToken, err := s.jwt.BuildNewJWT(user.UserID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
