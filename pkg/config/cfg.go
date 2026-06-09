package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Database       DB
	ListenAddress  string   `yaml:"listen_address" env:"LISTEN_ADDRESS" env-default:":8080"`
	AllowedOrigins []string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS"`
	JWTSecret      string   `yaml:"jwt_secret" env:"JWT_SECRET" env-required:"true"`
	MigrationsPath string   `yaml:"migrations_path" env:"MIGRATIONS_PATH" env-required:"true"`
}

type DB struct {
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Name     string `yaml:"name" env:"DB_NAME"`
	Driver   string `yaml:"driver" env:"DB_DRIVER"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE" env-default:"disable"`
}

func (cfg DB) DSN() string {
	switch cfg.Driver {
	case "sqlite", "sqlite3":
		return cfg.Name
	default:
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
		)
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
