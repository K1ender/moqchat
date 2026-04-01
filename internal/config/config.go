package config

import "github.com/ilyakaznacheev/cleanenv"

type Env string

const (
	Production  Env = "production"
	Development Env = "development"
)

type Config struct {
	Env Env `env:"ENV" env-default:"production"`
	TLS TLS
}

type TLS struct {
	CertPath string `env:"TLS_CERTIFICATE" env-required:"true"`
	KeyPath  string `env:"TLS_KEY"         env-required:"true"`
}

type Database struct {
	Host string `env:"DATABASE_HOST"     env-required:"true"`
	Port int    `env:"DATABASE_PORT"     env-required:"true"`
	User string `env:"DATABASE_USER"     env-required:"true"`
	Pass string `env:"DATABASE_PASSWORD" env-required:"true"`
	Name string `env:"DATABASE_NAME"     env-required:"true"`
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
