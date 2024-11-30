package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/transport"
)

type handler struct {
	log logger.Logger

	service transport.GophermartService
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

func NewHTTPHandler(
	log logger.Logger,
	service transport.GophermartService,
) transport.HTTPHandler {
	return &handler{log: log, service: service}
}
