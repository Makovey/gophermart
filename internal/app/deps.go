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

func (p *deps) Logger() logger.Logger {
	if p.logger == nil {
		p.logger = slog.NewLogger(slog.Local)
	}

	return p.logger
}

func (p *deps) Service() transport.GophermartService {
	if p.service == nil {
		p.service = gophermart.NewGophermartService(p.Repository(), p.Logger(), p.JWT())
	}

	return p.service
}

func (p *deps) Repository() service.GophermartRepository {
	if p.repo == nil {
		p.repo, _ = postgresql.NewPostgresRepo(p.Logger(), p.Config())
	}

	return p.repo
}

func (p *deps) Handler() transport.HTTPHandler {
	if p.handler == nil {
		p.handler = http.NewHTTPHandler(p.Logger(), p.Service())
	}

	return p.handler
}

func (p *deps) Config() config.Config {
	if p.cfg == nil {
		p.cfg = config.NewConfig(p.Logger())
	}

	return p.cfg
}

func (p *deps) JWT() *jwt.JWT {
	if p.jwt == nil {
		p.jwt = jwt.NewJWT(p.Logger())
	}

	return p.jwt
}

func (p *deps) CloseAll() error {
	if closer, ok := p.Repository().(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
