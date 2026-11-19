// Package config đọc cấu hình từ biến môi trường (giống process.env trong Node).
package config

import "os"

type Config struct {
	Port string
}

func Load() Config {
	return Config{Port: getEnv("PORT", "8080")}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
