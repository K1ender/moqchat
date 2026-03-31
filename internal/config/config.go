package config

import "github.com/ilyakaznacheev/cleanenv"

type Env string

const (
	Production  Env = "production"
	Development Env = "development"
)

type Config struct {
	Env Env
	TLS TLS
}

type TLS struct {
	CertPath string `env:"TLS_CERTIFICATE"`
	KeyPath  string `env:"TLS_KEY"`
}

func MustInit() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err == nil {
		return &cfg
	}

	err = cleanenv.ReadConfig(".env", &cfg)
	if err == nil {
		return &cfg
	}

	panic(err)
}
