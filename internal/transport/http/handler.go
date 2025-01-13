package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Makovey/gophermart/internal/transport/http/model"
	"net/http"
	"unicode/utf8"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/transport"
)

//go:generate mockgen -source=handler.go -destination=../../service/mocks/service_mock.go -package=mocks
type GophermartService interface {
	UserService
	OrderService
	BalanceService
	HistoryService
}

type UserService interface {
	RegisterNewUser(ctx context.Context, request model.AuthRequest) (string, error)
	LoginUser(ctx context.Context, request model.AuthRequest) (string, error)
}

type OrderService interface {
	ValidateOrderID(orderID string) bool
	ProcessNewOrder(ctx context.Context, userID, orderID string) error
	GetOrders(ctx context.Context, userID string) ([]model.OrderResponse, error)
}

type BalanceService interface {
	GetUsersBalance(ctx context.Context, userID string) (model.BalanceResponse, error)
	WithdrawUsersBalance(ctx context.Context, userID string, withdraw model.WithdrawRequest) error
}

type HistoryService interface {
	GetUsersWithdrawHistory(ctx context.Context, userID string) ([]model.WithdrawHistoryResponse, error)
}

type handler struct {
	log logger.Logger

	userService    UserService
	orderService   OrderService
	balanceService BalanceService
	historyService HistoryService
}

func NewHTTPHandler(
	log logger.Logger,
	service GophermartService,
) transport.HTTPHandler {
	return &handler{
		log:            log,
		userService:    service,
		orderService:   service,
		balanceService: service,
		historyService: service,
	}
}

const (
	UserIDLength = 36
)

func (h handler) writeResponseWithError(w http.ResponseWriter, statusCode int, message string) {
	fn := "http.writeResponseWithError"

	errResp := map[string]string{"error": message}
	err := writeJSON(w, statusCode, errResp)
	if err != nil {
		h.log.Error(fmt.Sprintf("[%s] failed to write response:", fn), "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h handler) writeResponse(w http.ResponseWriter, statusCode int, body any) {
	fn := "http.writeResponse"

	err := writeJSON(w, statusCode, body)
	if err != nil {
		h.log.Error(fmt.Sprintf("[%s] failed to write response:", fn), "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func getUserIDFromContext(ctx context.Context) (string, error) {
	if ctx.Value(middleware.CtxUserIDKey) == nil {
		return "", errors.New("user id not found in context")
	}

	key := ctx.Value(middleware.CtxUserIDKey).(string)
	if utf8.RuneCountInString(key) != UserIDLength {
		return "", errors.New("invalid user id")
	}

	return key, nil
}
