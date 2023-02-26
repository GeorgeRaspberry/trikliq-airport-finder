package database

import (
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

var database DBConfig

func init() {

	database = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		Database: os.Getenv("DB_NAME"),
	}

	Connect()
}
