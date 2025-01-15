package main

import (
	"github.com/Makovey/gophermart/internal/service/worker"
	"github.com/Makovey/gophermart/internal/transport/accrual"
	"os"

	"github.com/Makovey/gophermart/internal/app"
	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger/slog"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/repository/postgresql"
	"github.com/Makovey/gophermart/internal/service/gophermart"
	"github.com/Makovey/gophermart/internal/transport/http"
	"github.com/Makovey/gophermart/pkg/jwt"
)

func main() {
	log := slog.NewLogger(slog.Local)
	cfg := config.NewConfig(log)
	jwt := jwt.NewJWT(log)

	repo, err := postgresql.NewPostgresRepo(log, cfg)
	if err != nil {
		log.Error("failed to initialize postgres repository", "error", err)
		os.Exit(1)
	}

	appl := app.NewApp(
		log,
		cfg,
		worker.NewWorker(repo, accrual.NewHTTPClient(cfg, log), cfg, log),
		http.NewHTTPHandler(
			log,
			gophermart.NewUserService(repo, jwt),
			gophermart.NewOrderService(repo),
			gophermart.NewBalanceService(repo),
			gophermart.NewHistoryService(repo),
		),
		middleware.NewAuth(jwt, log),
	)

	appl.Run()

	if err = repo.Close(); err != nil {
		log.Error("closed all resources with error", "error", err.Error())
	}
}
