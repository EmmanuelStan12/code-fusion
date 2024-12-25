package utils

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	ErrorInvalidToken   = "JWT_INVALID_TOKEN"
	ErrorTokenExpired   = "JWT_TOKEN_EXPIRED"
	ErrorSigningFailure = "JWT_SIGNING_FAILURE"
	ErrorInvalidSubject = "JWT_INVALID_SUBJECT"
)

type JwtUtils struct {
	configs.JwtConfig
}

func (config *JwtUtils) GetExp() int64 {
	return time.Now().Add(time.Hour * 24).Unix()
}

func (config *JwtUtils) Create(subject string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": subject,
		"exp": config.GetExp(),
		"iss": config.Issuer,
		"aud": config.Audience,
	})
	tkStr, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		panic(errors.InternalServerError(ErrorSigningFailure, err))
	}
	return tkStr
}

func (config *JwtUtils) Verify(tkStr string) string {
	token, err := jwt.Parse(tkStr, func(token *jwt.Token) (interface{}, error) {
		return config.SecretKey, nil
	})
	if err != nil {
		panic(errors.Unauthorized(ErrorTokenExpired, err))
	}
	if !token.Valid {
		panic(errors.Unauthorized(ErrorInvalidToken, err))
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		panic(errors.Unauthorized(ErrorInvalidSubject, err))
	}
	return subject
}
