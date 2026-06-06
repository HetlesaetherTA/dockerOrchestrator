package infra

import (
	"hetlesaether.com/dockerOrchestrator/internal/domain"
	"hetlesaether.com/dockerOrchestrator/internal/fs"
)

type infraRepository struct {
	Disks    map[string]domain.Disk
	Mounts   []domain.Mount
	Networks map[string]domain.Network
}

func New() infraRepository {
	return infraRepository{}
}

func NewFromYAML(path string) (infraRepository, error) {
	ir := New()
	return ir, fs.ConstructFromYAML(&ir, path)
}
