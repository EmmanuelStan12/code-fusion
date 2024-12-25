package configs

import (
	"github.com/joho/godotenv"
	"os"
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
