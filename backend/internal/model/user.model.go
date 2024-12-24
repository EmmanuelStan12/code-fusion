package model

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"-"`
}
