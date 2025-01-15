package http

import (
	"fmt"
	"net/http"
)

func (h handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	fn := "http.GetBalance"

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, reloginAndTryAgain)
		return
	}

	balance, err := h.balanceService.GetUsersBalance(r.Context(), userID)
	if err != nil {
		h.log.Error(fmt.Sprintf("[%s] can't get users balance, userID - %s", fn, userID), "err", err.Error())
		h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
		return
	}

	h.writeResponse(w, http.StatusOK, balance)
}
