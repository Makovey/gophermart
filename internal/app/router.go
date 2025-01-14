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

	r.Post("/api/user/register", a.Handler().RegisterUser)
	r.Post("/api/user/login", a.Handler().LoginUser)

	r.Group(func(r chi.Router) {
		authMiddleware := middleware.NewAuth(a.JWT(), a.Logger())

		r.Use(authMiddleware.Authenticate)
		r.Get("/api/user/orders", a.Handler().GetOrders)
		r.Post("/api/user/orders", a.Handler().PostOrder)

		r.Get("/api/user/balance", a.Handler().GetBalance)
		r.Post("/api/user/balance/withdraw", a.Handler().PostWithdraw)

		r.Get("/api/user/withdrawals", a.Handler().GetWithdrawsHistory)
	})

	return r
}
