package client

import (
	"context"
	errors2 "errors"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebSocketClient struct {
	Manager     *db.PersistenceManager
	Upgrade     *websocket.Upgrader
	Sessions    map[model.SessionId]WebSocketCollaborators
	SessionCode map[model.SessionId]string
	mu          sync.RWMutex
}

const (
	ErrCannotInitWebSocket = "CANNOT_INIT_WEB_SOCKET"
)

type WebSocketCollaborator struct {
	Collaborator *model.CollaboratorModel
	Connection   *websocket.Conn
	CreatedAt    time.Time
}

type WebSocketCollaborators map[model.CollaboratorID]*WebSocketCollaborator

func NewWebSocketClient(manager *db.PersistenceManager) *WebSocketClient {
	return &WebSocketClient{
		Upgrade: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Sessions:    make(map[model.SessionId]WebSocketCollaborators),
		SessionCode: make(map[model.SessionId]string),
		Manager:     manager,
	}
}

func (client *WebSocketClient) InitWebSocket(w http.ResponseWriter, r *http.Request, collaborator *model.CollaboratorModel, session model.CodeSessionModel, dockerClient *DockerClient) {
	conn, err := client.Upgrade.Upgrade(w, r, nil)
	if err != nil {
		panic(errors.InternalServerError(ErrCannotInitWebSocket, err))
	}
	go client.HandleConnection(conn, dockerClient, session, collaborator)
}

func (client *WebSocketClient) AddCollaborator(session model.CodeSessionModel, collaborator *model.CollaboratorModel, conn *websocket.Conn) (*WebSocketCollaborator, error) {
	collaborators := client.Sessions[session.SessionId]
	if collaborators == nil {
		collaborators = make(WebSocketCollaborators)
		client.Sessions[session.SessionId] = collaborators
	}

	// this should not happen
	if socketC, ok := collaborators[collaborator.ID]; ok {
		return socketC, nil
	}

	webSocketCollaborator := &WebSocketCollaborator{
		Collaborator: collaborator,
		Connection:   conn,
		CreatedAt:    time.Now(),
	}
	collaborators[collaborator.ID] = webSocketCollaborator

	return webSocketCollaborator, nil
}

// IsCollaboratorActive Please call in a thread safe way to prevent deadlock
func (client *WebSocketClient) IsCollaboratorActive(sessionId model.SessionId, collaboratorId model.CollaboratorID) bool {
	collaborators := client.Sessions[sessionId]
	if collaborators == nil {
		return false
	}

	_, ok := collaborators[collaboratorId]
	return ok
}

func (client *WebSocketClient) RemoveCollaborator(session model.CodeSessionModel, collaboratorID model.CollaboratorID) (*websocket.Conn, bool) {
	if collaborators, ok := client.Sessions[session.SessionId]; ok {
		if collaborator, ok := collaborators[collaboratorID]; ok {
			duration := time.Since(collaborator.CreatedAt)
			activeDuration := uint(duration.Minutes()) + collaborator.Collaborator.ActiveDuration
			result := client.Manager.DB.Table("collaborator_models").Where("collaborator_models.id = ?", collaborator.Collaborator.ID).
				Update("last_active", time.Now()).
				Update("active_duration", activeDuration)
			if result.Error != nil {
				log.Printf("Error deleting or removing collaborator %+v\n", result.Error)
			}

			delete(collaborators, collaboratorID)
			if len(collaborators) == 0 {
				delete(client.Sessions, session.SessionId)
			}
			return collaborator.Connection, ok
		}
	}
	return nil, false
}

/*func (client *WebSocketClient) AddSessionQueue(session model.CodeSessionModel) {
	queue, ok := client.SessionQueues[session.SessionId]
	if ok {
		return
	}
	queue = NewCodeSessionQueue(session.SessionId, session.Code)
	queue.RegisterDebounce(2*time.Second, func(code string) {
		client.UpdateCollaborators(session, code)
	})
	client.SessionQueues[session.SessionId] = queue
	queue.Start()
}*/

func (client *WebSocketClient) InitSession(conn *websocket.Conn, session model.CodeSessionModel, collaborator *model.CollaboratorModel) error {
	client.mu.Lock()
	defer client.mu.Unlock()

	_, err := client.AddCollaborator(session, collaborator, conn)
	if err != nil {
		return err
	}

	collaborators := client.Sessions[session.SessionId]

	_, ok := client.SessionCode[session.SessionId]
	if !ok {
		client.SessionCode[session.SessionId] = session.Code
	}

	for i := range session.Collaborators {
		session.Collaborators[i].IsActive = client.IsCollaboratorActive(session.SessionId, session.Collaborators[i].ID)
	}

	code := client.SessionCode[session.SessionId]
	session.Code = code

	if err := conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(session, dto.ActionSessionInitialized)); err != nil {
		return err
	}

	collaborator.IsActive = true
	activeMessage := dto.BuildWebSocketSessionDTO[model.CollaboratorModel](*collaborator, dto.ActionCollaboratorActive)
	for cId, collab := range collaborators {
		if cId == collaborator.ID {
			continue
		}
		if err := collab.Connection.WriteMessage(websocket.TextMessage, activeMessage); err != nil {
			return err
		}
	}

	return nil
}

func (client *WebSocketClient) CloseSession(conn *websocket.Conn, session model.CodeSessionModel, collaborator *model.CollaboratorModel) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.RemoveCollaborator(session, collaborator.ID)
	conn.Close()

	collaborator.IsActive = false
	collaborators, ok := client.Sessions[session.SessionId]
	if ok {
		for _, collab := range collaborators {
			collab.Connection.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[model.CollaboratorModel](*collaborator, dto.ActionCollaboratorInactive))
		}
	} else {
		code, ok := client.SessionCode[session.SessionId]
		if ok {
			result := client.Manager.DB.Table("code_session_models").Where("code_session_models.session_id = ?", session.SessionId).
				Update("code", code)
			if result.Error != nil {
				log.Printf("Error updating code %+v\n", result.Error)
			}
			delete(client.SessionCode, session.SessionId)
		}
		/*queue, ok := client.SessionQueues[session.SessionId]
		if ok {
			delete(client.SessionQueues, session.SessionId)
			queue.CloseQueue()
		}*/
	}
}

func (client *WebSocketClient) UpdateCollaborators(session model.CodeSessionModel, collaborator *model.CollaboratorModel, message *dto.CodeUpdateMessage) {
	client.mu.Lock()
	defer client.mu.Unlock()

	collaborators, ok := client.Sessions[session.SessionId]
	if ok {
		for _, collab := range collaborators {
			if collaborator.ID == collab.Collaborator.ID {
				continue
			}
			collab.Connection.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(message, dto.ActionCodeUpdate))
		}
	}
	client.SessionCode[session.SessionId] = message.Code
}

func (client *WebSocketClient) CreateCollaborator(session model.CodeSessionModel, userId uint, role string, status string) (*model.CollaboratorModel, error) {
	var collaborator model.CollaboratorModel

	result := client.Manager.DB.Where("session_id = ? and user_id = ?", session.ID, userId).First(&collaborator)
	if collaborator.ID != 0 {
		return &collaborator, nil
	}

	var user model.UserModel
	result = client.Manager.DB.First(&user, userId)
	if result.Error != nil {
		return nil, result.Error
	}
	if user.Email == "" {
		return nil, errors2.New("user not found")
	}

	collaborator = model.CollaboratorModel{
		CodeSessionId:  session.ID,
		UserId:         userId,
		Role:           role,
		ActiveDuration: 0,
	}

	result = client.Manager.DB.Create(collaborator)
	if result.Error != nil {
		return nil, result.Error
	}

	return &collaborator, nil
}

func (client *WebSocketClient) GetImageName(language string) string {
	if strings.EqualFold(language, configs.LanguageJavaScript) || strings.EqualFold(language, configs.LanguageTypeScript) {
		return NodeSandboxImage
	}
	if strings.EqualFold(language, configs.LanguagePython) {
		return PythonSandboxImage
	}
	return ""
}

func (client *WebSocketClient) HandleConnection(conn *websocket.Conn, dockerClient *DockerClient, session model.CodeSessionModel, collaborator *model.CollaboratorModel) {
	image := client.GetImageName(string(session.Language))
	if image == "" {
		dockerClient.Logger.Printf("Cannot find image with language: %s\n", session.Language)
		return
	}
	dockerCon, err := dockerClient.AllocateContainer(image)
	if err != nil {
		log.Printf("An error occured %+v\n", err)
		conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(session, dto.ActionSessionError))
		return
	}

	defer func() {
		client.CloseSession(conn, session, collaborator)
		CloseCodeSession(dockerClient, session, dockerCon)
	}()

	err = client.InitSession(conn, session, collaborator)
	if err != nil {
		dockerClient.Logger.Printf("Cannot init session: %v\n", err)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			dockerClient.Logger.Printf("Read error for session %s: %v\n", session.SessionId, err)
			break
		}

		dockerClient.Logger.Printf("Received message: %s", string(message))

		switch messageType {
		case websocket.TextMessage:
			socketMessage, err := dto.ProcessWebSocketMessage(message)
			if err != nil {
				dockerClient.Logger.Printf("Invalid message format: %v\n", err)
				conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[any](nil, dto.ActionMessageTypeNotSupported))
				continue
			}
			client.HandleWebSocketMessageType(dockerCon, session, conn, socketMessage, collaborator)
		case websocket.CloseMessage:
			{
				break
			}
		default:
			dockerClient.Logger.Printf("Unsupported message type received: %d\n", messageType)
			err = conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[any](nil, dto.ActionMessageTypeNotSupported))
			if err != nil {
				dockerClient.Logger.Printf("Write error: %v\n", err)
			}
		}
	}
}

func (client *WebSocketClient) HandleWebSocketMessageType(
	container *DockerContainer,
	session model.CodeSessionModel,
	conn *websocket.Conn,
	message *dto.WebSocketRequestMessage,
	collaborator *model.CollaboratorModel,
) {
	switch message.MessageType {
	case string(dto.ActionCodeExecution):
		contextId := uuid.New().String()

		action, err := message.GetExecuteCodeMessage()
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[any](nil, dto.ActionCannotDecodeActionType))
			return
		}
		err = container.ExecuteCodeRequest(
			session.SessionId,
			contextId,
			string(session.Language),
			action,
			func(response *proto.ExecuteCodeResponse) {
				container.Logger.Printf("Processed message: %v\n", response)
				resultKey := GenerateGrpcResultKey(session.SessionId, contextId)

				conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(response, dto.ActionCodeExecutionSuccess))
				container.mu.Lock()
				delete(container.Results, resultKey)
				container.mu.Unlock()
			})

	case string(dto.ActionCodeUpdate):
		action, err := message.GetCodeUpdateMessage()
		container.Logger.Printf("Updating code: %v\n", err)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[any](nil, dto.ActionCannotDecodeActionType))
			return
		}
		client.UpdateCollaborators(session, collaborator, action)
	}
}

func CloseCodeSession(dockerClient *DockerClient, session model.CodeSessionModel, dockerCon *DockerContainer) {
	ctx, cancel := context.WithCancel(context.Background())
	dockerClient.Logger.Printf("Closing user session %s\n", session.SessionId)
	_, err := dockerCon.GrpcClient.CodeClient.CloseSession(ctx, &proto.CloseSessionRequest{
		SessionId: string(session.SessionId),
	})
	if err != nil {
		dockerClient.Logger.Printf("Unable to close user session %+v\n", err)
	} else {
		dockerClient.Logger.Printf("Successfully closed user session %s\n", session.SessionId)
	}
	cancel()
}
