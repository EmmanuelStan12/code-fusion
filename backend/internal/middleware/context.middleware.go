package middleware

import (
	"context"
	"github.com/EmmanuelStan12/code-fusion/client"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"net/http"
)

type AppContext struct {
	PersistenceManager *db.PersistenceManager
	Jwt                client.JwtClient
	LocaleConfig       *configs.LocaleConfig
	Logger             *client.Logger
	DockerClient       *client.DockerClient
	SocketClient       *client.WebSocketClient
}

const JwtContextKey = "JwtContextKey"
const PersistenceContextKey = "PersistenceContextKey"
const LocaleContextKey = "LocaleContextKey"
const LoggerContextKey = "LoggerContextKey"
const DockerClientKey = "DockerClientKey"
const WebSocketClient = "WebSocketClientKey"

func ContextMiddleware(ctx AppContext) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			newCtx := context.WithValue(r.Context(), JwtContextKey, ctx.Jwt)
			newCtx = context.WithValue(newCtx, PersistenceContextKey, ctx.PersistenceManager)
			newCtx = context.WithValue(newCtx, LocaleContextKey, ctx.LocaleConfig)
			newCtx = context.WithValue(newCtx, LoggerContextKey, ctx.Logger)
			newCtx = context.WithValue(newCtx, DockerClientKey, ctx.DockerClient)
			newCtx = context.WithValue(newCtx, WebSocketClient, ctx.SocketClient)
			next.ServeHTTP(w, r.WithContext(newCtx))
		}
		return http.HandlerFunc(handler)
	}
}
