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
	// cpu
	cpuUsage *prometheus.GaugeVec
	cpuTotal *prometheus.GaugeVec
}

func newVmExporter(parent *nodeExporter) *vmExporter {
	exp := &vmExporter{parent: parent}
	exp.build()
	return exp
}

func (exp *vmExporter) build() {
	const namespace = "vm"
	constLabels := prometheus.Labels{"node_name": exp.parent.name}
	labels := []string{"vm_id", "vm_name", "vm_status"}

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
		ConstLabels: constLabels,
	}, append(labels, "type"))
	// cpu
	exp.cpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "cpu_usage",
		Help:        "vm cpu usage ratio(precent)",
		ConstLabels: constLabels,
	}, labels)
	exp.cpuTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "cpu_total",
		Help:        "vm max cpu core count",
		ConstLabels: constLabels,
	}, labels)
}

func (exp *vmExporter) Describe(ch chan<- *prometheus.Desc) {
	// info
	exp.uptime.Describe(ch)
	exp.info.Describe(ch)
	// cpu
	exp.cpuUsage.Describe(ch)
	exp.cpuTotal.Describe(ch)
}

func (exp *vmExporter) Collect(ch chan<- prometheus.Metric) {
	exp.updateStatus()

	// info
	exp.uptime.Collect(ch)
	exp.info.Collect(ch)
	// cpu
	exp.cpuUsage.Collect(ch)
	exp.cpuTotal.Collect(ch)
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
			"vm_id":     vm.ID,
			"vm_name":   vm.Name,
			"vm_status": string(vm.Status),
		}
		// info
		exp.uptime.With(labels).Set(float64(vm.Uptime))
		exp.info.With(merge(labels, prometheus.Labels{"type": string(vm.Type)})).Set(1)
		// cpu
		exp.cpuUsage.With(labels).Set(vm.Cpu)
		exp.cpuTotal.With(labels).Set(float64(vm.MaxCpu))
	}
}
