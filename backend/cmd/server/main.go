package main

import (
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/routes"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

const (
	WelcomeMessage = "WELCOME_MESSAGE"
	ServerStatus   = "SERVER_STATUS"
	RouteNotFound  = "ROUTE_NOT_FOUND"
)

func main() {
	appConfig, err := configs.LoadConfig(os.Getenv("config_file_path"))
	if err != nil {
		panic(err)
	}
	dbManager := db.Init(appConfig.DB)
	localeConfig := configs.InitLocale(os.Getenv("locale_file_path"))
	jwt := utils.JwtUtils{
		JwtConfig: appConfig.JWT,
	}
	appContext := middleware.AppContext{
		PersistenceManager: dbManager,
		Jwt:                jwt,
		LocaleConfig:       localeConfig,
	}
	mainRouter := chi.NewRouter()
	mainRouter.Use(middleware.ErrorMiddleware(localeConfig))
	mainRouter.Use(middleware.ContextMiddleware(appContext))
	mainRouter.Use(middleware.RequestLoggerMiddleware(appConfig.Logger))
	mainRouter.Use(middleware.AuthMiddleware)
	mainRouter.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		utils.WriteResponse[any](writer, nil, true, http.StatusOK, WelcomeMessage, localeConfig)
	})
	mainRouter.Get("/status", func(writer http.ResponseWriter, request *http.Request) {
		utils.WriteResponse[map[string]any](
			writer,
			map[string]any{"dbStatus": dbManager.IsConnected()},
			true,
			http.StatusOK,
			ServerStatus,
			localeConfig,
		)
	})
	mainRouter.Route("/api/v1", func(r chi.Router) {
		r.Mount("/", routes.NewAuthRouter(appContext))
	})
	mainRouter.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		utils.WriteResponse[any](
			writer,
			nil,
			true,
			http.StatusNotFound,
			RouteNotFound,
			localeConfig,
		)
	})
	err = http.ListenAndServe(":3000", mainRouter)
	if err != nil {
		return
	}
}
