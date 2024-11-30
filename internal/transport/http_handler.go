package transport

import "net/http"

type HTTPHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
}
