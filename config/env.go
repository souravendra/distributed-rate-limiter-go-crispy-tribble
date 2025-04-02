package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type RateLimitConfig struct {
	Strategy string
	Rate     int
	Interval int
	Burst    int
}

var cfg RateLimitConfig

func LoadEnv() {
	_ = godotenv.Load()

	cfg.Strategy = getEnv("RATE_LIMIT_STRATEGY", "token")
	cfg.Rate = getEnvAsInt("RATE", 5)
	cfg.Interval = getEnvAsInt("INTERVAL", 1)
	cfg.Burst = getEnvAsInt("BURST", 5)
}

func Get() RateLimitConfig {
	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}
