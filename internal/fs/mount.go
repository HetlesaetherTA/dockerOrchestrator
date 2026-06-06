package fs

//
// import (
// 	"errors"
// 	"fmt"
// 	"golang.org/x/sys/unix"
// 	"hetlesaether.com/dockerOrchestrator/internal/domain"
// 	"log/slog"
// 	"os"
// 	"strings"
// )

// func MountDisk(mount *domain.Mount) error {
// 	disk := GetDiskFromUUID(mount.UUID)
//
// 	if disk.UUID == "" {
// 		return fmt.Errorf("Disk not online:  %s", mount.UUID)
// 	}
//
// 	devicePath := fmt.Sprintf("/dev/disk/by-uuid/%s", mount.UUID)
//
// 	if _, err := os.Stat(devicePath); err != nil {
// 		if errors.Is(err, os.ErrNotExist) {
// 			return fmt.Errorf("hardware partition missing: disk %s with UUID [%s] is not connected", disk.Name, disk.UUID)
// 		}
// 		return fmt.Errorf("failed to check hardware block device state: %w", err)
// 	}
//
// 	if err := os.MkdirAll(mount.MountPath, 0755); err != nil {
// 		return fmt.Errorf("failed to create mount target directory directory %s: %w", mount.MountPath, err)
// 	}
//
// 	alreadyMounted, err := IsMounted(mount.MountPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to verify active mount tables: %w", err)
// 	}
// 	if alreadyMounted {
// 		slog.Info("Storage gate clear: Disk is already mounted", "target", mount.MountPath)
// 		mount.Mounted = true
// 		return nil
// 	}
//
// 	err = unix.Mount(devicePath, mount.MountPath, "ext4", unix.MS_MGC_VAL, "")
// 	if err != nil {
// 		return fmt.Errorf("kernel system mount call failed for dev [%s]: %w", devicePath, err)
// 	}
//
// 	slog.Info("Storage drive mounted successfully", "disk", disk.Name, "path", mount.MountPath)
// 	mount.Mounted = true
// 	return nil
// }
//
// func IsMounted(path string) (bool, error) {
// 	data, err := os.ReadFile("/proc/mounts")
// 	if err != nil {
// 		return false, err
// 	}
// 	return strings.Contains(string(data), path), nil
// }
