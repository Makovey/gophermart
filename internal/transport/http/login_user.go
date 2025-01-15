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

func (h handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	fn := "http.LoginUser"

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Info(fmt.Sprintf("[%s] bad request received", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusBadRequest, badRequestError)
		return
	}

	var reqModel model.AuthRequest
	err = json.Unmarshal(body, &reqModel)
	if err != nil {
		h.log.Warn(fmt.Sprintf("[%s] can't unmarshal request body", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
		return
	}

	validate := validator.New()
	if err = validate.Struct(reqModel); err != nil {
		h.log.Info(fmt.Sprintf("[%s] login or password is too long", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusBadRequest, loginOrPasswordIsEmpty)
		return
	}

	token, err := h.userService.LoginUser(r.Context(), reqModel)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrExecStmt):
			h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
			return
		case errors.Is(err, service.ErrNotFound),
			errors.Is(err, service.ErrPasswordDoesntMatch):
			h.log.Info(fmt.Sprintf("[%s] user can't login with login: %s", fn, reqModel.Login), "error", err.Error())
			h.writeResponseWithError(w, http.StatusUnauthorized, "login or password is incorrect")
			return
		}
	}

	w.Header().Set("Authorization", token)
	w.WriteHeader(http.StatusOK)
}
