package domain

import (
	"time"
)

// Full applications constructed from one or more containers
type Service struct {
	Restart   string
	Container KnownContainer
}

// Known containers are registered though Labels. Daemon will attempt to repair itself if it loses the container session
type KnownContainer struct {
	Name           string           // homepage.name: "{name}" Label in docker-compose.yml
	DependsOn      []KnownContainer // All Containers must be running to start
	LastKnownState ContainerState   // Last known state
}

// Represent a docker container at some point in time
type ContainerState struct {
	ID      string    // Docker container ID
	Name    string    // Docker container Name
	State   string    // Docker container State
	Updated time.Time // Container struct synced with Docker container at
}
