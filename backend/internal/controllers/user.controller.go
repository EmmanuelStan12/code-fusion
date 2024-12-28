package controllers

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/service"
	"net/http"
)

type UserController struct {
	UserService *service.UserService
	Locale      *configs.LocaleConfig
}

const (
	UserRetrieved = "USER_RETRIEVED"
)

func NewUserController(context middleware.AppContext) *UserController {
	return &UserController{
		UserService: service.NewUserService(context.Jwt, context.PersistenceManager),
		Locale:      context.LocaleConfig,
	}
}

func (controller *UserController) GetAuthUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(model.UserModel)
	utils.WriteResponse[model.UserModel](w, user, true, http.StatusOK, UserRetrieved, controller.Locale)
}
