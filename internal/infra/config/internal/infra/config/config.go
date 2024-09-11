package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RateLimitIP            int    `mapstructure:"RATE_LIMIT_IP"`
	RateLimitToken         int    `mapstructure:"RATE_LIMIT_TOKEN"`
	TimeLimitType          string `mapstructure:"TIME_LIMIT_TYPE"`
	TimeBlockType          string `mapstructure:"TIME_BLOCK_TYPE"`
	BlockLimitTimeDuration int
}

func LoadConfig(envFile string) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: %v\n", err)
	}

	ipMaxRequestsPerSecond, err := strconv.Atoi(getEnv("RATE_LIMIT_IP", "10"))
	if err != nil {
		return nil, err
	}

	tokenMaxRequestsPerSecond, err := strconv.Atoi(getEnv("RATE_LIMIT_TOKEN", "10"))
	if err != nil {
		return nil, err
	}

	blockDurationSeconds, err := strconv.Atoi(getEnv("BLOCK_LIMIT_TIME_DURATION", "60"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		RateLimitIP:            ipMaxRequestsPerSecond,
		RateLimitToken:         tokenMaxRequestsPerSecond,
		TimeLimitType:          "second",
		TimeBlockType:          "second",
		BlockLimitTimeDuration: blockDurationSeconds,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
