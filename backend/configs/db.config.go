package configs

import (
	"fmt"
)

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	SSLMode  string
	Port     string
}

func (config *DBConfig) DSN() string {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		" password=%s dbname=%s sslmode=%s", config.Host, config.Port,
		config.User, config.Password, config.Name, config.SSLMode)
	return psqlInfo
}

func InitDBConfig() DBConfig {
	return DBConfig{
		Name:     GetEnvVar("DB_NAME"),
		User:     GetEnvVar("DB_USER"),
		Password: GetEnvVar("DB_PASSWORD"),
		Host:     GetEnvVar("DB_HOST"),
		SSLMode:  GetEnvVar("DB_SSL_MODE"),
		Port:     GetEnvVar("DB_PORT"),
	}
}
