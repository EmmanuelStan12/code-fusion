package configs

import (
	"fmt"
	"os"
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
		Name:     os.Getenv("db_name"),
		User:     os.Getenv("db_user"),
		Password: os.Getenv("db_password"),
		Host:     os.Getenv("db_host"),
		SSLMode:  os.Getenv("db_ssl_mode"),
		Port:     os.Getenv("db_port"),
	}
}
