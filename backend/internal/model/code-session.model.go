package model

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"time"
)

type CodeSessionModel struct {
	ID          uint                `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Title       string              `json:"title"`
	Language    configs.Language    `json:"language"`
	SessionId   string              `json:"sessionId"`
	MemoryLimit configs.MemoryLimit `json:"memoryLimit"`
	Timeout     configs.Timeout     `json:"timeout"`
	Code        string              `json:"code"`
	Status      string              `json:"status"`
	UserId      uint                `json:"userId"`
	User        UserModel           `gorm:"foreignKey:UserId" json:"-"`
}
