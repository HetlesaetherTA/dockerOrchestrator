package infra

import (
	"fmt"
	"slices"

	"hetlesaether.com/dockerOrchestrator/internal/domain"
)

func (ir *infraRepository) AddMount(mount domain.Mount) error {
	for _, v := range ir.Mounts {
		if v.Dst == mount.Dst {
			return fmt.Errorf("mount destination %q already mounted by %q", mount.Dst, v.Src)
		}
	}

	ir.Mounts = append(ir.Mounts, mount)

	return nil
}

func (ir *infraRepository) RemoveMount(mount domain.Mount) {
	ir.Mounts = slices.DeleteFunc(ir.Mounts, func(e domain.Mount) bool {
		return mount == e
	})
}

func (ir *infraRepository) GetMounts() []domain.Mount {
	return ir.Mounts
}
