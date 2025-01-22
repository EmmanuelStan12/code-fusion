package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/model"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const (
	NodeSandboxImage   = "node_sandbox:latest"
	PythonSandboxImage = "python_sandbox:latest"
)

type DockerContainerKey string

type DockerClient struct {
	Client         *client.Client
	HostIP         string
	AllocatedPorts map[string]bool
	Containers     map[DockerContainerKey]*DockerContainer
	Config         *configs.DockerConfig
	Logger         *log.Logger
	mu             sync.RWMutex
}

type DockerContainer struct {
	ImageName           string
	Port                string
	Id                  string
	GrpcClient          *DockerGrpcClient
	Results             map[GrpcResultKey]func(*proto.ExecuteCodeResponse)
	CodeExecutionStream *CodeExecutionStream
	mu                  sync.RWMutex
	Logger              *log.Logger
}

func calculateCPUPercentage(stats *container.StatsResponse) float64 {
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	if systemDelta > 0 && cpuDelta > 0 {
		return (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}
	return 0.0
}

func (dc *DockerClient) GetContainerStatus(key DockerContainerKey) (*DockerContainer, bool) {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	con, exists := dc.Containers[key]
	return con, exists
}

func (dc *DockerClient) AddContainer(key DockerContainerKey, container *DockerContainer) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.Containers[key] = container
}

func (dc *DockerClient) CanAllocateMoreTasks(imageName string) *DockerContainer {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	for containerId, con := range dc.Containers {
		if con.ImageName != imageName {
			continue
		}
		stats, err := dc.Client.ContainerStats(context.Background(), string(containerId), false)
		if err != nil {
			dc.Logger.Printf("Error getting container stats: %v", err)
			continue
		}

		var statsData container.StatsResponse
		if err := json.NewDecoder(stats.Body).Decode(&statsData); err != nil {
			dc.Logger.Printf("Error decoding stats data: %v", err)
			continue
		}
		memoryLimit := float64(statsData.MemoryStats.Limit) / (1024 * 1024)
		memoryUsage := float64(statsData.MemoryStats.Usage) / (1024 * 1024)
		dc.Logger.Printf("Memory Usage: %.2f MiB / %.2f MiB\n", memoryUsage, memoryLimit)
		percentage := memoryUsage / memoryLimit
		if percentage < 0.75 {
			cpuUsage := calculateCPUPercentage(&statsData)
			if cpuUsage < 0.75 {
				stats.Body.Close()
				return con
			}
		}
	}
	return nil
}

func checkImageLocally(cli *client.Client, imageName string) bool {
	ctx := context.Background()

	images, err := cli.ImageList(ctx, image.ListOptions{
		All: true,
	})
	if err != nil {
		log.Fatalf("Error listing images: %v", err)
	}

	for _, localImage := range images {
		for _, tag := range localImage.RepoTags {
			if tag == imageName {
				return true
			}
		}
	}
	return false
}

func NewDockerClient() *DockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	if !checkImageLocally(cli, NodeSandboxImage) {
		panic(errors.New("Node sandbox image does not exist on the local machine"))
	}

	return &DockerClient{
		Client:         cli,
		HostIP:         configs.GetEnvVar(configs.EnvHostIP),
		AllocatedPorts: make(map[string]bool),
		Containers:     make(map[DockerContainerKey]*DockerContainer),
		Config:         configs.NewDockerConfig(),
		Logger:         log.New(os.Stdout, "[DockerContainer]", log.LstdFlags),
		mu:             sync.RWMutex{},
	}
}

func (dc *DockerClient) Dispose() error {
	dc.mu.RLock()
	defer dc.mu.RUnlock()
	for containerId, _ := range dc.Containers {
		err := dc.Client.ContainerRemove(context.Background(), string(containerId), container.RemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   true,
			Force:         true,
		})
		if err != nil {
			dc.Logger.Printf("Error getting container stats: %s\n", err)
		} else {
			dc.Logger.Printf("Successfully removed container: %s\n", containerId)
		}
	}
	return nil
}

func (dc *DockerClient) createDockerContainer(imageName string) (*DockerContainer, error) {
	ctx := context.Background()
	port, err := dc.generateFreePort()
	if err != nil {
		log.Printf("Cannot generate free port %+v\n", err)
		return nil, err
	}
	grpcPort := fmt.Sprintf("%d", port)
	res, err := dc.Client.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Env:   []string{fmt.Sprintf("PORT=%s", grpcPort)},
			ExposedPorts: nat.PortSet{
				nat.Port(grpcPort): struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(grpcPort): []nat.PortBinding{
					{HostIP: dc.HostIP, HostPort: grpcPort},
				},
			},
		},
		nil,
		nil,
		"",
	)

	if err != nil {
		log.Printf("Cannot create container %+v\n", err)
		return nil, err
	}

	err = dc.Client.ContainerStart(context.Background(), res.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Cannot start container %+v\n", err)
		e := dc.Client.ContainerRemove(context.Background(), res.ID, container.RemoveOptions{})
		if e != nil {
			log.Printf("Cannot remove container %+v\n", err)
			return nil, e
		}
		return nil, err
	}

	grpcConn, err := InitGrpcClient(dc.HostIP, grpcPort)

	if err != nil {
		dc.Logger.Printf("Unable to create grpc port for container %s, closing container...\n", res.ID)
		return nil, err
	}

	dcContainer := &DockerContainer{
		ImageName:  imageName,
		Port:       grpcPort,
		Id:         res.ID,
		Results:    make(map[GrpcResultKey]func(*proto.ExecuteCodeResponse)),
		GrpcClient: grpcConn,
		Logger:     dc.Logger,
	}
	containerKey := DockerContainerKey(res.ID)
	dc.Containers[containerKey] = dcContainer

	return dcContainer, nil
}

func (dc *DockerClient) AllocateContainer(imageName string) (*DockerContainer, error) {
	con := dc.CanAllocateMoreTasks(imageName)
	if con == nil {
		newContainer, err := dc.createDockerContainer(imageName)
		if err != nil {
			log.Printf("Cannot allocate container %+v\n", err)
			return nil, err
		}
		con = newContainer
	}

	err := WaitForServerReady(con.GrpcClient.GrpcClient, time.Minute)
	if err != nil {
		return nil, err
	}

	if con.CodeExecutionStream == nil {
		ctx, cancel := context.WithCancel(context.Background())

		executionStream, err := InitCodeExecutionStream(con.GrpcClient.CodeClient, ctx)
		if err != nil {
			cancel()
			return nil, err
		}
		log.Printf("Initialized execution stream...")
		con.CodeExecutionStream = executionStream

		// There is a backpressure issue here, can cause memory issues
		go func() {
			defer cancel()
			for {
				result, err := executionStream.stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					cancel()
					dc.Logger.Printf("Error receiving stream: %v\n", err)
					return
				}
				resultKey := GenerateGrpcResultKey(model.SessionId(result.SessionId), result.ContextId)
				con.mu.RLock()
				callback, exists := con.Results[resultKey]
				con.mu.RUnlock()
				if exists {
					go callback(result)
				} else {
					dc.Logger.Printf("Channel does not exist for key %v, ignoring result %v...", resultKey, result)
				}
			}
		}()
	}

	return con, nil
}

func (dc *DockerContainer) ExecuteCodeRequest(
	sessionId model.SessionId,
	contextId string,
	language string,
	executeRequest *proto.ExecuteCodeRequest,
	callback func(response *proto.ExecuteCodeResponse),
) error {
	executeRequest.SessionId = string(sessionId)
	executeRequest.ContextId = contextId
	executeRequest.Language = language

	resultKey := GenerateGrpcResultKey(sessionId, contextId)

	dc.mu.Lock()
	dc.Results[resultKey] = callback
	dc.mu.Unlock()
	err := dc.CodeExecutionStream.stream.Send(executeRequest)
	if err != nil {
		return err
	}

	return nil
}

func (dc *DockerClient) generateFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, fmt.Errorf("failed to generate free port: %w", err)
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}

func WaitForServerReady(conn *grpc.ClientConn, timeout time.Duration) error {
	healthClient := grpc_health_v1.NewHealthClient(conn)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

		resp, err := healthClient.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: "CODE_EXECUTION_SERVICE"})
		if err == nil && resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			log.Println("gRPC server is healthy and ready.")
			cancel()
			return nil
		}

		log.Printf("Waiting for gRPC server to be ready: %v\n", err)
		cancel()
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("gRPC server did not become ready within the timeout")
}
