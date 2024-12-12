package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Makovey/gophermart/internal/service"
)

const (
	reloginAndTryAgain = "relogin and try again"
	orderIDIsInvalid   = "order id is invalid"
)

func (h handler) PostOrder(w http.ResponseWriter, r *http.Request) {
	fn := "http.PostOrder"

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

	if len(body) == 0 {
		h.log.Info(fmt.Sprintf("%s: tryied to post order with empty body", fn), "userID", userID)
		h.writeResponseWithError(w, http.StatusBadRequest, badRequestError)
		return
	}

	if isValid := h.service.ValidateOrderID(string(body)); !isValid {
		h.log.Info(fmt.Sprintf("%s: invalid order id received", fn), "orderID", string(body))
		h.writeResponseWithError(w, http.StatusUnprocessableEntity, orderIDIsInvalid)
		return
	}

	err = h.service.ProcessNewOrder(r.Context(), userID, string(body))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderConflict):
			h.log.Info(fmt.Sprintf("%s: order already belong to another user", fn), "orderID", string(body))
			h.writeResponseWithError(w, http.StatusConflict, err.Error())
			return
		case errors.Is(err, service.ErrOrderAlreadyPosted):
			h.log.Info(fmt.Sprintf("%s: order already posted", fn), "orderID", string(body))
			w.WriteHeader(http.StatusOK)
			return
		default:
			h.log.Info(fmt.Sprintf("%s: error processing order request", fn), "error", err.Error())
			h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}
