package infra

import (
	"fmt"
	"maps"

	"hetlesaether.com/dockerOrchestrator/internal/domain"
)

func (ir *infraRepository) AddDisk(name string, disk domain.Disk) error {
	for k, v := range ir.Disks {
		if k == name {
			return fmt.Errorf("disk with name %q already exists", name)
		}

		if v == disk {
			return fmt.Errorf("disk is duplicate of disk with name %q", k)
		}
	}

	ir.Disks[name] = disk

	return nil
}

func (ir *infraRepository) RemoveDiskByName(name string) {
	delete(ir.Disks, name)
}

func (ir *infraRepository) RemoveDisk(disk domain.Disk) {
	maps.DeleteFunc(ir.Disks, func(_ string, v domain.Disk) bool {
		return disk == v
	})
}

func (ir *infraRepository) GetDisks() map[string]domain.Disk {
	return ir.Disks
}
