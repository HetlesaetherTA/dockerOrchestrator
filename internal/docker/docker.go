package docker

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"hetlesaether.com/dockerOrchestrator/internal/domain"
	"io"
	"log/slog"
	"strings"
	"time"
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

func (c *Client) CreateContainer(
	ctx context.Context,
	name string,
	blueprint domain.Blueprint,
) (domain.Container, error) {
	config, hostConfig := c.MapServiceToMobyConfigs(blueprint)

	slog.Info("Pulling image", "image", config.Image)
	out, err := c.dockerAPI.ImagePull(ctx, config.Image, client.ImagePullOptions{})
	if err != nil {
		return domain.Container{}, fmt.Errorf("failed to pull image %s: %w", config.Image, err)
	}

	_, _ = io.Copy(io.Discard, out)
	out.Close()

	opts := client.ContainerCreateOptions{
		Name:             name,
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: &network.NetworkingConfig{},
	}

	resp, err := c.dockerAPI.ContainerCreate(ctx, opts)

	if err != nil {
		return domain.Container{}, fmt.Errorf("failed to create container: %w", err)
	}

	res, err := c.dockerAPI.ContainerStart(
		ctx,
		resp.ID,
		client.ContainerStartOptions{},
	)

	slog.Info("started container", "res", res, "err", err)

	return domain.Container{
		Name:      name,
		Blueprint: blueprint,
		LastKnown: c.GetContainerByID(resp.ID),
	}, err
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

func (c *Client) GetContainerByID(id string) domain.ContainerState {
	filters := client.Filters{}.Add("id", id)
	res, err := fetchContainerList(
		c,
		context.Background(),
		client.ContainerListOptions{
			All:     true,
			Filters: filters,
		})

	if err != nil {
		return domain.ContainerState{}
	}

	return parseDockerAPIContainer(res[0])
}

// Start given container (by ID)
func (c *Client) StartContainer(container domain.ContainerState) error {
	_, err := c.dockerAPI.ContainerStart(context.Background(), container.ID, client.ContainerStartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container (%s): %w", container.Name, err)
	}

	slog.Info("Container sucessfully start", "Container", container)
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
	slog.Info("Container successfully stopped", "id", container.ID, "name", container.Name)
	return nil
}

func (c *Client) MapServiceToMobyConfigs(blueprint domain.Blueprint) (*container.Config, *container.HostConfig) {
	var image string
	if image = blueprint.Image; image == "" {
		image = "" // TODO: Build image
	}

	mobyConfig := &container.Config{
		Image:        image,
		Cmd:          blueprint.Command,
		Env:          blueprint.Environment.ToMapping().Values(),
		Labels:       blueprint.Labels,
		ExposedPorts: buildExposedPorts(blueprint),
	}

	mobyHostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyMode("no"),
		},
		Binds:   buildVolumeBinds(blueprint.Volumes),
		CapAdd:  blueprint.CapAdd,
		CapDrop: blueprint.CapDrop,
	}

	return mobyConfig, mobyHostConfig
}

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

func buildExposedPorts(blueprint domain.Blueprint) network.PortSet {
	exposedPorts := make(network.PortSet)

	for _, portConfig := range blueprint.ServiceConfig.Ports {
		protocol := portConfig.Protocol
		if protocol == "" {
			protocol = "tcp"
		}

		portNum := uint16(portConfig.Target)

		mobyPort, ok := network.PortFrom(portNum, network.IPProtocol(protocol))
		if !ok {
			continue
		}

		exposedPorts[mobyPort] = struct{}{}
	}

	return exposedPorts
}

func buildVolumeBinds(volumeConfigs []types.ServiceVolumeConfig) []string {
	var binds []string

	for _, vol := range volumeConfigs {
		if vol.Target == "" {
			continue
		}

		if vol.Source == "" {
			continue
		}

		bindStr := fmt.Sprintf("%s:%s", vol.Source, vol.Target)

		// TODO: support ":ao" (append only)
		if vol.ReadOnly {
			bindStr += ":ro"
		}

		binds = append(binds, bindStr)
	}

	return binds
}
