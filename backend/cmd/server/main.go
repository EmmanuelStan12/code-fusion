package main

import (
	"fmt"
	"github.com/EmmanuelStan12/code-fusion/client"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
	"os"
	"path/filepath"
)

const (
	WelcomeMessage = "WELCOME_MESSAGE"
	ServerStatus   = "SERVER_STATUS"
	RouteNotFound  = "ROUTE_NOT_FOUND"
)

func resolveResourcePath(name string) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := filepath.Join(wd, "resources", name)
	return path
}

func initMigrations(manager *db.PersistenceManager) {
	manager.RegisterEntity(&model.UserModel{})
	manager.Migrate()
}

func main() {
	configs.LoadEnv()
	configPath := resolveResourcePath("config.json")
	appConfig, err := configs.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	dbManager := db.Init(appConfig.DB)
	initMigrations(dbManager)
	localePath := resolveResourcePath("messages.json")
	localeConfig := configs.InitLocale(localePath)
	jwt := client.JwtClient{
		JwtConfig: appConfig.JWT,
	}
	logger := client.NewLogger(appConfig.LogLevel)
	dockerClient := client.NewDockerClient()
	socketClient := client.NewWebSocketClient()
	appContext := middleware.AppContext{
		PersistenceManager: dbManager,
		Jwt:                jwt,
		LocaleConfig:       localeConfig,
		Logger:             logger,
		DockerClient:       dockerClient,
		SocketClient:       socketClient,
	}
	mainRouter := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))
	mainRouter.Use(middleware.ErrorMiddleware(localeConfig, logger))
	mainRouter.Use(middleware.ContextMiddleware(appContext))
	mainRouter.Use(middleware.RequestLoggerMiddleware(logger))
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
		r.Mount("/users", routes.NewUserRouter(appContext))
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
	port := configs.GetEnvVar("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Printf("Server listening on PORT :%s\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), mainRouter)
	if err != nil {
		panic(err)
	}
}
