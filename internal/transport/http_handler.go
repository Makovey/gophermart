package transport

import "net/http"

type HTTPHandler interface {
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)

	PostOrder(w http.ResponseWriter, r *http.Request)
	GetOrders(w http.ResponseWriter, r *http.Request)

	GetBalance(w http.ResponseWriter, r *http.Request)
	PostWithdraw(w http.ResponseWriter, r *http.Request)
	GetWithdrawsHistory(w http.ResponseWriter, r *http.Request)
}
