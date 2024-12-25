package controllers

import (
	"encoding/json"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/service"
	"net/http"
)

const (
	LoginSuccessful    = "LOGIN_SUCCESSFUL"
	RegisterSuccessful = "REGISTER_SUCCESSFUL"
	ErrDecoding        = "DECODING"
)

type AuthController struct {
	AuthService *service.AuthService
	Locale      *configs.LocaleConfig
}

func NewAuthController(context middleware.AppContext) *AuthController {
	return &AuthController{
		AuthService: &service.AuthService{
			Manager: context.PersistenceManager,
			Jwt:     context.Jwt,
		},
		Locale: context.LocaleConfig,
	}
}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var data dto.LoginDTO
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(errors.BadRequest(ErrDecoding, err))
	}
	result := controller.AuthService.Login(&data)
	utils.WriteResponse[dto.AuthDTO](w, result, true, http.StatusOK, LoginSuccessful, controller.Locale)
}

func (controller *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var data dto.RegisterDTO
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {

	}
	result := controller.AuthService.Register(&data)
	utils.WriteResponse[dto.AuthDTO](w, result, true, http.StatusCreated, RegisterSuccessful, controller.Locale)
}
