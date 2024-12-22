package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

const (
	internalError          = "internal server error"
	badRequestError        = "bad request"
	loginOrPasswordIsEmpty = "login or password is empty, or greater than 30 symbols"
)

func (h handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	fn := "http.RegisterUser"

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
		h.log.Error(fmt.Sprintf("%s: login or password is too long", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusBadRequest, loginOrPasswordIsEmpty)
		return
	}

	token, err := h.userService.RegisterNewUser(r.Context(), reqModel)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrLoginIsAlreadyExist):
			h.writeResponseWithError(w, http.StatusConflict, "user is already registered")
			return
		case errors.Is(err, service.ErrGeneratePass),
			errors.Is(err, service.ErrExecStmt):
			h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
			return
		}
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
}
