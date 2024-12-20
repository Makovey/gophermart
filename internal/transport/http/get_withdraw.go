package http

import (
	"fmt"
	"net/http"
)

func (h handler) GetWithdrawsHistory(w http.ResponseWriter, r *http.Request) {
	fn := "http.GetWithdrawsHistory"

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, reloginAndTryAgain)
		return
	}

	withdraws, err := h.historyService.GetUsersWithdrawHistory(r.Context(), userID)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s: failed to get user withdraws", fn), "error", err.Error())
		h.writeResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(withdraws) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h.writeResponse(w, http.StatusOK, withdraws)
}
