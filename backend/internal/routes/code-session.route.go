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
	router.Get("/", controller.GetUserCodeSessions)
	router.Get("/{sessionId}", controller.GetCodeSessionById)
	router.Post("/create", controller.CreateSession)
	router.Get("/init/{sessionId}", controller.InitCodeSession)
	return router
}
