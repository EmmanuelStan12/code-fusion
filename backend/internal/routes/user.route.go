package routes

import (
	"github.com/EmmanuelStan12/code-fusion/internal/controllers"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewUserRouter(context middleware.AppContext) http.Handler {
	router := chi.NewRouter()
	controller := controllers.NewUserController(context)
	router.Get("/me", controller.GetAuthUser)
	router.Get("/", controller.GetUsers)
	return router
}
