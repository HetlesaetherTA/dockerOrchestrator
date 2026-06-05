package hardware

import (
	"hetlesaether.com/dockerOrchestrator/internal/domain"
)

type hardwareRepository struct {
	Disks  []domain.Disk
	Mounts []domain.Mount
}

func CreateRepository() hardwareRepository {
	return hardwareRepository{}
}
