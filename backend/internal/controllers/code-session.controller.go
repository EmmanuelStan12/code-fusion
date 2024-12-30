package controllers

import (
	"encoding/json"
	"github.com/EmmanuelStan12/code-fusion/client"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/utils"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/middleware"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	ErrInvalidSessionId = "INVALID_SESSION_ID"
	SessionsRetrieved   = "SESSIONS_RETRIEVED"
)

type CodeSessionController struct {
	CodeSessionService *service.CodeSessionService
	Locale             *configs.LocaleConfig
}

func NewCodeSessionController(context middleware.AppContext) *CodeSessionController {
	return &CodeSessionController{
		CodeSessionService: &service.CodeSessionService{
			Manager: context.PersistenceManager,
		},
		Locale: context.LocaleConfig,
	}
}

func (c *CodeSessionController) GetCodeSessionById(w http.ResponseWriter, r *http.Request) {
	sessionId := chi.URLParam(r, "sessionId")
	codeSession := c.CodeSessionService.GetCodeSessionById(sessionId)
	utils.WriteResponse[model.CodeSessionModel](w, *codeSession, true, http.StatusOK, SessionsRetrieved, c.Locale)
}

func (c *CodeSessionController) CreateSession(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(model.UserModel)
	dockerClient := r.Context().Value(middleware.DockerClientKey).(*client.DockerClient)
	var data dto.CreateCodeSessionDTO
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		panic(errors.BadRequest(ErrDecoding, err))
	}
	session := c.CodeSessionService.CreateSession(user.ID, &data, dockerClient.Config)
	utils.WriteResponse[model.CodeSessionModel](w, *session, true, http.StatusOK, SessionsRetrieved, c.Locale)
}

func (c *CodeSessionController) InitCodeSession(w http.ResponseWriter, r *http.Request) {
	sessionId := chi.URLParam(r, "sessionId")
	if sessionId == "" {
		panic(errors.BadRequest(ErrInvalidSessionId, nil))
	}
	codeSession := c.CodeSessionService.GetCodeSessionById(sessionId)
	socketClient := r.Context().Value(middleware.WebSocketClient).(*client.WebSocketClient)
	dockerClient := r.Context().Value(middleware.DockerClientKey).(*client.DockerClient)

	socketClient.InitWebSocket(w, r, *codeSession, dockerClient)
}

func (c *CodeSessionController) GetUserCodeSessions(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(middleware.UserKey).(model.UserModel)
	sessions := c.CodeSessionService.GetCodeSessionsByUserId(user.ID)
	utils.WriteResponse[[]model.CodeSessionModel](w, sessions, true, http.StatusOK, SessionsRetrieved, c.Locale)
}
