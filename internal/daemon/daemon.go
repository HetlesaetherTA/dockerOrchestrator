package daemon

import (
	"fmt"
	"hetlesaether.com/dockerOrchestrator/internal/fs"
	"log/slog"
	"os"
	"strings"
	// "hetlesaether.com/dockerOrchestrator/internal/docker"
)

type test struct {
	Test string
}

func Start() {
	env := os.Getenv("APP_ENV")

	if strings.TrimSpace(env) == "" {
		env = "dev"
	}

	slog.Info("Starting application daemon", "environment", env)

	var t test
	fmt.Println(fs.ConstructFromYAML(t, os.Getenv("APP_PATH")+"/hardware.yml"))
	fmt.Println(t.Test)
}
