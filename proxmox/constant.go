package proxmox

// ResourceType resource type
type ResourceType string

const (
	ResourceNone    ResourceType = "none"
	ResourceVM      ResourceType = "vm"
	ResourceStorage ResourceType = "storage"
	ResourceNode    ResourceType = "node"
	ResourceSdn     ResourceType = "sdn"
	ResourcePool    ResourceType = "pool"
	ResourceQemu    ResourceType = "qemu"
	ResourceLxc     ResourceType = "lxc"
	ResourceOpenvz  ResourceType = "openvz"
)

// Status resource status
type Status string

const (
	StatusRunning   Status = "running"
	StatusStopped   Status = "stopped"
	StatusAvailable Status = "available"
)
