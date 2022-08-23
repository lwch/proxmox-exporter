package exporter

import (
	"exporter/proxmox"
	"fmt"
	"sync"
	"time"

	"github.com/lwch/logging"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter exporter struct
type Exporter struct {
	sync.RWMutex
	cli  *proxmox.Client
	node *nodeExporter
}

func getNode(cli *proxmox.Client, exp *Exporter) *nodeExporter {
	const maxCount = 10
	for i := 0; i < maxCount; i++ {
		status, err := cli.ClusterStatus()
		if err != nil {
			logging.Warning("get cluster/status %d times: %v", i+1, err)
			time.Sleep(time.Second)
			continue
		}
		for _, st := range status {
			if st.Type != proxmox.ResourceNode {
				continue
			}
			if st.Local != 0 {
				logging.Info("current node name: %s", st.Name)
				return newNodeExporter(exp, st.Name)
			}
		}
		logging.Warning("get cluster/status %d times...", i+1)
		time.Sleep(time.Second)
	}
	panic(fmt.Sprintf("get cluster/status more than %d times", maxCount))
}

// New create exporter
func New(cli *proxmox.Client) *Exporter {
	exp := &Exporter{cli: cli}
	exp.node = getNode(cli, exp)
	return exp
}

func (exp *Exporter) Describe(ch chan<- *prometheus.Desc) {
	exp.node.Describe(ch)
}

func (exp *Exporter) Collect(ch chan<- prometheus.Metric) {
	exp.node.Collect(ch)
}
