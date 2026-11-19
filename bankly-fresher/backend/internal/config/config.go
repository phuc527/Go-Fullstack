// Package config đọc cấu hình từ biến môi trường (giống process.env trong Node).
package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8090"),
		// Mặc định trỏ vào Postgres của compose.yaml (cổng 5433 trên máy, tránh đụng Postgres cài sẵn).
		DatabaseURL: getEnv("DATABASE_URL", "postgres://bankly:bankly@localhost:5433/bankly?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
