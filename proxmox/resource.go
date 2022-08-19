package proxmox

import (
	"net/url"

	"github.com/lwch/runtime"
)

type Resource struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Node      string        `json:"node"`
	Type      ResourceType  `json:"type"`
	Status    Status        `json:"status"`  // Resource type dependent status
	Uptime    int           `json:"uptime"`  // Node uptime in seconds (when type in node,qemu,lxc).
	Cpu       float64       `json:"cpu"`     // CPU utilization (when type in node,qemu,lxc)
	Memory    runtime.Bytes `json:"mem"`     // Used memory in bytes (when type in node,qemu,lxc)
	Disk      runtime.Bytes `json:"disk"`    // Used disk space in bytes (when type in storage), used root image spave for VMs (type in qemu,lxc)
	MaxCpu    int           `json:"maxcpu"`  // Number of available CPUs (when type in node,qemu,lxc)
	MaxDisk   runtime.Bytes `json:"maxdisk"` // Storage size in bytes (when type in storage), root image size for VMs (type in qemu,lxc)
	MaxMemory runtime.Bytes `json:"maxmem"`  // Number of available memory in bytes (when type in node,qemu,lxc)
	// not in document
	DiskRead  runtime.Bytes `json:"diskread"`
	DiskWrite runtime.Bytes `json:"diskwrite"`
	NetIn     runtime.Bytes `json:"netin"`
	NetOut    runtime.Bytes `json:"netout"`
}

// Resources get resources
func (cli *Client) Resources(t ResourceType) ([]Resource, error) {
	var data struct {
		Data []Resource `json:"data"`
	}
	args := make(url.Values)
	if t != ResourceNone {
		args.Set("type", string(t))
	}
	err := cli.get("/cluster/resources", nil, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}
