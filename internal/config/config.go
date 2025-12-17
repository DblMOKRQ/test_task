package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr     string
	LogLevel     string
	UserRepo     string
	PasswordRepo string
	HostRepo     string
	PortRepo     string
	DBName       string
	SSLMode      string
}

func MustLoad() *Config {
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, reading from environment variables")
	}

	return &Config{
		HTTPAddr:     os.Getenv("HTTP_ADDR"),
		LogLevel:     os.Getenv("LOG_LEVEL"),
		UserRepo:     os.Getenv("DB_USER"),
		PasswordRepo: os.Getenv("DB_PASSWORD"),
		HostRepo:     os.Getenv("DB_HOST"),
		PortRepo:     os.Getenv("DB_PORT"),
		DBName:       os.Getenv("DB_NAME"),
		SSLMode:      os.Getenv("DB_SSLMODE"),
	}
}
