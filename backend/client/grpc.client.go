package client

import (
	"fmt"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcResultKey string

type DockerGrpcClient struct {
	GrpcClient *grpc.ClientConn
	CodeClient proto.CodeExecutionServiceClient
}

// InitGrpcClient initializes a gRPC client connection.
func InitGrpcClient(ip string, port string) (*DockerGrpcClient, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", ip, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Cannot init grpc client for port: %s, %+v\n", port, err)
		return nil, err
	}
	codeClient := proto.NewCodeExecutionServiceClient(conn)
	return &DockerGrpcClient{
		GrpcClient: conn,
		CodeClient: codeClient,
	}, nil
}

// GenerateGrpcResultKey generates a unique result key for gRPC communication.
func GenerateGrpcResultKey(sessionId model.SessionId, contextId string) GrpcResultKey {
	return GrpcResultKey(fmt.Sprintf("%s:%s", sessionId, contextId))
}

// GetIdsFromResultKey splits a GrpcResultKey into session and context IDs.
func GetIdsFromResultKey(resultKey GrpcResultKey) (string, string) {
	r := strings.Split(string(resultKey), ":")
	return r[0], r[1]
}
