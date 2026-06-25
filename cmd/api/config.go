package main

import "os"

type Config struct {
	Port   string
	DBPath string
}

func NewConfig() *Config {
	return &Config{
		Port: getEnv("APP_PORT", "8080"),
		DBPath: getEnv("APP_DB_PATH", "store.db"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}