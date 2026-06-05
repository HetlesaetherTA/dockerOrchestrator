package domain

type Mount struct {
	UUID      string // Domain.Disk to mount
	MountPath string // Absolute path to mount to
	Mounted   bool   // Is Disk currently mounted?
}

// Validated and constructed from .env: DISK_{Name}={UUID}
type Disk struct {
	Name string // DISK_{Name} from .env
	UUID string // /dev/disk/by-uuid
}
