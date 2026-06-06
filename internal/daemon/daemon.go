package daemon

import (
	"context"
	"log/slog"
	"os"

	"hetlesaether.com/dockerOrchestrator/internal/docker"
	"hetlesaether.com/dockerOrchestrator/internal/infra"
)

func Start() {
	slog.Info("Starting application daemon")

	infraPath := os.Getenv("APP_PATH") + "/infra.yml"
	infraRepository, err := infra.NewFromYAML(infraPath)

	if err != nil {
		slog.Error("Failed to get or parse infrastructure from yaml", "error", err)
	}

	slog.Info("Found infrastructure from YAML", "read from", infraPath, "infra", infraRepository)

	dockerAPI, err := docker.New()

	if err != nil {
		slog.Error("Failed to initiate docker API Client", "error", err)
	}

	blueprints := SyncBlueprintsInSource()
	container, err := dockerAPI.CreateContainer(context.Background(), blueprints[0].Name, blueprints[0])

	slog.Info("Created Container", "container", container, "err", err)

}
