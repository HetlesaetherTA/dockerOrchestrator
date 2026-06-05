package hardware

import (
	"errors"
	"slices"

	"hetlesaether.com/dockerOrchestrator/internal/domain"
)

func (hr *hardwareRepository) AddMount(mount domain.Mount) error {
	if slices.Contains(hr.Mounts, mount) {
		return errors.New("Mount already registered")
	}

	hr.Mounts = append(hr.Mounts, mount)

	return nil
}

func (hr *hardwareRepository) RemoveMount(mount domain.Mount) {
	hr.Mounts = slices.DeleteFunc(hr.Mounts, func(m domain.Mount) bool {
		return mount == m
	})
}

func (hr *hardwareRepository) GetMounts() []domain.Mount {
	return hr.Mounts
}
