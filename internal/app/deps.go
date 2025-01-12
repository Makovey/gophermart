package app

import (
	"io"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/logger/slog"
	"github.com/Makovey/gophermart/internal/repository/postgresql"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/service/gophermart"
	"github.com/Makovey/gophermart/internal/transport"
	"github.com/Makovey/gophermart/internal/transport/http"
	"github.com/Makovey/gophermart/pkg/jwt"
)

type deps struct {
	logger logger.Logger
	cfg    config.Config
	jwt    *jwt.JWT

	handler transport.HTTPHandler
	service transport.GophermartService
	repo    service.GophermartRepository
}

func newDeps() *deps {
	return &deps{}
}

func (d *deps) initDependencies() error {
	log := d.makeLogger()
	cfg := d.makeConfig(log)
	jwt := d.makeJWT(log)

	repo, err := d.makeRepository(log, cfg)
	if err != nil {
		return err
	}

	serv := d.makeService(repo, log, jwt)
	d.makeHandler(log, serv)

	return nil
}

func (d *deps) makeLogger() logger.Logger {
	d.logger = slog.NewLogger(slog.Local)
	return d.logger
}

func (d *deps) makeService(
	repo service.GophermartRepository,
	logger logger.Logger,
	jwt *jwt.JWT,
) transport.GophermartService {
	d.service = gophermart.NewGophermartService(repo, logger, jwt)
	return d.service
}

func (d *deps) makeRepository(
	logger logger.Logger,
	cfg config.Config,
) (service.GophermartRepository, error) {
	repo, err := postgresql.NewPostgresRepo(logger, cfg)
	if err != nil {
		return nil, err
	}
	d.repo = repo
	return d.repo, nil
}

func (d *deps) makeHandler(logger logger.Logger, service transport.GophermartService) transport.HTTPHandler {
	d.handler = http.NewHTTPHandler(logger, service)
	return d.handler
}

func (d *deps) makeConfig(logger logger.Logger) config.Config {
	d.cfg = config.NewConfig(logger)
	return d.cfg
}

func (d *deps) makeJWT(logger logger.Logger) *jwt.JWT {
	d.jwt = jwt.NewJWT(logger)
	return d.jwt
}

func (d *deps) CloseAll() error {
	if closer, ok := d.repo.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
