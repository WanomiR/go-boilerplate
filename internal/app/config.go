package app

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		Log
		TG
		PG
	}

	Log struct {
		// TODO: instead of required add default log level
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	TG struct {
		Token string `env-required:"true" env:"TG_TOKEN"`
	}

	PG struct {
		Host          string `env-required:"true" env:"PG_HOST"`
		Port          string `env-required:"true" env:"PG_PORT"`
		User          string `env-required:"true" env:"PG_USER"`
		Password      string `env-required:"true" env:"PG_PASSWORD"`
		UserAdmin     string `env-required:"true" env:"PG_USER_ADMIN"`
		PasswordAdmin string `env-required:"true" env:"PG_PASSWORD_ADMIN"`
		Database      string `env-required:"true" env:"PG_DATABASE"`
	}
)

func NewConfig() (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
