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
		_ = nodes
		exp.Lock()
		defer exp.Unlock()
	}
	for {
		get()
		<-tk.C
	}
}

func (exp *Exporter) Describe(ch chan<- *prometheus.Desc) {
}

func (exp *Exporter) Collect(ch chan<- prometheus.Metric) {
}
