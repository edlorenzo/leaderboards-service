package config

import "os"

// DbConfig struct
type DbConfig struct {
	DbHost string
	DbUser string
	DbPort string
	DbPass string
	DbName string
}

// Config struct
type Config struct {
	DB     DbConfig
	AppURL string
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		DB: DbConfig{
			DbHost: os.Getenv("DB_HOST"),
			DbPort: os.Getenv("DB_PORT"),
			DbPass: os.Getenv("DB_PASS"),
			DbUser: os.Getenv("DB_USER"),
			DbName: os.Getenv("DB_DBNAME"),
		},
		AppURL: os.Getenv("APP_URL"),
	}
}
