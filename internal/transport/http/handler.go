package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/service/gophermart"
	"github.com/Makovey/gophermart/internal/transport"
)

type handler struct {
	log logger.Logger

	service transport.GophermartService
}

func NewHTTPHandler(
	log logger.Logger,
	service transport.GophermartService,
) transport.HTTPHandler {
	return &handler{log: log, service: service}
}

func (h handler) writeResponseWithError(w http.ResponseWriter, statusCode int, message string) {
	fn := "http.writeResponseWithError"

	errResp := map[string]string{"error": message}
	err := writeJSON(w, statusCode, errResp)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s: failed to write response:", fn), "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h handler) writeResponse(w http.ResponseWriter, statusCode int, body any) {
	fn := "http.writeResponse"

	err := writeJSON(w, statusCode, body)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s: failed to write response:", fn), "error", err.Error())
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
	if utf8.RuneCountInString(key) != gophermart.UserIDLength {
		return "", errors.New("invalid user id")
	}

	return key, nil
}
