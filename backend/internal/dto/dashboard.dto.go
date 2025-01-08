package dto

import (
	"github.com/EmmanuelStan12/code-fusion/internal/model"
)

type AnalyticsStats struct {
	TotalActiveDuration float64 `gorm:"column:total_active_duration" json:"totalMinutes"`
	TotalSessions       int     `gorm:"column:total_sessions" json:"totalSessions"`
	TotalLanguages      int     `gorm:"column:total_languages" json:"totalLanguagesUsed"`
	TotalCollaborators  int     `gorm:"column:total_collaborators" json:"totalCollaborators"`
}

type DashboardDTO struct {
	RecentSessions      []model.CodeSessionModel  `json:"recentSessions"`
	RecentCollaborators []model.CollaboratorModel `json:"recentCollaborators"`
	Analytics           AnalyticsStats            `json:"analytics"`
}
