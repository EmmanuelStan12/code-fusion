package service

import (
	errors2 "errors"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"strings"
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
	ErrInvalidRequest          = "INVALID_REQUEST_FORMAT"
)

func NewCodeSessionService(manager *db.PersistenceManager) *CodeSessionService {
	codeSessionService := CodeSessionService{
		Manager: manager,
	}
	return &codeSessionService
}

func (css *CodeSessionService) GetCodeSessionById(sessionId string, userId uint) *model.CodeSessionModel {
	codeSession := model.CodeSessionModel{}
	result := css.Manager.DB.
		Model(&model.CodeSessionModel{}).Preload("Collaborators").
		Preload("Collaborators.User").
		Joins("LEFT JOIN collaborator_models ON collaborator_models.code_session_id = code_session_models.id").
		Where("code_session_models.session_id = ?", sessionId).
		Where("collaborator_models.user_id = ?", userId).
		First(&codeSession)
	if result.Error != nil {
		panic(errors.InternalServerError("CAN_T_FIND_SESSIONS", result.Error))
	}
	return &codeSession
}

func (css *CodeSessionService) CreateSession(userId uint, data *dto.CreateCodeSessionDTO, config *configs.DockerConfig) (*model.CodeSessionModel, *model.CollaboratorModel) {
	if !config.IsValidLanguage(data.Language) {
		panic(errors.BadRequest(ErrInvalidLanguage, nil))
	}

	if data.Title == "" {
		panic(errors.BadRequest(ErrInvalidTitle, nil))
	}

	codeSession := model.CodeSessionModel{
		Title:     data.Title,
		Language:  data.Language,
		SessionId: model.SessionId(uuid.New().String()),
		Code:      "",
	}

	tx := css.Manager.DB.Begin()
	result := css.Manager.DB.Create(&codeSession)
	if result.Error != nil {
		tx.Rollback()
		panic(errors.InternalServerError(ErrCannotCreateCodeSession, result.Error))
	}
	var collaborators []model.CollaboratorModel
	owner := model.CollaboratorModel{
		CodeSessionId: codeSession.ID,
		UserId:        userId,
		Role:          model.RoleOwner,
	}
	collaborators = append(collaborators, owner)

	if data.CollaboratorIds != "" {
		ids := strings.Split(data.CollaboratorIds, ",")
		for _, collabId := range ids {
			if collabId = strings.TrimSpace(collabId); collabId == "" {
				continue
			}
			uId, err := strconv.Atoi(collabId)
			if err != nil {
				tx.Rollback()
				panic(errors.BadRequest(ErrInvalidRequest, err))
			}
			collaborators = append(collaborators, model.CollaboratorModel{
				CodeSessionId: codeSession.ID,
				UserId:        uint(uId),
				Role:          model.RoleCollaborator,
			})
		}
	}
	result = css.Manager.DB.Create(collaborators)
	if result.Error != nil {
		tx.Rollback()
		panic(errors.InternalServerError(ErrCannotCreateCodeSession, result.Error))
	}
	tx.Commit()
	return &codeSession, &owner
}

func (css *CodeSessionService) GetCodeSessionsByUserId(userId uint) []model.CodeSessionModel {
	var sessions []model.CodeSessionModel
	result := css.Manager.DB.Find("userId = ?", userId, &sessions)
	if result.Error != nil {
		panic(errors.InternalServerError("CAN_T_FIND_SESSIONS", result.Error))
	}
	return sessions
}

func (css *CodeSessionService) FindAllCollaborators(session model.CodeSessionModel) []model.CollaboratorModel {
	var collaborators []model.CollaboratorModel
	css.Manager.DB.Where("session_id = ?", session.SessionId).Find(&collaborators)
	return collaborators
}

func (css *CodeSessionService) FindCollaborator(sessionId model.SessionId, userId uint) (*model.CollaboratorModel, error) {
	var collaborator model.CollaboratorModel

	result := css.Manager.DB.Model(&model.CollaboratorModel{}).Preload("User").
		Joins("LEFT JOIN code_session_models cs ON cs.id = collaborator_models.code_session_id").
		Where("cs.session_id = ? and user_id = ?", sessionId, userId).First(&collaborator)
	if errors2.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &collaborator, nil
}
