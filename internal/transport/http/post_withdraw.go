package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport/http/model"
)

func (h handler) PostWithdraw(w http.ResponseWriter, r *http.Request) {
	fn := "http.PostWithdraw"

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, reloginAndTryAgain)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Info(fmt.Sprintf("%s: bad request received", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusBadRequest, badRequestError)
		return
	}

	var reqModel model.WithdrawRequest
	err = json.Unmarshal(body, &reqModel)
	if err != nil {
		h.log.Warn(fmt.Sprintf("%s: can't unmarshal request body", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
		return
	}

	err = h.balanceService.WithdrawBalance(r.Context(), userID, reqModel)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderConflict),
			errors.Is(err, service.ErrNotFound):
			h.writeResponseWithError(w, http.StatusUnprocessableEntity, "order not found or belong to another user")
		case errors.Is(err, service.ErrNotEnoughFounds):
			h.writeResponseWithError(w, http.StatusPaymentRequired, "not enough accrual on balance")
		default:
			h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
		}
	}

	w.WriteHeader(http.StatusOK)
}
