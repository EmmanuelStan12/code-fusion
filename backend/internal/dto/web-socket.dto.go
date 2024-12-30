package dto

import (
	"encoding/json"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
)

type ActionType string

const (
	ErrCannotMarshalSessionData              = "CANNOT_MARSHAL_SESSION_DATA"
	ActionSessionInitialized      ActionType = "SESSION_INITIALIZED"
	ActionSessionPaused           ActionType = "SESSION_PAUSED"
	ActionSessionError            ActionType = "SESSION_PAUSED"
	ActionSessionClosed           ActionType = "SESSION_CLOSED"
	ActionCollaboratorAdded       ActionType = "COLLABORATOR_ADDED"
	ActionCollaboratorRemoved     ActionType = "COLLABORATOR_REMOVED"
	ActionCollaboratorUpdated     ActionType = "COLLABORATOR_UPDATED"
	ActionCodeExecutionStarted    ActionType = "CODE_EXECUTION_STARTED"
	ActionCodeExecutionSuccess    ActionType = "CODE_EXECUTION_SUCCESS"
	ActionMessageTypeNotSupported ActionType = "MESSAGE_TYPE_NOT_SUPPORTED"
	ActionCodeExecutionFailed     ActionType = "CODE_EXECUTION_FAILED"
	ActionError                   ActionType = "ERROR"
)

type WebSocketSessionDTO[T any] struct {
	Action ActionType `json:"action"`
	Data   T          `json:"data"`
}

func BuildWebSocketSessionDTO[T any](data T, action ActionType) []byte {
	session := WebSocketSessionDTO[T]{
		Action: action,
		Data:   data,
	}
	result, err := json.Marshal(session)
	if err != nil {
		panic(errors.InternalServerError(ErrCannotMarshalSessionData, err))
	}
	return result
}
