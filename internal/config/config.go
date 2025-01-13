package config

import (
	"strings"
	"time"

	"github.com/Makovey/gophermart/internal/logger"
)

type Config interface {
	RunAddress() string
	DatabaseURI() string
	AccrualAddress() string
	AccrualFileLocation() string
	TickerTimer() time.Duration
	AccrualClientTimeout() time.Duration
}

const (
	defaultAddr                = "localhost:8080"
	defaultAccrual             = ":8085"
	defaultAccrualFileLocation = "./cmd/accrual/accrual_darwin_arm64"
	defaultTickerTimer         = "1s"
	defaultClientTimeout       = "10s"
)

type config struct {
	runAddress           string
	databaseURI          string
	accrualAddress       string
	accrualFileLocation  string
	tickerTimer          time.Duration
	accrualClientTimeout time.Duration
}

func (cfg *config) RunAddress() string {
	return cfg.runAddress
}

func (cfg *config) DatabaseURI() string {
	return cfg.databaseURI
}

func (cfg *config) AccrualAddress() string {
	return cfg.accrualAddress
}

func (cfg *config) AccrualFileLocation() string {
	return cfg.accrualFileLocation
}

func (cfg *config) TickerTimer() time.Duration {
	return cfg.tickerTimer
}

func (cfg *config) AccrualClientTimeout() time.Duration {
	return cfg.accrualClientTimeout
}

func NewConfig(log logger.Logger) Config {
	envCfg := newEnvConfig(log)
	flags := newFlagsValue()

	tickerTimer, err := time.ParseDuration(resolveValue(flags.tickerTimer, envCfg.TickerTimer, defaultTickerTimer))
	if err != nil {
		log.Warn("invalid ticker timer duration, used used: " + defaultTickerTimer)
		tickerTimer, _ = time.ParseDuration(defaultTickerTimer)
	}

	clientTimeout, err := time.ParseDuration(resolveValue(flags.accrualClientTimeout, envCfg.AccrualClientTimeout, defaultClientTimeout))
	if err != nil {
		log.Warn("invalid client timeout duration, used default: " + defaultClientTimeout)
		clientTimeout, _ = time.ParseDuration(defaultClientTimeout)
	}

	cfg := &config{
		runAddress:           resolveValue(flags.runAddress, envCfg.RunAddress, defaultAddr),
		databaseURI:          resolveDatabaseURI(flags.databaseURI, envCfg.DatabaseURI),
		accrualAddress:       resolveValue(flags.accrualAddress, envCfg.AccrualAddress, defaultAccrual),
		accrualFileLocation:  resolveValue(flags.accrualFileLocation, envCfg.AccrualFileLocation, defaultAccrualFileLocation),
		tickerTimer:          tickerTimer,
		accrualClientTimeout: clientTimeout,
	}

	log.Debug("RunAddress: " + cfg.runAddress)
	log.Debug("Database DSN: " + cfg.databaseURI)
	log.Debug("AccrualAddr: " + cfg.accrualAddress)
	log.Debug("AccrualFileLocation: " + cfg.accrualFileLocation)
	log.Debug("TickerTimer: " + cfg.tickerTimer.String())
	log.Debug("AccrualClientTimeout: " + cfg.accrualClientTimeout.String())

	return cfg
}

func resolveValue(flagValue, envValue, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if envValue != "" {
		return envValue
	}
	return defaultValue
}

func resolveDatabaseURI(flagValue, envValue string) string {
	dsn := resolveValue(flagValue, envValue, "")
	if dsn != "" && !strings.Contains(dsn, "?sslmode=disable") {
		dsn += "?sslmode=disable"
	}
	return dsn
}
