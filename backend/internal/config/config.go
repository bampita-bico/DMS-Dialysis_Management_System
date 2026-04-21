package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env      string
	HTTPAddr string

	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	DBMaxConns int32

	JWTSecret string
}

func Load() (Config, error) {
	// Optional .env
	_ = godotenv.Load()

	cfg := Config{}
	cfg.Env = getenv("APP_ENV", "dev")
	cfg.HTTPAddr = getenv("HTTP_ADDR", ":8080")

	cfg.DBHost = getenv("DB_HOST", "localhost")
	port, err := mustInt("DB_PORT", getenv("DB_PORT", "5432"))
	if err != nil {
		return Config{}, err
	}
	cfg.DBPort = port
	cfg.DBName = getenv("DB_NAME", "dms")
	cfg.DBUser = getenv("DB_USER", "dms")
	cfg.DBPassword = getenv("DB_PASSWORD", "dms_dev_password")
	cfg.DBSSLMode = getenv("DB_SSLMODE", "disable")
	maxConns, err := mustInt("DB_MAX_CONNS", getenv("DB_MAX_CONNS", "4"))
	if err != nil {
		return Config{}, err
	}
	cfg.DBMaxConns = int32(maxConns)

	cfg.JWTSecret = getenv("JWT_SECRET", "")
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func mustInt(key, v string) (int, error) {
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %q", key, v)
	}
	return i, nil
}
