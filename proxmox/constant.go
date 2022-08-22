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
	ResourceCluster ResourceType = "cluster"
)

// Status resource status
type ResourceStatus string

const (
	StatusRunning   ResourceStatus = "running"
	StatusStopped   ResourceStatus = "stopped"
	StatusAvailable ResourceStatus = "available"
)
