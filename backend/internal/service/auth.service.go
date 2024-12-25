package service

import (
	"errors"
	appErrors "github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"regexp"
	"strconv"
)

var (
	ErrInvalidEmail           = "INVALID_EMAIL"
	ErrEmptyField             = "EMPTY_FIELD"
	ErrInvalidPassword        = "INVALID_PASSWORD"
	ErrInvalidEmailOrPassword = "INVALID_EMAIL_OR_PASSWORD"
	ErrEmailAlreadyExists     = "EMAIL_ALREADY_EXISTS"
)

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return false
	}
	return true
}

func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return true
}

type AuthService struct {
	Manager *db.PersistenceManager
	Jwt     utils.JwtUtils
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
	user := model.UserModel{
		Email: data.Email,
	}
	result := s.Manager.DB.First(&user)
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
	user := model.UserModel{Email: data.Email}
	result := s.Manager.DB.First(&user)
	if result.Error == nil {
		panic(appErrors.BadRequest(ErrEmailAlreadyExists, nil))
	}
	if !validatePassword(data.Password) {
		panic(appErrors.BadRequest(ErrInvalidPassword, nil))
	}
	if len(data.LastName) == 0 {
		panic(appErrors.ValidationError(ErrEmptyField, "lastName"))
	}
	if len(data.FirstName) == 0 {
		panic(appErrors.ValidationError(ErrEmptyField, "firstName"))
	}
	hashedPassword := hashPassword(data.Password)
	user = model.UserModel{
		Firstname: data.FirstName,
		Lastname:  data.LastName,
		Email:     data.Email,
		Password:  hashedPassword,
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
