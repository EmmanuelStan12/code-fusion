package controllers

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"net/http"
)

type AnalyticsController struct {
	Manager *db.PersistenceManager
	Locale  *configs.LocaleConfig
}

func NewAnalyticsController(context middleware.AppContext) *AnalyticsController {
	return &AnalyticsController{
		Locale:  context.LocaleConfig,
		Manager: context.PersistenceManager,
	}
}

func (controller *AnalyticsController) GetDashboardAnalytics(w http.ResponseWriter, r *http.Request) {

}
