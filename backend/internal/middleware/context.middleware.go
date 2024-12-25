package middleware

import (
	"context"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"net/http"
)

type AppContext struct {
	PersistenceManager *db.PersistenceManager
	Jwt                utils.JwtUtils
	LocaleConfig       *configs.LocaleConfig
	Logger             *utils.Logger
}

const JwtContextKey = "JwtContextKey"
const PersistenceContextKey = "PersistenceContextKey"
const LocaleContextKey = "LocaleContextKey"
const LoggerContextKey = "LoggerContextKey"

func ContextMiddleware(ctx AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			newCtx := context.WithValue(r.Context(), JwtContextKey, ctx.Jwt)
			newCtx = context.WithValue(newCtx, PersistenceContextKey, ctx.PersistenceManager)
			newCtx = context.WithValue(newCtx, LocaleContextKey, ctx.LocaleConfig)
			newCtx = context.WithValue(newCtx, LoggerContextKey, ctx.Logger)
			next.ServeHTTP(w, r.WithContext(newCtx))
		}
		return http.HandlerFunc(handler)
	}
}
