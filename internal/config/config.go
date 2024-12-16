package config

import (
	"strings"

	"github.com/Makovey/gophermart/internal/logger"
)

type Config interface {
	RunAddress() string
	DatabaseURI() string
	AccrualAddress() string
	AccrualFileLocation() string
}

const (
	defaultAddr                = "localhost:8080"
	defaultAccrual             = ":8085"
	defaultAccrualFileLocation = "./cmd/accrual/accrual_darwin_arm64"
)

type config struct {
	runAddress          string
	databaseURI         string
	accrualAddress      string
	accrualFileLocation string
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

func (cfg config) AccrualFileLocation() string {
	return cfg.accrualFileLocation
}

func NewConfig(log logger.Logger) Config {
	envCfg := newEnvConfig(log)
	flags := newFlagsValue()

	addr := runAddress(envCfg, flags)
	dsn := databaseURI(envCfg, flags)
	accrualAddr := accrualAddress(envCfg, flags)
	accrualFileLoc := accrualLocation(envCfg, flags)

	log.Debug("RunAddress: " + addr)
	log.Debug("Database DSN: " + dsn)
	log.Debug("AccrualAddr: " + accrualAddr)
	log.Debug("AccrualFileLocation: " + accrualFileLoc)

	return &config{
		runAddress:          addr,
		databaseURI:         dsn,
		accrualAddress:      accrualAddr,
		accrualFileLocation: accrualFileLoc,
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

func accrualLocation(envCfg envConfig, flags flagsValue) string {
	loc := defaultAccrualFileLocation

	if flags.accrualFileLocation != "" {
		loc = flags.accrualFileLocation
	} else if envCfg.AccrualFileLocation != "" {
		loc = envCfg.AccrualFileLocation
	}

	return loc
}
