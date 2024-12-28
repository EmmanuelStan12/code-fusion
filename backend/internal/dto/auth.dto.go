package dto

import "github.com/EmmanuelStan12/code-fusion/internal/model"

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterDTO struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
}

type AuthDTO struct {
	User  model.UserModel `json:"user"`
	Token string          `json:"token"`
}
