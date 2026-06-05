package hardware

import (
	"errors"
	"slices"

	"hetlesaether.com/dockerOrchestrator/internal/domain"
)

func (hr *hardwareRepository) AddDisk(disk domain.Disk) error {
	if slices.Contains(hr.Disks, disk) {
		return errors.New("Disk already registered")
	}

	hr.Disks = append(hr.Disks, disk)

	return nil
}

func (hr *hardwareRepository) RemoveDisk(disk domain.Disk) {
	hr.Disks = slices.DeleteFunc(hr.Disks, func(d domain.Disk) bool {
		return disk == d
	})
}

func (hr *hardwareRepository) GetDisks() []domain.Disk {
	return hr.Disks
}
