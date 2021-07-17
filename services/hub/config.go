package hub

import "github.com/bfoody/Walmart-Scraper/utils/config"

// A Config contains various credentials, etc loaded from environment
// variables.
type Config struct {
	Env              string `env:"SCR_ENV" default:"dev"` // "dev" or "prod"
	DatabaseURL      string `env:"SCR_DATABASE_URL"`
	DatabasePort     string `env:"SCR_DATABASE_PORT"`
	DatabaseName     string `env:"SCR_DATABASE_NAME"`
	DatabaseUsername string `env:"SCR_DATABASE_USERNAME"`
	DatabasePassword string `env:"SCR_DATABASE_PASSWORD"`
	AMQPURL          string `env:"SCR_AMQP_URL"`
	AMQPExchange     string `env:"SCR_AMQP_EXCHANGE"`
}

// LoadConfig loads all config options from environment variables into
// a *Config.
func LoadConfig() (*Config, error) {
	cfg := Config{}

	err := config.LoadConfigFromEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
