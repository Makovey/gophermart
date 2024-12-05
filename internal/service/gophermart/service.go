package gophermart

import (
	"context"
	"fmt"
	"github.com/Makovey/gophermart/internal/middleware/utils"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Makovey/gophermart/internal/logger"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

const (
	tokenExp = time.Hour * 24
)

type serv struct {
	repo   service.GophermartRepository
	logger logger.Logger
}

func (s serv) RegisterUser(ctx context.Context, request model.AuthRequest) (string, error) {
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

	jwtToken, err := s.buildNewJWT(user.UserID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s serv) buildNewJWT(userID string) (string, error) {
	fn := "gophermart.buildNewJWT:"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, utils.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte("gophermart_key"))
	if err != nil {
		s.logger.Warn(fmt.Sprintf("%s: can't sign token", fn), "error", err.Error())
		return "", err
	}

	return tokenString, nil
}

func NewGophermartService(
	repo service.GophermartRepository,
	logger logger.Logger,
) transport.GophermartService {
	return &serv{repo, logger}
}
