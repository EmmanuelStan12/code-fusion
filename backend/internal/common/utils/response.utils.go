package utils

import (
	"encoding/json"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"net/http"
)

type ApiResponse[T any] struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Data       T      `json:"data"`
}

func BuildResponse[T any](data T, success bool, code int, resType string, config *configs.LocaleConfig) ApiResponse[T] {
	message := config.Translate(resType)
	return ApiResponse[T]{
		Success:    success,
		StatusCode: code,
		Message:    message,
		Type:       resType,
		Data:       data,
	}
}

func WriteResponse[T any](writer http.ResponseWriter, data T, success bool, code int, resType string, config *configs.LocaleConfig) {
	message := config.Translate(resType)
	response := ApiResponse[T]{
		Success:    success,
		StatusCode: code,
		Message:    message,
		Type:       resType,
		Data:       data,
	}
	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		panic(err)
	}
}
