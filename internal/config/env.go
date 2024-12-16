package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"

	"github.com/Makovey/gophermart/internal/logger"
)

type envConfig struct {
	RunAddress          string `env:"RUN_ADDRESS"`
	DatabaseURI         string `env:"DATABASE_URI"`
	AccrualAddress      string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	AccrualFileLocation string `env:"ACCRUAL_SYSTEM_LOCATION"`
}

func newEnvConfig(log logger.Logger) envConfig {
	fn := "env.newEnvConfig"

	var cfg envConfig
	err := env.Parse(&cfg)
	if err != nil {
		log.Error(fmt.Sprintf("%s: could not parse environment variables", fn), "error", err)
	}

	return cfg
}
