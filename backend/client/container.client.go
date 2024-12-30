package client

import (
	"context"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
	"google.golang.org/grpc"
)

type CodeExecutionStream struct {
	stream grpc.BidiStreamingClient[proto.ExecuteCodeRequest, proto.ExecuteCodeResponse]
}

// InitCodeExecutionStream initializes the execution stream for the Docker container.
func InitCodeExecutionStream(codeClient proto.CodeExecutionServiceClient, ctx context.Context) (*CodeExecutionStream, error) {
	stream, err := codeClient.ExecuteCode(ctx)
	if err != nil {
		return nil, err
	}

	return &CodeExecutionStream{
		stream: stream,
	}, nil
}
