package config

import (
	"log/slog"
	"os"
	"path"

	"github.com/caarlos0/env"
	dotenv "github.com/joho/godotenv"
)

type Config struct {
	Logger  *slog.Logger
	SsePort int `env:"SSE_PORT" envDefault:"8989"`
	WebPort int `env:"WEB_PORT" envDefault:"8990"`
}

func xdgConfig() string {
	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		return path.Join(xdgHome, "mcp_human", "config")
	}

	return path.Join(os.Getenv("HOME"), ".config", "mcp_human", "config")
}

func MustLoad(logger *slog.Logger) *Config {
	conf, err := Load(logger)
	if err != nil {
		panic(err)
	}

	return conf
}

func Load(logger *slog.Logger) (*Config, error) {
	conf := new(Config)

	for _, name := range []string{".env", xdgConfig()} {
		if _, err := os.Stat(name); err != nil {
			logger.Debug("config file not found",
				slog.String("name", name))

			continue
		}

		logger.Info("loading config file",
			slog.String("name", name))

		if err := dotenv.Load(name); err != nil {
			return nil, err
		}

		break
	}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	conf.Logger = logger

	return conf, nil
}
