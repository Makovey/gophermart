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
)

const (
	accrualSystemAddrFlag = "-a"
)

type App struct {
	deps *deps
	wg   *sync.WaitGroup
}

func NewApp() *App {
	return &App{deps: newDeps(), wg: &sync.WaitGroup{}}
}

func (a *App) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a.initDependencies()
	a.startHTTPServer(ctx)
	a.startAccrualSystem(ctx)

	a.wg.Wait()
}

func (a *App) initDependencies() {
	a.deps.Config()
	a.deps.Logger()
}

func (a *App) startAccrualSystem(ctx context.Context) {
	a.wg.Add(1)
	go a.runAccrualSystem(ctx)
}

func (a *App) startHTTPServer(ctx context.Context) {
	a.wg.Add(1)
	go a.runHTTPServer(ctx)
}

func (a *App) runAccrualSystem(ctx context.Context) {
	fileLoc := a.deps.Config().AccrualFileLocation()
	port := a.deps.Config().AccrualAddress()

	fullPath, err := filepath.Abs(fileLoc)
	if err != nil {
		a.deps.Logger().Error(fmt.Sprintf("can't abs absolute path from: %s", fullPath), "error", err.Error())
	}

	//cmd := exec.Command(fmt.Sprintf("%s %s=%s", fullPath, accrualSystemAddrFlag, port))
	cmd := exec.Command(fullPath, accrualSystemAddrFlag, port)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		a.deps.Logger().Error("can't run accrual system", "err", err.Error())
	}
	cmd.Wait()

	<-ctx.Done()
	if cmd.ProcessState != nil {
		if err = cmd.Process.Kill(); err != nil {
			a.deps.Logger().Error("can't kill process with accrual system", "err", err.Error())
		}
	}
	a.wg.Done()
}

func (a *App) runHTTPServer(ctx context.Context) {
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

	<-ctx.Done()
	a.deps.Logger().Debug("shutting down http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.deps.CloseAll(); err != nil {
		a.deps.Logger().Error("closed all resources with error", "error", err.Error())
	}

	if err := srv.Shutdown(ctx); err != nil {
		a.deps.Logger().Error("server forced to shutdown: %v", "error", err.Error())
	}
	a.wg.Done()
}
