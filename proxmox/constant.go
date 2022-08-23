package proxmox

import (
	"bytes"
)

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

// StorageContent storage content type
type StorageContent string

const (
	ContentTemplate StorageContent = "vztmpl"   // lxc container template
	ContentIso      StorageContent = "iso"      // iso image
	ContentBackup   StorageContent = "backup"   // backup
	ContentSnippets StorageContent = "snippets" // snippets
	ContentRootDir  StorageContent = "rootdir"  // container rootdir
	ContentImages   StorageContent = "images"   // kvm disk images
)

// StorageContents storage contents
type StorageContents []StorageContent

func (content *StorageContents) UnmarshalJSON(value []byte) error {
	value = bytes.TrimPrefix(value, []byte{'"'})
	value = bytes.TrimSuffix(value, []byte{'"'})
	for _, v := range bytes.Split(value, []byte{','}) {
		*content = append(*content, StorageContent(v))
	}
	return nil
}

type StorageType string
