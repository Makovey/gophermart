package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/Makovey/gophermart/internal/middleware"
)

func (a *App) initRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.NewCompressor().Compress)

	r.Post("/api/user/register", a.deps.Handler().Register)
	r.Post("/api/user/login", a.deps.Handler().Login)

	r.Group(func(r chi.Router) {
		authMiddleware := middleware.NewAuth(a.deps.JWT(), a.deps.Logger())

		r.Use(authMiddleware.Authenticate)
		r.Get("/api/user/orders", a.deps.Handler().GetOrders)
		r.Post("/api/user/orders", a.deps.Handler().PostOrder)

		r.Get("/api/user/balance", func(w http.ResponseWriter, r *http.Request) {})
		r.Post("/api/user/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {})

		r.Get("/api/user/withdrawals", func(w http.ResponseWriter, r *http.Request) {})
	})

	return r
}
