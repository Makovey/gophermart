package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/Makovey/gophermart/internal/repository"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

const (
	internalError          = "internal server error"
	badRequestError        = "bad request"
	userAlreadyExists      = "user is already registered"
	loginOrPasswordIsEmpty = "login or password is empty, or greater than 30 symbols"
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

	validate := validator.New()
	if err = validate.Struct(reqModel); err != nil {
		h.log.Error(fmt.Sprintf("%s: login or password is empty or greater than 30 symbols", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusBadRequest, loginOrPasswordIsEmpty)
		return
	}

	header, err := h.service.RegisterUser(r.Context(), reqModel)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrLoginIsAlreadyExist):
			h.writeResponseWithError(w, http.StatusConflict, userAlreadyExists)
			return
		case errors.Is(err, service.ErrGeneratePass),
			errors.Is(err, repository.ErrExecStmt):
			h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
			return
		}
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", header))
	w.WriteHeader(http.StatusOK)
}
