package http

import (
	"fmt"
	"net/http"
)

func (h handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	fn := "http.GetOrders"

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		h.writeResponseWithError(w, http.StatusBadRequest, reloginAndTryAgain)
		return
	}

	orders, err := h.service.GetOrders(r.Context(), userID)
	if err != nil {
		h.log.Error(fmt.Sprintf("%s: can't get orders", fn), "err", err.Error())
		h.writeResponseWithError(w, http.StatusInternalServerError, internalError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	h.writeResponse(w, http.StatusOK, orders)
}
