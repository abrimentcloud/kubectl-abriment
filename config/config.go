package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

const (
	KubeconfigFileName = "config"

	AbrimentCluster = "abriment-cluster"
	AbrimentContext = "abriment-context"
	AbrimentUser    = "abriment-user"
)

type Config struct {
	LoginEndpoint  string `env:"LOGIN_ENDPOINT" envDefault:"https://backend.abriment.com/dashboard/api/login/"`
	ConfigEndpoint string `env:"CONFIG_ENDPOINT" envDefault:"https://backend.abriment.com/dashboard/api/v1/paas/kubeconfig/"`
}

func ParseCfg() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing env file into config | %v", err)
	}

	return &cfg, nil
}
