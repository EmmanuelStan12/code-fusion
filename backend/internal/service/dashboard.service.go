package service

import (
	errors2 "errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"gorm.io/gorm"
)

const (
	ErrGettingDashboardStats = "GETTING_DASHBOARD_STATS"
)

type AnalyticsService struct {
	Manager *db.PersistenceManager
}

func NewAnalyticsService(context middleware.AppContext) *AnalyticsService {
	return &AnalyticsService{
		Manager: context.PersistenceManager,
	}
}

func (s *AnalyticsService) GetSummary(user model.UserModel) *dto.AnalyticsStats {
	var stats dto.AnalyticsStats

	err := s.Manager.DB.
		Table("collaborator_models").
		Select(`
        COALESCE(SUM(collaborator_models.active_duration)) as total_active_duration,
        COUNT(cs.id) as total_sessions,
        COUNT(DISTINCT cs.language) as total_languages,
        COUNT(DISTINCT other_collaborators.user_id) AS total_collaborators
    `).
		Joins("LEFT JOIN user_models um ON um.id = collaborator_models.user_id").
		Joins("LEFT JOIN code_session_models cs ON cs.id = collaborator_models.code_session_id").
		Joins(`
        LEFT JOIN collaborator_models other_collaborators 
        ON other_collaborators.code_session_id = cs.id 
        AND other_collaborators.user_id <> collaborator_models.user_id
    `).
		Where("um.id = ?", user.ID).
		Scan(&stats).Error

	if err != nil {
		panic(errors.InternalServerError(ErrGettingDashboardStats, err))
	}

	return &stats
}

func (s *AnalyticsService) GetRecentSessions(user model.UserModel) []model.CodeSessionModel {
	var sessions []model.CodeSessionModel

	err := s.Manager.DB.
		Model(&model.CodeSessionModel{}).
		Select("code_session_models.*").
		Joins("LEFT JOIN collaborator_models c ON code_session_models.id = c.code_session_id").
		Where("c.user_id = ? AND c.role = ?", user.ID, model.RoleOwner).
		Order("code_session_models.created_at DESC").
		Limit(10).
		Find(&sessions).Error

	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		panic(errors.InternalServerError(ErrGettingDashboardStats, err))
	}

	return sessions
}

func (s *AnalyticsService) GetRecentCollaborators(user model.UserModel) []model.CollaboratorModel {
	var collaborators []model.CollaboratorModel

	err := s.Manager.DB.Model(&model.CollaboratorModel{}).
		Select("others.*").
		Joins("LEFT JOIN code_session_models cs ON cs.id = collaborator_models.code_session_id").
		Joins(`LEFT JOIN collaborator_models others 
					ON others.code_session_id = cs.id
		`).
		Where("collaborator_models.user_id = ? AND collaborator_models.role = ? AND others.user_id <> ?", user.ID, model.RoleOwner, user.ID).
		Order("others.last_active DESC").
		Limit(10).
		First(&collaborators).Error

	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		panic(errors.InternalServerError(ErrGettingDashboardStats, err))
	}

	return collaborators
}
