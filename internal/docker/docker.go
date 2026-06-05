package docker

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"hetlesaether.com/dockerOrchestrator/internal/domain"
)

type Client struct {
	dockerAPI *client.Client
}

// Returns a new docker api instance
func New() (*Client, error) {
	apiClient, err := client.New(client.FromEnv)

	if err != nil {
		return nil, fmt.Errorf("Failed to create Docker API connection: %w", err)
	}

	return &Client{
		dockerAPI: apiClient,
	}, nil
}

// Get all running containers
func (c *Client) GetRunningContainers() ([]domain.ContainerState, error) {
	resFilters := client.Filters{}

	resFilters.Add("status", "running")

	res, err := fetchContainerList(
		c,
		context.Background(),
		client.ContainerListOptions{
			All:     true,
			Filters: resFilters,
		})

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch containers %w", err)
	}

	return parseDockerAPIContainerList(res), nil
}

// Get all containers (also not running)
func (c *Client) GetContainers() ([]domain.ContainerState, error) {
	res, err := fetchContainerList(
		c,
		context.Background(),
		client.ContainerListOptions{
			All: true,
		})

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch containers %w", err)
	}

	return parseDockerAPIContainerList(res), nil
}

// Start given container (by ID)
func (c *Client) StartContainer(container domain.ContainerState) error {
	_, err := c.dockerAPI.ContainerStart(context.Background(), container.ID, client.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container (%s): %w", container.Name, err)
	}

	slog.Info("Container sucessfully start", container)
	return nil
}

// Gracefully stop given container
func (c *Client) StopContainer(container domain.ContainerState) error {
	timeoutSeconds := 10
	_, err := c.dockerAPI.ContainerStop(
		context.Background(),
		container.ID,
		client.ContainerStopOptions{
			Timeout: &timeoutSeconds,
		})

	if err != nil {
		return fmt.Errorf("Failed to stop container ID(%s): %w", container.ID, err)
	}
[]
	slog.Info("Container successfully stopped", "id", container.ID, "name", container.Name)
	return nil
}

// Find containers with labals: homepage.name is defined
func (c *Client) ScanForKnownContainers() (map[string]domain.ContainerState, error) {
	queryFilters := client.Filters{}
	queryFilters.Add("label", "homepage.name")

	containers, err := fetchContainerList(
		c,
		context.Background(),
		client.ContainerListOptions{
			All:     true,
			Filters: queryFilters,
		})

	if err != nil {
		return nil, fmt.Errorf("failed to list containers by label: %w", err)
	}

	tmp := make(map[string]domain.ContainerState)

	for _, c := range containers {
		if labelValue, exists := c.Labels["homepage.name"]; exists {
			tmp[labelValue] = parseDockerAPIContainer(c)
		}
	}

	return tmp, nil

}

// Private functions

func parseDockerAPIContainerList(c []container.Summary) []domain.ContainerState {
	tmp := make([]domain.ContainerState, len(c))

	for inx, cur := range c {
		tmp[inx] = parseDockerAPIContainer(cur)
	}

	return tmp
}

func fetchContainerList(c *Client, ctx context.Context, opt client.ContainerListOptions) ([]container.Summary, error) {
	res, err := c.dockerAPI.ContainerList(ctx, opt)

	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func parseDockerAPIContainer(c container.Summary) domain.ContainerState {
	var name string

	if len(c.Names) > 0 {
		name = strings.TrimPrefix(c.Names[0], "/")
	}

	return domain.ContainerState{
		ID:      c.ID,
		State:   string(c.State),
		Name:    name,
		Updated: time.Now(),
	}
}
