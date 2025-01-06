package app

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		Target
		Log
		Debug
		TG
		PG
	}

	Target struct {
		Addr string `env:"TARGET_ADDR" env-default:"0.0.0.0:8888"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"DEBUG"`
	}

	Debug struct {
		ServerAddr string `env:"DEBUG_SERVER_ADDR" env-default:"0.0.0.0:8080"`
	}

	TG struct {
		Token string `env:"TG_TOKEN" env-default:"7769410503:AAEmqePfLePAEU7OCjI38x75mnb4M-7bNGs"`
	}

	PG struct {
		Host          string `env:"PG_HOST" env-default:"postgres"`
		Port          string `env:"PG_PORT" env-default:"5432"`
		User          string `env:"PG_USER" env-default:"user"`
		Password      string `env:"PG_PASSWORD" env-default:"password"`
		UserAdmin     string `env:"PG_USER_ADMIN" env-default:"user"`
		PasswordAdmin string `env:"PG_PASSWORD_ADMIN" env-default:"password"`
		Database      string `env:"PG_DATABASE" env-default:"db"`
	}
)

func NewConfig() (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
