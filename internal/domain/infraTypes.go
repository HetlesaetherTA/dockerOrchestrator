package domain

type InfraRepositoryYAMLResult struct {
	Disks    map[string]Disk    `yaml:"disks"`
	Mounts   []Mount            `yaml:"mounts"`
	Networks map[string]Network `yaml:"networks"`
}

type Disk struct {
	UUID string `yaml:"uuid"`
	RO   bool   `yaml:"ro"`
	AO   bool   `yaml:"ao"`
}

type Mount struct {
	Src string `yaml:"source"`
	Dst string `yaml:"destination"`
}

type Network struct {
	Driver      string `yaml:"driver"`
	Internal    bool   `yaml:"internal"`
	External    bool   `yaml:"external"`
	Attachable  bool   `yaml:"attachable"`
	Ipv4        bool   `yaml:"ipv4"`
	Ipv6        bool   `yaml:"ipv6"`
	Description string `yaml:"description"`

	Label string
}
