package proxmox

import (
	"fmt"
	"net/url"

	"github.com/lwch/runtime"
)

type usage struct {
	Used  runtime.Bytes `json:"used"`  // used
	Free  runtime.Bytes `json:"free"`  // free
	Total runtime.Bytes `json:"total"` // total
}

// NodeStatus node status
type NodeStatus struct {
	Uptime      int      `json:"uptime"`  // Node uptime in seconds (when type in node,qemu,lxc).
	Cpu         float64  `json:"cpu"`     // CPU utilization (when type in node,qemu,lxc)
	LoadAverage []string `json:"loadavg"` // CPU load average
	CpuInfo     struct {
		Model         string `json:"model"`   // CPU model
		Flags         string `json:"flags"`   // CPU flags
		Sockets       uint64 `json:"sockets"` // CPU sockets
		Cores         uint64 `json:"cores"`   // CPU cores
		Threads       uint64 `json:"cpus"`    // CPU threads
		Frequency     string `json:"mhz"`     // CPU frequency
		UserFrequency int    `json:"user_hz"` // CPU limit frequency?
	} `json:"cpuinfo"`
	Memory        usage   `json:"memory"`     // Physical memory info
	Swap          usage   `json:"swap"`       // Swap memory info
	RootFs        usage   `json:"rootfs"`     // rootfs usage
	Idle          float64 `json:"idle"`       // IO idle?
	Wait          float64 `json:"wait"`       // IO wait ratio
	KernelVersion string  `json:"kversion"`   // kernel version
	PveVersion    string  `json:"pveversion"` // proxmox version
}

// NodeStatus get node status
func (cli *Client) NodeStatus(name string) (NodeStatus, error) {
	var data struct {
		Data NodeStatus `json:"data"`
	}
	err := cli.get(fmt.Sprintf("/nodes/%s/status", name), nil, &data)
	return data.Data, err
}

// Task task info
type Task struct {
	ID      string `json:"id"`
	UnionID string `json:"upid"`      // union process id
	User    string `json:"user"`      // create user
	Status  string `json:"status"`    // task status, empty means pending
	Node    string `json:"node"`      // task on node
	Type    string `json:"type"`      // task type
	Pid     int    `json:"pid"`       // process id
	PStart  int    `json:"pstart"`    // ??
	Start   int64  `json:"starttime"` // begin time
	End     int64  `json:"endtime"`   // end time
}

// TaskInfo response task info
type TaskInfo struct {
	Total int    `json:"total"`
	Tasks []Task `json:"data"`
}

// Tasks get task list on node
func (cli *Client) NodeTasks(name string, start, limit int) (TaskInfo, error) {
	if start < 0 {
		start = 0
	}
	if limit < 0 {
		limit = 50
	}
	args := make(url.Values)
	args.Set("start", fmt.Sprintf("%d", start))
	args.Set("limit", fmt.Sprintf("%d", limit))
	var info TaskInfo
	err := cli.get(fmt.Sprintf("/nodes/%s/tasks", name), args, &info)
	return info, err
}

// NodeStorage node storage
type NodeStorage struct {
	Content   StorageContents `json:"content"`       // Allowed storage content types
	Storage   string          `json:"storage"`       // The storage identifier
	Type      string          `json:"type"`          // Storage type, TODO: struct
	Active    int             `json:"active"`        // Set when storage is accessible
	Enabled   int             `json:"enabled"`       // Set when storage is enabled (not disabled)
	Shared    int             `json:"shared"`        // Shared flag from storage configuration
	Used      runtime.Bytes   `json:"used"`          // Used storage space in bytes
	Available runtime.Bytes   `json:"avail"`         // Available storage space in bytes
	Total     runtime.Bytes   `json:"total"`         // Total storage space in bytes
	Ratio     float64         `json:"used_fraction"` // Used fraction (used/total)
}

// NodeStorage get storage list of node
func (cli *Client) NodeStorage(name string) ([]NodeStorage, error) {
	var data struct {
		Data []NodeStorage `json:"data"`
	}
	err := cli.get(fmt.Sprintf("/nodes/%s/storage", name), nil, &data)
	return data.Data, err
}

// NodeRrdData node rrddata
type NodeRrdData struct {
	Time        int64   `json:"time"`
	Cpu         float64 `json:"cpu"`
	MaxCpu      float64 `json:"maxcpu"`
	LoadAverage float64 `json:"loadavg"`
	MemoryUsed  float64 `json:"memused"`
	MemoryTotal float64 `json:"memtotal"`
	SwapUsed    float64 `json:"swapused"`
	SwapTotal   float64 `json:"swaptotal"`
	RootUsed    float64 `json:"rootused"`
	RootTotal   float64 `json:"roottotal"`
	NetIn       float64 `json:"netin"`
	NetOut      float64 `json:"netout"`
	IoWait      float64 `json:"iowait"`
}

// NodeRrdData get rrddata of node from hour
func (cli *Client) NodeRrdData(name string) ([]NodeRrdData, error) {
	args := make(url.Values)
	args.Set("timeframe", "hour")
	var data struct {
		Data []NodeRrdData `json:"data"`
	}
	err := cli.get(fmt.Sprintf("/nodes/%s/rrddata", name), args, &data)
	return data.Data, err
}
