package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/Makovey/gophermart/internal/middleware"
)

type App struct {
	deps *deps
}

func NewApp() *App {
	return &App{deps: newDeps()}
}

func (a App) Run() {
	a.runHTTPServer()
}

func (a App) initRouter() http.Handler {
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

func (a App) runHTTPServer() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, os.Kill, syscall.SIGTERM)

	cfg := a.deps.Config()
	a.deps.Logger().Info("starting http server on port: " + cfg.RunAddress())

	srv := &http.Server{
		Addr:    cfg.RunAddress(),
		Handler: a.initRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.deps.Logger().Info("server closed", "error", err.Error())
		}
	}()

	<-shutdown
	a.deps.Logger().Debug("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.deps.CloseAll(); err != nil {
		a.deps.Logger().Error("closed all resources with error", "error", err.Error())
	}

	if err := srv.Shutdown(ctx); err != nil {
		a.deps.Logger().Error("server forced to shutdown: %v", "error", err.Error())
	}
}
