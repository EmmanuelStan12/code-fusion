package middleware

import (
	"context"
	"github.com/EmmanuelStan12/code-fusion/client"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var excludedPaths = []string{
	"^/api/v1/login",
	"^/api/v1/register",
	"^/$",
	"^/status$",
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

func IsWSPath(path string) bool {
	exp := regexp.MustCompile("^/api/v1/sessions/init")
	return exp.MatchString(path)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if RequireAuth(r.URL.Path) {
			token := ""
			if IsWSPath(r.URL.Path) {
				queryParams := r.URL.Query()
				token = queryParams.Get("token")
			} else {
				token = r.Header.Get("Authorization")
			}
			if token == "" {
				panic(errors.Unauthorized(EmptyOrInvalidToken, nil))
			}
			token = strings.Replace(token, "Bearer ", "", 1)
			jwtUtils := r.Context().Value(JwtContextKey).(client.JwtClient)
			userId, err := strconv.Atoi(jwtUtils.Verify(token))
			if err != nil {
				panic(err)
			}
			manager := r.Context().Value(PersistenceContextKey).(*db.PersistenceManager)
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
