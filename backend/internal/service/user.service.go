package service

import (
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
)

const (
	ErrUserDoesNotExist = "USER_DOES_NOT_EXIST"
)

type UserService struct {
	BaseService
}

func NewUserService(jwt utils.JwtUtils, manager *db.PersistenceManager) *UserService {
	authService := UserService{}
	authService.Jwt = jwt
	authService.Manager = manager
	return &authService
}

func (us *UserService) GetUserById(userId int) *model.UserModel {
	var user model.UserModel
	us.Manager.DB.First(&user, userId)
	if user.Email == "" {
		panic(errors.Unauthorized(ErrUserDoesNotExist, nil))
	}

	return &user
}
