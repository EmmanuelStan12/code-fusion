package configs

import (
	"github.com/joho/godotenv"
	"os"
)

const (
	EnvHostIP     = "HOST_IP"
	EnvDBName     = "DB_NAME"
	EnvDBUser     = "DB_USER"
	EnvDBPassword = "DB_PASSWORD"
	EnvDBHost     = "DB_HOST"
	EnvSSLMode    = "SSL_MODE"
	EnvDBPort     = "DB_PORT"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func GetEnvVar(key string) string {
	return os.Getenv(key)
}
