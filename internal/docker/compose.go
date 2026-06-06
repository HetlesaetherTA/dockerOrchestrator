package docker

import (
	"context"

	"fmt"
	"github.com/compose-spec/compose-go/v2/cli"
	"github.com/compose-spec/compose-go/v2/types"
)

func LoadComposeFromPath(ctx context.Context, path string) (types.ServiceConfig, error) {
	opts, err := cli.NewProjectOptions([]string{path})
	if err != nil {
		return types.ServiceConfig{}, fmt.Errorf("failed to parse file options: %w", err)
	}

	project, err := opts.LoadProject(ctx)
	if err != nil {
		return types.ServiceConfig{}, fmt.Errorf("failed to load project: %w", err)
	}

	for _, serviceConfig := range project.Services {
		return serviceConfig, nil
	}

	return types.ServiceConfig{}, fmt.Errorf("no service configuration found in file: %s", path)
}

func CompileBlueprint(projectOptions cli.ProjectOptions) (map[string]any, error) {
	return projectOptions.LoadModel(context.Background())
}
