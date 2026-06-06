package daemon

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"hetlesaether.com/dockerOrchestrator/internal/docker"
	"hetlesaether.com/dockerOrchestrator/internal/domain"
)

func ConstructBlueprint(ctx context.Context, path string) (domain.Blueprint, error) {
	config, err := docker.LoadComposeFromPath(ctx, path)

	if err != nil {
		return domain.Blueprint{}, err
	}

	slog.Info("Containerblueprint", "config", config)
	return domain.Blueprint{
		config,
	}, nil
}

func SyncBlueprintsInSource() []domain.Blueprint {
	var tmp []domain.Blueprint

	err := filepath.WalkDir(os.Getenv("BLUEPRINT_PATH"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			err = fmt.Errorf("prevent panic by handling error at %q: %v\n", path, err)
			return err
		}

		if strings.HasSuffix(strings.ToLower(d.Name()), ".yml") || strings.HasSuffix(strings.ToLower(d.Name()), ".yaml") {
			blueprint, err := ConstructBlueprint(context.Background(), path)
			tmp = append(tmp, blueprint)
			slog.Info("Adding containers from blueprint", "blueprint", blueprint, "error", err)
		}

		return nil
	})
	if err != nil {
		slog.Error("ReadDir failed", "error", err)
	}
	return tmp
}
