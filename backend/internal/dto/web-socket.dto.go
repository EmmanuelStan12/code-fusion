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
	ActionSessionError            ActionType = "SESSION_ERROR"
	ActionSessionClosed           ActionType = "SESSION_CLOSED"
	ActionCollaboratorActive      ActionType = "COLLABORATOR_ACTIVE"
	ActionCollaboratorInactive    ActionType = "COLLABORATOR_INACTIVE"
	ActionCodeExecutionSuccess    ActionType = "CODE_EXECUTION_SUCCESS"
	ActionCodeExecution           ActionType = "CODE_EXECUTION"
	ActionCodeUpdate              ActionType = "CODE_UPDATE"
	ActionAddCollaborator         ActionType = "ADD_COLLABORATOR"
	ActionAddCollaboratorError    ActionType = "ADD_COLLABORATOR_ERROR"
	ActionAddCollaboratorSuccess  ActionType = "ADD_COLLABORATOR_SUCCESS"
	ActionMessageTypeNotSupported ActionType = "MESSAGE_TYPE_NOT_SUPPORTED"
	ActionCannotDecodeActionType  ActionType = "CANNOT_DECODE_ACTION_TYPE"
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
