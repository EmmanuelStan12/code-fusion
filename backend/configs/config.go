package configs

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	DB     DBConfig
	JWT    JwtConfig    `json:"jwt"`
	Logger LoggerConfig `json:"logging"`
}

func LoadConfig(configPath string) (*AppConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	var config AppConfig
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	dbConfig := InitDBConfig()
	config.DB = dbConfig
	return &config, nil
}
