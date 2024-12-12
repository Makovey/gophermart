package gophermart

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/Makovey/gophermart/internal/logger"
	repoModel "github.com/Makovey/gophermart/internal/repository/model"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/luhn"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"github.com/Makovey/gophermart/pkg/jwt"
)

var (
	UserIDLength = 10
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
		UserID:       uuid.NewString()[:UserIDLength],
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

func (s serv) ValidateOrderID(orderID string) bool {
	orderInt, err := strconv.Atoi(orderID)
	if err != nil {
		return false
	}

	return luhn.IsValid(orderInt)
}

func (s serv) ProcessNewOrder(ctx context.Context, userID, orderID string) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			return s.repo.PostNewOrder(ctx, orderID, userID)
		default:
			return err
		}
	}

	if order.OwnerUserID != userID {
		return service.ErrOrderConflict
	}

	return service.ErrOrderAlreadyPosted
}

func (s serv) GetOrders(ctx context.Context, userID string) ([]model.Order, error) {
	repoOrders, err := s.repo.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	var models []model.Order
	for _, repOrder := range repoOrders {
		var accrual *float64
		if repOrder.Accrual != nil {
			float := repOrder.Accrual.Round(2).InexactFloat64()
			accrual = &float
		}

		order := model.Order{
			Number:     repOrder.OrderID,
			Status:     repOrder.Status,
			Accrual:    accrual,
			UploadedAt: repOrder.CreatedAt.Format(time.RFC3339),
		}

		models = append(models, order)
	}

	return models, nil
}
