package configs

import (
	"encoding/json"
	"os"
)

type LocaleConfig struct {
	Messages map[string]string
}

func (locale *LocaleConfig) Translate(key string) string {
	value, ok := locale.Messages[key]
	if !ok {
		return ""
	}
	return value
}

func InitLocale(localePath string) *LocaleConfig {
	file, err := os.Open(localePath)
	var result map[string]string

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&result)
	if err != nil {
		panic(err)
	}
	return &LocaleConfig{
		Messages: result,
	}
}
