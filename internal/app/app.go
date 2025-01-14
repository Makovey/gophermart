package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/Makovey/gophermart/internal/config"
	"github.com/Makovey/gophermart/internal/logger"
	"github.com/Makovey/gophermart/internal/middleware"
	"github.com/Makovey/gophermart/internal/service"
	"github.com/Makovey/gophermart/internal/transport"
)

const (
	accrualSystemAddrFlag = "-a"
)

type App struct {
	logger         logger.Logger
	cfg            config.Config
	worker         service.Worker
	handler        transport.HTTPHandler
	authMiddleware middleware.Auth

	wg *sync.WaitGroup
}

func NewApp(
	log logger.Logger,
	cfg config.Config,
	worker service.Worker,
	handler transport.HTTPHandler,
	authMiddleware middleware.Auth,
) *App {
	return &App{
		logger:         log,
		cfg:            cfg,
		worker:         worker,
		handler:        handler,
		authMiddleware: authMiddleware,
		wg:             &sync.WaitGroup{},
	}
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	readyCh := make(chan struct{})
	defer close(readyCh)

	a.startHTTPServer(ctx)
	a.startAccrualSystem(ctx, readyCh)
	a.startAccrualWorker(ctx, readyCh)

	a.wg.Wait()
}

func (a *App) startHTTPServer(ctx context.Context) {
	a.wg.Add(1)
	go a.runHTTPServer(ctx)
}

func (a *App) startAccrualSystem(ctx context.Context, readyCh chan<- struct{}) {
	a.wg.Add(1)
	go a.runAccrualSystem(ctx, readyCh)
}

func (a *App) startAccrualWorker(ctx context.Context, readyCh <-chan struct{}) {
	a.wg.Add(1)
	defer a.wg.Done()

	<-readyCh
	go func() {
		a.worker.ProcessNewOrders(ctx)

		<-ctx.Done()
		a.worker.DownProcess()
	}()
}

func (a *App) runAccrualSystem(ctx context.Context, ready chan<- struct{}) {
	defer a.wg.Done()
	fileLoc := a.cfg.AccrualFileLocation()
	port := a.cfg.AccrualAddress()

	fullPath, err := filepath.Abs(fileLoc)
	if err != nil {
		a.logger.Error(fmt.Sprintf("can't abs absolute path from: %s", fullPath), "error", err.Error())
	}

	cmd := exec.Command(fullPath, accrualSystemAddrFlag, port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		a.logger.Error("can't run accrual system", "err", err.Error())
	}

	time.AfterFunc(time.Second, func() {
		ready <- struct{}{}
	})

	go func() {
		<-ctx.Done()
		if cmd.ProcessState != nil {
			if err = cmd.Process.Kill(); err != nil {
				a.logger.Error("can't kill process with accrual system", "err", err.Error())
			}
		}
	}()

}

func (a *App) runHTTPServer(ctx context.Context) {
	defer a.wg.Done()
	a.logger.Info("starting http server on port: " + a.cfg.RunAddress())

	srv := &http.Server{
		Addr:    a.cfg.RunAddress(),
		Handler: a.initRouter(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			a.logger.Info("server closed", "error", err.Error())
		}
	}()

	<-ctx.Done()
	a.logger.Debug("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		a.logger.Error("server forced to shutdown: %v", "error", err.Error())
	}
}
