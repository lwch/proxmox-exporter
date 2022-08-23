package exporter

import (
	"exporter/proxmox"

	"github.com/lwch/logging"
	"github.com/prometheus/client_golang/prometheus"
)

type vmExporter struct {
	parent *nodeExporter

	// stats
	uptime *prometheus.GaugeVec
	info   *prometheus.GaugeVec
}

func newVmExporter(parent *nodeExporter) *vmExporter {
	exp := &vmExporter{parent: parent}
	exp.build()
	return exp
}

func (exp *vmExporter) build() {
	const namespace = "vm"
	constLabels := prometheus.Labels{"node_name": exp.parent.name}
	labels := []string{"vm_id", "vm_name"}

	// info
	exp.uptime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "uptime",
		Help:        "vm uptime",
		ConstLabels: constLabels,
	}, labels)
	exp.info = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "info",
		Help: `vm info, labels:
type: lxc or qemu`,
	}, append(labels, "type"))
}

func (exp *vmExporter) Describe(ch chan<- *prometheus.Desc) {
	// info
	exp.uptime.Describe(ch)
	exp.info.Describe(ch)
}

func (exp *vmExporter) Collect(ch chan<- prometheus.Metric) {
	exp.updateStatus()

	// info
	exp.uptime.Collect(ch)
	exp.info.Collect(ch)
}

func (exp *vmExporter) updateStatus() {
	vms, err := exp.parent.parent.cli.ClusterResources(proxmox.ResourceVM)
	if err != nil {
		logging.Error("get vm resource: %v", err)
		return
	}
	merge := func(a, b prometheus.Labels) prometheus.Labels {
		ret := make(prometheus.Labels)
		for k, v := range a {
			ret[k] = v
		}
		for k, v := range b {
			ret[k] = v
		}
		return ret
	}
	for _, vm := range vms {
		labels := prometheus.Labels{
			"vm_id":   vm.ID,
			"vm_name": vm.Name,
		}
		// info
		exp.uptime.With(labels).Set(float64(vm.Uptime))
		exp.info.With(merge(labels, prometheus.Labels{"type": string(vm.Type)})).Set(1)
	}
}
