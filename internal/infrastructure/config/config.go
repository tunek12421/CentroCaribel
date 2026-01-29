package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
	Admin  AdminConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret               string
	ExpirationHours      int
	RefreshExpirationHrs int
}

type AdminConfig struct {
	Email    string
	Password string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "caribel"),
			Password: getEnv("DB_PASSWORD", "caribel_secret_2024"),
			Name:     getEnv("DB_NAME", "centro_caribel"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", "cambiar-este-secret-en-produccion"),
			ExpirationHours:      getEnvInt("JWT_EXPIRATION_HOURS", 8),
			RefreshExpirationHrs: getEnvInt("JWT_REFRESH_EXPIRATION_HOURS", 72),
		},
		Admin: AdminConfig{
			Email:    getEnv("ADMIN_EMAIL", "admin@centrocaribel.com"),
			Password: getEnv("ADMIN_PASSWORD", "Admin123!"),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
