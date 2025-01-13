package main

import (
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

	serv := gophermart.NewGophermartService(repo, log, jwt)
	handler := http.NewHTTPHandler(log, serv)
	auth := middleware.NewAuth(jwt, log)

	appl := app.NewApp(log, cfg, repo, handler, auth)

	appl.Run()
}
