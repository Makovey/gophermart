package config

import (
	"strings"

	"github.com/Makovey/gophermart/internal/logger"
)

type Config interface {
	RunAddress() string
	DatabaseURI() string
	AccrualAddress() string
}

const (
	defaultAddr    = "localhost:8080"
	defaultAccrual = "localhost:8081"
)

type config struct {
	runAddress     string
	databaseURI    string
	accrualAddress string
}

func (cfg config) RunAddress() string {
	return cfg.runAddress
}

func (cfg config) DatabaseURI() string {
	return cfg.databaseURI
}

func (cfg config) AccrualAddress() string {
	return cfg.accrualAddress
}

func NewConfig(log logger.Logger) Config {
	envCfg := newEnvConfig(log)
	flags := newFlagsValue()

	addr := runAddress(envCfg, flags)
	dsn := databaseURI(envCfg, flags)
	accrualAddr := accrualAddress(envCfg, flags)

	log.Debug("RunAddress: " + addr)
	log.Debug("Database DSN: " + dsn)
	log.Debug("AccrualAddr: " + accrualAddr)

	return &config{
		runAddress:     addr,
		databaseURI:    dsn,
		accrualAddress: accrualAddr,
	}
}

func runAddress(envCfg envConfig, flags flagsValue) string {
	addr := defaultAddr

	if flags.runAddress != "" {
		addr = flags.runAddress
	} else if envCfg.RunAddress != "" {
		addr = envCfg.RunAddress
	}

	return addr
}

func databaseURI(envCfg envConfig, flags flagsValue) string {
	var databaseDSN string

	if flags.databaseURI != "" {
		databaseDSN = flags.databaseURI
	} else if envCfg.DatabaseURI != "" {
		databaseDSN = envCfg.DatabaseURI
	}

	if databaseDSN != "" && !strings.Contains(databaseDSN, "?sslmode=disable") {
		databaseDSN = databaseDSN + "?sslmode=disable"
	}

	return databaseDSN
}

func accrualAddress(envCfg envConfig, flags flagsValue) string {
	addr := defaultAccrual

	if flags.accrualAddress != "" {
		addr = flags.accrualAddress
	} else if envCfg.AccrualAddress != "" {
		addr = envCfg.AccrualAddress
	}

	return addr
}
