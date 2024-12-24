package middleware

import (
	"context"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"net/http"
	"regexp"
	"strconv"
)

var excludedPaths = []string{
	"^/api/v1/login",
	"^/api/v1/signup",
}

const (
	EmptyOrInvalidToken = "EMPTY_OR_INVALID_TOKEN"
	UserKey             = "UserKey"
)

func RequireAuth(path string) bool {
	for _, excludedPath := range excludedPaths {
		exp := regexp.MustCompile(excludedPath)
		if exp.MatchString(path) {
			return false
		}
	}
	return true
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RequireAuth(r.URL.Path) {
			token := r.Header.Get("Authorization")
			if token == "" {
				panic(errors.Unauthorized(EmptyOrInvalidToken, nil))
			}
			jwtUtils := r.Context().Value(JwtContextKey).(utils.JwtUtils)
			userId, err := strconv.Atoi(jwtUtils.Verify(token))
			if err != nil {
				panic(err)
			}
			manager := r.Context().Value(PersistenceContextKey).(db.PersistenceManager)
			var user model.UserModel
			result := manager.DB.First(&user, userId)
			if result.Error != nil {
				panic(result.Error)
			}
			ctx := context.WithValue(r.Context(), UserKey, user)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
