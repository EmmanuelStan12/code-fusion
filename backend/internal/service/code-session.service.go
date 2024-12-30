package service

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/google/uuid"
)

type CodeSessionService struct {
	Manager *db.PersistenceManager
}

const (
	ErrInvalidMemoryLimit      = "INVALID_MEMORY_LIMIT"
	ErrInvalidLanguage         = "INVALID_LANGUAGE"
	ErrInvalidTimeout          = "INVALID_TIMEOUT"
	ErrInvalidTitle            = "INVALID_TITLE"
	ErrCannotCreateCodeSession = "CANNOT_CREATE_CODE_SESSION"
)

func NewCodeSessionService(manager *db.PersistenceManager) *CodeSessionService {
	codeSessionService := CodeSessionService{
		Manager: manager,
	}
	return &codeSessionService
}

func (css *CodeSessionService) GetCodeSessionById(sessionId int) *model.CodeSessionModel {
	codeSession := model.CodeSessionModel{}
	result := css.Manager.DB.Find(&codeSession, "sessionId = ?", sessionId)
	if result.Error != nil {
		panic(errors.InternalServerError("CAN_T_FIND_SESSIONS", result.Error))
	}
	return &codeSession
}

func (css *CodeSessionService) CreateSession(userId uint, data *dto.CreateCodeSessionDTO, config *configs.DockerConfig) *model.CodeSessionModel {
	if !config.IsValidLanguage(data.Language) {
		panic(errors.BadRequest(ErrInvalidLanguage, nil))
	}

	if !config.IsValidMemoryLimit(data.MemoryLimit) {
		panic(errors.BadRequest(ErrInvalidMemoryLimit, nil))
	}

	if !config.IsValidTimeout(data.Timeout) {
		panic(errors.BadRequest(ErrInvalidTimeout, nil))
	}

	if data.Title == "" {
		panic(errors.BadRequest(ErrInvalidTitle, nil))
	}

	codeSession := model.CodeSessionModel{
		Title:       data.Title,
		Language:    data.Language,
		SessionId:   uuid.New().String(),
		MemoryLimit: data.MemoryLimit,
		Timeout:     data.Timeout,
		Code:        "",
		Status:      "INACTIVE",
		UserId:      userId,
	}

	result := css.Manager.DB.Create(&codeSession)
	if result.Error != nil {
		panic(errors.InternalServerError(ErrCannotCreateCodeSession, result.Error))
	}
	return &codeSession
}

func (css *CodeSessionService) GetCodeSessionsByUserId(userId uint) []model.CodeSessionModel {
	var sessions []model.CodeSessionModel
	result := css.Manager.DB.Find("userId = ?", userId, &sessions)
	if result.Error != nil {
		panic(errors.InternalServerError("CAN_T_FIND_SESSIONS", result.Error))
	}
	return sessions
}
