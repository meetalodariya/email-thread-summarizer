package config

import "os"

type DBConfig struct {
	Host     string
	Password string
	User     string
	Database string
	Port     string
	SSL      string
}

func LoadDBConfig() *DBConfig {
	config := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}

	return config
}
