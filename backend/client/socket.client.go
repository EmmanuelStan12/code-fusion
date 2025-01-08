package client

import (
	"encoding/json"
	errors2 "errors"
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

type WebSocketClient struct {
	Manager  *db.PersistenceManager
	Upgrade  *websocket.Upgrader
	Sessions map[model.SessionId]WebSocketCollaborators
	mu       sync.RWMutex
}

type WebSocketRequestMessage struct {
	MessageType string                   `json:"messageType"`
	Data        proto.ExecuteCodeRequest `json:"data"`
	UserId      uint                     `json:"userId"`
}

type WebSocketResponseMessage[T any] struct {
	MessageType string `json:"messageType"`
	Data        T      `json:"data"`
}

const (
	ErrCannotInitWebSocket  = "CANNOT_INIT_WEB_SOCKET"
	ErrCannotSendMessage    = "CANNOT_SEND_MESSAGE"
	ErrInvalidMessageFormat = "INVALID_MESSAGE_FORMAT"
	ErrCodeExecutionFailed  = "CODE_EXECUTION_FAILED"
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
		Sessions: make(map[model.SessionId]WebSocketCollaborators),
		Manager:  manager,
	}
}

func (client *WebSocketClient) InitWebSocket(w http.ResponseWriter, r *http.Request, collaborator *model.CollaboratorModel, session model.CodeSessionModel, dockerClient *DockerClient) {
	conn, err := client.Upgrade.Upgrade(w, r, nil)
	if err != nil {
		panic(errors.InternalServerError(ErrCannotInitWebSocket, err))
	}
	socketCollaborator, err := client.AddCollaborator(session, collaborator, conn)
	if err != nil {
		panic(errors.InternalServerError(ErrCannotInitWebSocket, err))
	}
	go client.HandleConnection(conn, dockerClient, session, socketCollaborator)
}

func (client *WebSocketClient) AddCollaborator(session model.CodeSessionModel, collaborator *model.CollaboratorModel, conn *websocket.Conn) (*WebSocketCollaborator, error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	collaborators := client.Sessions[session.SessionId]
	if collaborators == nil {
		collaborators = make(WebSocketCollaborators)
		client.Sessions[session.SessionId] = collaborators
	}

	webSocketCollaborator := &WebSocketCollaborator{
		Collaborator: collaborator,
		Connection:   conn,
		CreatedAt:    time.Now(),
	}
	collaborators[collaborator.ID] = webSocketCollaborator

	return webSocketCollaborator, nil
}

func (client *WebSocketClient) IsCollaboratorActive(sessionId model.SessionId, collaboratorId model.CollaboratorID) bool {
	client.mu.RLock()
	defer client.mu.RUnlock()

	collaborators := client.Sessions[sessionId]
	if collaborators == nil {
		return false
	}

	_, ok := collaborators[collaboratorId]
	return ok
}

func (client *WebSocketClient) RemoveCollaborator(session model.CodeSessionModel, collaboratorID model.CollaboratorID) (*websocket.Conn, bool) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if collaborators, ok := client.Sessions[session.SessionId]; ok {
		if collaborator, ok := collaborators[collaboratorID]; ok {
			duration := time.Since(collaborator.CreatedAt)
			activeDuration := uint(duration.Minutes()) + collaborator.Collaborator.ActiveDuration
			result := client.Manager.DB.Table("collaborator_models").Where("collaborator_models.id = ?", collaborator.Collaborator.ID).
				Update("last_active", time.Now()).
				Update("active_duration", activeDuration)
			if result.Error != nil {
				log.Printf("Error deleting or removing collaborator %+v\n", result.Error)
				return nil, false
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

func (client *WebSocketClient) InitSession(conn *websocket.Conn, session model.CodeSessionModel, collaborator *WebSocketCollaborator) error {
	err := conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(session, dto.ActionSessionInitialized))
	if err != nil {
		return err
	}

	client.mu.RLock()
	defer client.mu.RUnlock()

	collaborators, ok := client.Sessions[session.SessionId]
	if ok {
		for cId, collab := range collaborators {
			if cId == collaborator.Collaborator.ID {
				continue
			}
			collab.Connection.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[model.CollaboratorModel](*collaborator.Collaborator, dto.ActionCollaboratorActive))
		}
	}

	return nil
}

func (client *WebSocketClient) CloseSession(conn *websocket.Conn, session model.CodeSessionModel, collaborator *WebSocketCollaborator) {
	client.RemoveCollaborator(session, collaborator.Collaborator.ID)
	conn.Close()

	client.mu.RLock()
	defer client.mu.RUnlock()

	collaborators, ok := client.Sessions[session.SessionId]
	if ok {
		for _, collab := range collaborators {
			collab.Connection.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[model.CollaboratorModel](*collaborator.Collaborator, dto.ActionCollaboratorInactive))
		}
	}
}

func (client *WebSocketClient) UpdateSession(session model.CodeSessionModel, collaborator model.CollaboratorModel, code string) error {
	client.mu.RLock()
	defer client.mu.RUnlock()

	collaborators, ok := client.Sessions[session.SessionId]
	if ok {
		for cId, collab := range collaborators {
			if cId == collaborator.ID {
				continue
			}
			collab.Connection.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(struct {
				Collaborator model.CollaboratorModel `json:"collaborator"`
				Session      model.CodeSessionModel  `json:"session"`
				Code         string                  `json:"code"`
			}{Collaborator: collaborator, Session: session, Code: code}, dto.ActionCodeUpdate))
		}
	}

	return nil
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

func (client *WebSocketClient) HandleConnection(conn *websocket.Conn, dockerClient *DockerClient, session model.CodeSessionModel, collaborator *WebSocketCollaborator) {
	defer func() {
		client.CloseSession(conn, session, collaborator)
	}()
	dockerCon, err := dockerClient.AllocateContainer(NodeSandboxImage)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(session, dto.ActionSessionError))
		return
	}

	err = client.InitSession(conn, session, collaborator)
	if err != nil {
		dockerClient.Logger.Printf("Cannot write message: %v\n", err)
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
			contextId := uuid.New().String()
			var webSocketMessage WebSocketRequestMessage
			err := json.Unmarshal(message, &webSocketMessage)
			if err != nil {
				dockerClient.Logger.Printf("Invalid message format: %v\n", err)
				err = conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[any](nil, dto.ActionMessageTypeNotSupported))
				continue
			}

			switch webSocketMessage.MessageType {
			case string(dto.ActionCodeExecution):
				err = dockerCon.ExecuteCodeRequest(session.SessionId, contextId, &webSocketMessage.Data, func(response *proto.ExecuteCodeResponse) {
					dockerClient.Logger.Printf("Processed message: %v\n", response)
					resultKey := GenerateGrpcResultKey(session.SessionId, contextId)

					err = conn.WriteMessage(messageType, dto.BuildWebSocketSessionDTO(response, dto.ActionCodeExecutionSuccess))

					dockerCon.mu.Lock()
					delete(dockerCon.Results, resultKey)
					dockerCon.mu.Unlock()
				})

			case string(dto.ActionAddCollaborator):
				newCollaborator, err := client.CreateCollaborator(session, webSocketMessage.UserId, model.RoleCollaborator, model.StatusInactive)
				if err != nil {
					dockerClient.Logger.Printf("Unable to add collaborator: %v\n", err)
					conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[any](nil, dto.ActionAddCollaboratorError))
					continue
				}
				conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[model.CollaboratorModel](*newCollaborator, dto.ActionAddCollaboratorSuccess))

			case string(dto.ActionCodeUpdate):
				err := client.UpdateSession(session, *collaborator.Collaborator, webSocketMessage.Data.Code)
				if err != nil {
					dockerClient.Logger.Printf("Error updating code: %v\n", err)
				}
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
