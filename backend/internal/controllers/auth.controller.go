package controllers

import (
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/service"
	"net/http"
)

type AuthController struct {
	AuthService *service.AuthService
}

func NewAuthController(context middleware.AppContext) *AuthController {
	return &AuthController{
		AuthService: &service.AuthService{
			Manager: context.PersistenceManager,
			Jwt:     context.Jwt,
		},
	}
}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	result := controller.AuthService.Login()
}

func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	result := controller.AuthService.Register()
}
