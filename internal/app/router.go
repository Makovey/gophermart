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

	r.Post("/api/user/register", a.handler.RegisterUser)
	r.Post("/api/user/login", a.handler.LoginUser)

	r.Group(func(r chi.Router) {
		r.Use(a.authMiddleware.Authenticate)
		r.Get("/api/user/orders", a.handler.GetOrders)
		r.Post("/api/user/orders", a.handler.PostOrder)

		r.Get("/api/user/balance", a.handler.GetBalance)
		r.Post("/api/user/balance/withdraw", a.handler.PostWithdraw)

		r.Get("/api/user/withdrawals", a.handler.GetWithdrawsHistory)
	})

	return r
}
