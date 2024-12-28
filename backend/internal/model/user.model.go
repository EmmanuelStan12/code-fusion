package model

import (
	"database/sql"
	"time"
)

type UserModel struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"-"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	Firstname string       `json:"firstName"`
	Lastname  string       `json:"lastName"`
	Email     string       `json:"email"`
	Username  string       `json:"username"`
	Password  string       `json:"-"`
}
