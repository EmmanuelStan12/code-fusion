package client

import (
	"github.com/EmmanuelStan12/code-fusion/internal/common/errors"
	"github.com/EmmanuelStan12/code-fusion/internal/db"
	"github.com/EmmanuelStan12/code-fusion/internal/dto"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

type WebSocketClient struct {
	Manager *db.PersistenceManager
	Upgrade *websocket.Upgrader
}

const (
	ErrCannotInitWebSocket  = "CANNOT_INIT_WEB_SOCKET"
	ErrCannotSendMessage    = "CANNOT_SEND_MESSAGE"
	ErrInvalidMessageFormat = "INVALID_MESSAGE_FORMAT"
	ErrCodeExecutionFailed  = "CODE_EXECUTION_FAILED"
)

func NewWebSocketClient() *WebSocketClient {
	return &WebSocketClient{
		Upgrade: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (ws *WebSocketClient) InitWebSocket(w http.ResponseWriter, r *http.Request, session model.CodeSessionModel, dockerClient *DockerClient) {
	conn, err := ws.Upgrade.Upgrade(w, r, nil)
	if err != nil {
		panic(errors.InternalServerError(ErrCannotInitWebSocket, err))
	}
	go ws.HandleConnection(conn, dockerClient, session)
}

func (ws *WebSocketClient) HandleConnection(conn *websocket.Conn, dockerClient *DockerClient, session model.CodeSessionModel) {
	defer func() {
		conn.Close()
	}()
	dockerCon, err := dockerClient.AllocateContainer(NodeSandboxImage)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(session, dto.ActionSessionError))
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO(session, dto.ActionSessionInitialized))
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
			err = dockerCon.ExecuteCodeRequest(session.SessionId, contextId, message, func(response *proto.ExecuteCodeResponse) {
				dockerClient.Logger.Printf("Processed message: %v\n", response)
				resultKey := GenerateGrpcResultKey(session.SessionId, contextId)
				err = conn.WriteMessage(messageType, dto.BuildWebSocketSessionDTO(response, dto.ActionCodeExecutionSuccess))
				dockerCon.mu.Lock()
				delete(dockerCon.Results, resultKey)
				dockerCon.mu.Unlock()
			})
		default:
			dockerClient.Logger.Printf("Unsupported message type received: %d\n", messageType)
			err = conn.WriteMessage(websocket.TextMessage, dto.BuildWebSocketSessionDTO[any](nil, dto.ActionMessageTypeNotSupported))
			if err != nil {
				dockerClient.Logger.Printf("Write error: %v", err)
			}
			break
		}
	}
}
