package routes

import (
	"github.com/EmmanuelStan12/code-fusion/internal/controllers"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewCodeSessionRouter(context middleware.AppContext) http.Handler {
	router := chi.NewRouter()
	controller := controllers.NewCodeSessionController(context)
	router.Get("/sessions", controller.GetUserCodeSessions)
	router.Post("/sessions/{sessionId}", controller.GetCodeSessionById)
	router.Get("/sessions/create", controller.CreateSession)
	return router
}