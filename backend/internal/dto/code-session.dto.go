package dto

import (
	"encoding/json"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
)

type CreateCodeSessionDTO struct {
	Title           string
	Language        configs.Language
	CollaboratorIds string
}

type CodeSessionOperation struct {
	Type      string `json:"type"`
	Position  int    `json:"position"`
	Text      string `json:"text"`
	Length    int    `json:"length"`
	Timestamp int    `json:"timestamp"`
}

type WebSocketRequestMessage struct {
	MessageType string          `json:"messageType"`
	Data        json.RawMessage `json:"data"`
}

type CodeUpdateMessage struct {
	Operations []CodeSessionOperation `json:"operations"`
	Code       string                 `json:"code"`
}

type WebSocketResponseMessage[T any] struct {
	MessageType string `json:"messageType"`
	Data        T      `json:"data"`
}

func (message *WebSocketRequestMessage) GetExecuteCodeMessage() (*proto.ExecuteCodeRequest, error) {
	var action proto.ExecuteCodeRequest
	if err := json.Unmarshal(message.Data, &action); err != nil {
		return nil, err
	}
	return &action, nil
}

func (message *WebSocketRequestMessage) GetCodeUpdateMessage() (*CodeUpdateMessage, error) {
	var action CodeUpdateMessage
	if err := json.Unmarshal(message.Data, &action); err != nil {
		return nil, err
	}
	return &action, nil
}

func ProcessWebSocketMessage(data []byte) (*WebSocketRequestMessage, error) {
	var message WebSocketRequestMessage
	err := json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}
	return &message, err
}
