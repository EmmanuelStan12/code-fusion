package main

import (
	"context"
	"fmt"
	pb "github.com/EmmanuelStan12/code-fusion/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.NewClient("localhost:4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewCodeExecutionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)

	defer cancel()

	stream, err := client.ExecuteCode(ctx)

	if err != nil {
		panic(err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			result, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				panic(err)
			}
			log.Printf("Code execution result %+v\n", result)
		}
	}()

	for i := 0; i < 7; i++ {
		err = stream.Send(&pb.ExecuteCodeRequest{
			Code:      fmt.Sprintf("console.log(%d);\nconst i = %d; 5 * i", i, i),
			SessionId: strconv.Itoa(rand.Intn(1000)),
		})
		if err != nil {
			panic(err)
		}
	}
	err = stream.CloseSend()
	<-waitc
	if err != nil {
		panic(err)
	}
}
