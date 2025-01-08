package routes

import (
	"github.com/EmmanuelStan12/code-fusion/internal/controllers"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewDashboardRouter(context middleware.AppContext) http.Handler {
	router := chi.NewRouter()
	controller := controllers.NewAnalyticsController(context)
	router.Get("/", controller.GetDashboardAnalytics)
	return router
}
