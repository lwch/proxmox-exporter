package exporter

import (
	"exporter/proxmox"
	"sync"
	"time"

	"github.com/lwch/logging"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter exporter struct
type Exporter struct {
	sync.RWMutex
	cli   *proxmox.Client
	nodes map[string]*nodeExporter
}

// New create exporter
func New(cli *proxmox.Client) *Exporter {
	exp := &Exporter{
		cli:   cli,
		nodes: make(map[string]*nodeExporter),
	}
	go exp.collectNodes()
	return exp
}

func (exp *Exporter) collectNodes() {
	tk := time.NewTicker(10 * time.Second)
	get := func() {
		nodes, err := exp.cli.ClusterResources(proxmox.ResourceNode)
		if err != nil {
			logging.Error("can not get node resource: %v", err)
			return
		}
		exp.Lock()
		defer exp.Unlock()
		for _, node := range nodes {
			if _, ok := exp.nodes[node.Node]; !ok {
				exp.nodes[node.Node] = newNodeExporter(exp, node.Node)
			}
		}
	}
	for {
		get()
		<-tk.C
	}
}

func (exp *Exporter) Describe(ch chan<- *prometheus.Desc) {
	nodes := make([]*nodeExporter, 0, len(exp.nodes))
	exp.RLock()
	for _, node := range exp.nodes {
		nodes = append(nodes, node)
	}
	exp.RUnlock()

	for _, node := range nodes {
		node.Describe(ch)
	}
}

func (exp *Exporter) Collect(ch chan<- prometheus.Metric) {
	nodes := make([]*nodeExporter, 0, len(exp.nodes))
	exp.RLock()
	for _, node := range exp.nodes {
		nodes = append(nodes, node)
	}
	exp.RUnlock()

	for _, node := range nodes {
		node.Collect(ch)
	}
}
