package model

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"time"
)

type SessionId string

type CodeSessionModel struct {
	ID            uint                `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
	Title         string              `json:"title"`
	Language      configs.Language    `json:"language"`
	SessionId     SessionId           `json:"sessionId"`
	MemoryLimit   configs.MemoryLimit `json:"memoryLimit"`
	Timeout       configs.Timeout     `json:"timeout"`
	Code          string              `json:"code"`
	Collaborators []CollaboratorModel `gorm:"foreignKey:CodeSessionId;constraint:OnDelete:CASCADE" json:"collaborators"`
}
