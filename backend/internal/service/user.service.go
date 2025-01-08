package service

import (
	errors2 "errors"
	"github.com/EmmanuelStan12/code-fusion/client"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"gorm.io/gorm"
)

const (
	ErrUserDoesNotExist = "USER_DOES_NOT_EXIST"
)

type UserService struct {
	BaseService
}

func NewUserService(jwt client.JwtClient, manager *db.PersistenceManager) *UserService {
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

func (s *UserService) GetUsers(user model.UserModel) []model.UserModel {
	var users []model.UserModel

	err := s.Manager.DB.Model(&model.UserModel{}).
		Where("user_models.id <> ?", user.ID).
		Scan(&users).Error

	if err != nil && !errors2.Is(err, gorm.ErrRecordNotFound) {
		panic(errors.InternalServerError("UNABLE_TO_RETRIEVE_USERS", err))
	}
	return users
}
