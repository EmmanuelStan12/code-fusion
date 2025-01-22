package model

import "time"

type CollaboratorID uint

const (
	RoleOwner        = "OWNER"
	RoleCollaborator = "COLLABORATOR"

	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
)

type CollaboratorModel struct {
	ID             CollaboratorID   `gorm:"primarykey" json:"id"`
	CodeSessionId  uint             `json:"codeSessionId"`
	CodeSession    CodeSessionModel `gorm:"foreignKey:CodeSessionId" json:"-"`
	UserId         uint             `json:"userId"`
	User           UserModel        `gorm:"foreignKey:UserId" json:"user"`
	Role           string           `json:"role"`
	LastActive     time.Time        `json:"lastActive"`
	ActiveDuration uint             `json:"activeDuration"`
	IsActive       bool             `gorm:"-" json:"isActive"`
}

type CollaboratorDTO struct {
}
