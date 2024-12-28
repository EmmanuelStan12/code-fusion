package routes

import (
	"github.com/EmmanuelStan12/code-fusion/internal/controllers"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewAuthRouter(context middleware.AppContext) http.Handler {
	router := chi.NewRouter()
	controller := controllers.NewAuthController(context)
	router.Post("/login", controller.Login)
	router.Post("/register", controller.Register)
	return router
}
