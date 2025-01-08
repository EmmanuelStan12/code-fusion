package controllers

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/service"
	"net/http"
)

type AnalyticsController struct {
	Manager *db.PersistenceManager
	Locale  *configs.LocaleConfig
	Service *service.AnalyticsService
}

func NewAnalyticsController(context middleware.AppContext) *AnalyticsController {
	return &AnalyticsController{
		Locale:  context.LocaleConfig,
		Manager: context.PersistenceManager,
		Service: service.NewAnalyticsService(context),
	}
}

func (controller *AnalyticsController) GetDashboardAnalytics(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(model.UserModel)
	stats := controller.Service.GetSummary(user)
	recentSessions := controller.Service.GetRecentSessions(user)
	recentCollaborators := controller.Service.GetRecentCollaborators(user)

	result := dto.DashboardDTO{
		Analytics:           *stats,
		RecentCollaborators: recentCollaborators,
		RecentSessions:      recentSessions,
	}

	utils.WriteResponse[dto.DashboardDTO](w, result, true, http.StatusOK, LoginSuccessful, controller.Locale)
}
