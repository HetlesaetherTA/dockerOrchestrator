package main

import (
	"hetlesaether.com/dockerOrchestrator/internal/daemon"
	"log/slog"
	"os"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	daemon.Start()
}
