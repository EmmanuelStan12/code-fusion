package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/EmmanuelStan12/code-fusion/configs"
	"github.com/EmmanuelStan12/code-fusion/internal/proto"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	NodeSandboxImage = "node_sandbox"
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
	Name                string
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
	for containerKey, con := range dc.Containers {
		imageKey, containerId := GetIdsFromContainerKey(containerKey)
		if imageKey != imageName {
			continue
		}
		stats, err := dc.Client.ContainerStats(context.Background(), containerId, false)
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

func (dc *DockerClient) createDockerContainer(imageName string) (*DockerContainer, error) {
	ctx := context.Background()
	port, err := dc.generateFreePort()
	if err != nil {
		return nil, err
	}
	containerName := fmt.Sprintf("%s-%s", imageName, uuid.New().String())
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
		containerName,
	)

	grpcConn, err := InitGrpcClient(dc.HostIP, grpcPort)

	if err != nil {
		dc.Logger.Printf("Unable to create grpc port for container %s, closing container...\n", res.ID)
		return nil, err
	}

	dcContainer := &DockerContainer{
		ImageName:  imageName,
		Port:       grpcPort,
		Id:         res.ID,
		Name:       containerName,
		GrpcClient: grpcConn,
		Logger:     dc.Logger,
	}
	containerKey := GetDockerContainerKey(res.ID, imageName)
	dc.Containers[containerKey] = dcContainer

	return dcContainer, nil
}

func (dc *DockerClient) AllocateContainer(imageName string) (*DockerContainer, error) {
	con := dc.CanAllocateMoreTasks(imageName)
	if con == nil {
		con, err := dc.createDockerContainer(imageName)
		if err != nil {
			return nil, err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)

		executionStream, err := InitCodeExecutionStream(con.GrpcClient.CodeClient, ctx)
		con.CodeExecutionStream = executionStream

		// There is a backpressure issue here, can cause memory issues
		go func() {
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
				resultKey := GenerateGrpcResultKey(result.SessionId, result.ContextId)
				con.mu.RLock()
				callback, exists := con.Results[resultKey]
				con.mu.RUnlock()
				if exists {
					go func() {
						callback(result)
					}()
				} else {
					dc.Logger.Printf("Channel does not exist for key %v, ignoring result %v...", resultKey, result)
				}
			}
		}()

		if err != nil {
			return nil, err
		}
		return con, nil
	}

	return con, nil
}

func (dc *DockerContainer) ExecuteCodeRequest(sessionId, contextId string, message []byte, callback func(response *proto.ExecuteCodeResponse)) error {
	var executeRequest proto.ExecuteCodeRequest
	err := json.Unmarshal(message, &executeRequest)
	if err != nil {
		dc.Logger.Printf("Invalid message format: %v\n", err)
		return err
	}

	resultKey := GenerateGrpcResultKey(sessionId, contextId)
	err = dc.CodeExecutionStream.stream.Send(&executeRequest)
	if err != nil {
		return err
	}

	dc.mu.Lock()
	dc.Results[resultKey] = callback
	dc.mu.Unlock()
	err = dc.CodeExecutionStream.stream.Send(&executeRequest)
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

func GetDockerContainerKey(containerId, imageName string) DockerContainerKey {
	return DockerContainerKey(fmt.Sprintf("%s:%s", containerId, imageName))
}

func GetIdsFromContainerKey(key DockerContainerKey) (string, string) {
	r := strings.Split(string(key), ":")
	return r[0], r[1]
}
