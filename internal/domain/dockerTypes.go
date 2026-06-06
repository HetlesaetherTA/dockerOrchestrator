package domain

import (
	"time"

	"github.com/compose-spec/compose-go/v2/types"
)

type Service struct {
	Restart   string
	Container Container
}

type Blueprint struct {
	types.ServiceConfig
}

type Container struct {
	Name      string
	Blueprint Blueprint
	LastKnown ContainerState
}

type ContainerState struct {
	ID      string
	Name    string
	State   string
	Updated time.Time
}
