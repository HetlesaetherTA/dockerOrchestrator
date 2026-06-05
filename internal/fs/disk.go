package fs

import (
	"hetlesaether.com/dockerOrchestrator/internal/domain"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	onlineDisks  []domain.Disk
	offlineDisks []domain.Disk
	disksOnce    sync.Once
)

func readDisksFromEnv() {
	disksOnce.Do(func() {
		slog.Info("Started reading Disks from .env")
		var tmp []domain.Disk

		for _, env := range os.Environ() {
			pair := strings.SplitN(env, "=", 2)
			if len(pair) < 2 {
				continue
			}

			key := pair[0]
			val := pair[1]

			if !strings.HasPrefix(strings.ToUpper(key), "DISK_") {
				continue
			}

			name := strings.TrimPrefix(strings.ToUpper(key), "DISK_")

			disk := domain.Disk{
				Name: strings.TrimSpace(name),
				UUID: strings.TrimSpace(val),
			}

			path := "/dev/disk/by-uuid/" + disk.UUID

			if !pathExists(path) {
				offlineDisks = append(offlineDisks, disk)
				slog.Warn("Couldn't find Disk's dev file. Disk is offline", "file not found", path, "key", key, "value", val)
				continue
			}

			slog.Info("/dev/ file found! Disk is online", "key", key, "value", val)
			tmp = append(tmp, disk)

		}
		onlineDisks = tmp
		slog.Info("Finished reading Disks from .env", "onlineDisks", onlineDisks, "offlineDisks", offlineDisks)
	})

}

func GetDisks() []domain.Disk {
	readDisksFromEnv()
	recheckDisks()

	snapshot := make([]domain.Disk, len(onlineDisks))
	copy(snapshot, onlineDisks)
	return snapshot
}

func GetDiskFromUUID(uuid string) domain.Disk {
	recheckDisks()

	for _, v := range onlineDisks {
		if v.UUID == uuid {
			return v
		}
	}
	return domain.Disk{}
}

func recheckDisks() {
	if len(offlineDisks) == 0 {
		return
	}

	slog.Debug("Starting periodic disk health recheck lifecycle")

	allDisks := append(onlineDisks, offlineDisks...)

	var nextOnline []domain.Disk
	var nextOffline []domain.Disk

	for _, disk := range allDisks {
		path := "/dev/disk/by-uuid/" + disk.UUID
		currentlyPhysical := pathExists(path)

		wasOnline := isDiskInSlice(disk, onlineDisks)

		if currentlyPhysical {
			nextOnline = append(nextOnline, disk)

			if !wasOnline {
				slog.Info("Disk status changed: Drive is now ONLINE",
					"disk", disk.Name,
					"uuid", disk.UUID,
					"node", path,
				)
			}
		} else {
			nextOffline = append(nextOffline, disk)

			if wasOnline {
				slog.Error("Disk status changed: Drive has gone OFFLINE",
					"disk", disk.Name,
					"uuid", disk.UUID,
					"expected_node", path,
				)
			}
		}
		if len(offlineDisks) != 0 {
			slog.Info("Finished rechecking disk, you have offline disks!", "offlineDisks", offlineDisks)
		}
	}

	onlineDisks = nextOnline
	offlineDisks = nextOffline
}

func isDiskInSlice(target domain.Disk, list []domain.Disk) bool {
	for _, d := range list {
		if d.UUID == target.UUID {
			return true
		}
	}
	return false
}

func pathExists(path string) bool {
	_, err := os.Stat(path)

	if err == nil {
		return true
	}

	return false
}
