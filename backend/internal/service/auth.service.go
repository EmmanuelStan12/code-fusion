package service

import (
	"errors"
	"github.com/EmmanuelStan12/code-fusion/client"
	appErrors "github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/mail"
	"strconv"
)

var (
	ErrInvalidEmail           = "INVALID_EMAIL"
	ErrInvalidLastname        = "INVALID_LASTNAME"
	ErrInvalidFirstname       = "INVALID_FIRSTNAME"
	ErrInvalidUsername        = "INVALID_USERNAME"
	ErrInvalidPassword        = "INVALID_PASSWORD"
	ErrInvalidEmailOrPassword = "INVALID_EMAIL_OR_PASSWORD"
	ErrEmailAlreadyExists     = "EMAIL_ALREADY_EXISTS"
	ErrUsernameAlreadyExists  = "USERNAME_ALREADY_EXISTS"
)

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	log.Printf("ERROR: %v\n", err)
	return err == nil
}

func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

type AuthService struct {
	BaseService
}

func NewAuthService(jwt client.JwtClient, manager *db.PersistenceManager) *AuthService {
	authService := AuthService{}
	authService.Jwt = jwt
	authService.Manager = manager
	return &authService
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func isValidPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *AuthService) Login(data *dto.LoginDTO) dto.AuthDTO {
	var user model.UserModel
	result := s.Manager.DB.Where("email = ?", data.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			panic(appErrors.Unauthorized(ErrInvalidEmailOrPassword, nil))
		}
		panic(result.Error)
	}
	if !isValidPassword(user.Password, data.Password) {
		panic(appErrors.Unauthorized(ErrInvalidEmailOrPassword, nil))
	}
	token := s.Jwt.Create(strconv.Itoa(int(user.ID)))
	return dto.AuthDTO{
		User:  user,
		Token: token,
	}
}

func (s *AuthService) Register(data *dto.RegisterDTO) dto.AuthDTO {
	if !validateEmail(data.Email) {
		panic(appErrors.BadRequest(ErrInvalidEmail, nil))
	}
	user := model.UserModel{}
	result := s.Manager.DB.Find(&user, "email = ?", data.Email)
	if user.Email != "" {
		panic(appErrors.BadRequest(ErrEmailAlreadyExists, nil))
	}
	if !validatePassword(data.Password) {
		panic(appErrors.BadRequest(ErrInvalidPassword, nil))
	}
	if len(data.LastName) == 0 {
		panic(appErrors.BadRequest(ErrInvalidLastname, nil))
	}
	if len(data.FirstName) == 0 {
		panic(appErrors.BadRequest(ErrInvalidFirstname, nil))
	}
	if len(data.Username) == 0 {
		panic(appErrors.BadRequest(ErrInvalidUsername, nil))
	}
	result = s.Manager.DB.Find(&user, "username = ?", data.Username)
	if user.Username != "" {
		panic(appErrors.BadRequest(ErrUsernameAlreadyExists, nil))
	}
	hashedPassword := hashPassword(data.Password)
	user = model.UserModel{
		Firstname: data.FirstName,
		Lastname:  data.LastName,
		Email:     data.Email,
		Password:  hashedPassword,
		Username:  data.Username,
	}
	result = s.Manager.DB.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	token := s.Jwt.Create(strconv.Itoa(int(user.ID)))
	return dto.AuthDTO{
		User:  user,
		Token: token,
	}
}
