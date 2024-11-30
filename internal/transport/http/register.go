package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Makovey/gophermart/internal/repository"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

const (
	internalError     = "internal server error"
	badRequestError   = "bad request"
	userAlreadyExists = "user is already registered"
)

func (h handler) Register(w http.ResponseWriter, r *http.Request) {
	fn := "http.Register"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Info(fmt.Sprintf("%s: bad request received", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusBadRequest, badRequestError)
		return
	}

	var reqModel model.AuthRequest
	err = json.Unmarshal(body, &reqModel)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s: can't unmarshal request body", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
		return
	}

	header, err := h.service.RegisterUser(r.Context(), reqModel)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrLoginIsAlreadyExist):
			h.writeResponseWithError(w, http.StatusConflict, userAlreadyExists)
			return
		case errors.Is(err, service.ErrGeneratePass),
			errors.Is(err, repository.ErrPrepareStmt),
			errors.Is(err, repository.ErrPrepareStmt):
			h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
			return
		}
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", header))
	w.WriteHeader(http.StatusOK)
}
