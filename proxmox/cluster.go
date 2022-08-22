package proxmox

import (
	"net/url"

	"github.com/lwch/runtime"
)

// Resource resource info
type Resource struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Node      string         `json:"node"`
	Type      ResourceType   `json:"type"`
	Status    ResourceStatus `json:"status"`  // Resource type dependent status
	Uptime    int            `json:"uptime"`  // Node uptime in seconds (when type in node,qemu,lxc).
	Cpu       float64        `json:"cpu"`     // CPU utilization (when type in node,qemu,lxc)
	Memory    runtime.Bytes  `json:"mem"`     // Used memory in bytes (when type in node,qemu,lxc)
	Disk      runtime.Bytes  `json:"disk"`    // Used disk space in bytes (when type in storage), used root image spave for VMs (type in qemu,lxc)
	MaxCpu    int            `json:"maxcpu"`  // Number of available CPUs (when type in node,qemu,lxc)
	MaxDisk   runtime.Bytes  `json:"maxdisk"` // Storage size in bytes (when type in storage), root image size for VMs (type in qemu,lxc)
	MaxMemory runtime.Bytes  `json:"maxmem"`  // Number of available memory in bytes (when type in node,qemu,lxc)
	// not in document
	DiskRead  runtime.Bytes `json:"diskread"`
	DiskWrite runtime.Bytes `json:"diskwrite"`
	NetIn     runtime.Bytes `json:"netin"`
	NetOut    runtime.Bytes `json:"netout"`
}

// ClusterResources get resources of cluster
func (cli *Client) ClusterResources(t ResourceType) ([]Resource, error) {
	var data struct {
		Data []Resource `json:"data"`
	}
	args := make(url.Values)
	if t != ResourceNone {
		args.Set("type", string(t))
	}
	err := cli.get("/cluster/resources", args, &data)
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

// ClusterTasks get tasks of cluster
func (cli *Client) ClusterTasks() ([]Task, error) {
	var data struct {
		Data []Task `json:"data"`
	}
	err := cli.get("/cluster/tasks", nil, &data)
	return data.Data, err
}

// ClusterStatus cluster status
type ClusterStatus struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Type    ResourceType `json:"type"`    // cluster or node
	Nodes   int          `json:"nodes"`   // [cluster] Nodes count, including offline nodes
	Quorate int          `json:"quorate"` // [cluster] Indicates if there is a majority of nodes online to make decisions
	Version int          `json:"version"` // [cluster] Current version of the corosync configuration file
	Ip      string       `json:"ip"`      // [node] IP address
	Level   string       `json:"level"`   // [node] Proxmox VE Subscription level
	Local   int          `json:"local"`   // [node] Indicates if this is the responding node
	NodeID  int          `json:"nodeid"`  // [node] ID of the node from the corosync configuration
	Online  int          `json:"online"`  // [node] Indicates if the node is online or offline
}

// ClusterStatus get cluster status
func (cli *Client) ClusterStatus() ([]ClusterStatus, error) {
	var data struct {
		Data []ClusterStatus `json:"data"`
	}
	err := cli.get("/cluster/status", nil, &data)
	return data.Data, err
}
