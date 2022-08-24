package exporter

import (
	"exporter/proxmox"
	"fmt"

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
	// memory
	memoryUsed  *prometheus.GaugeVec
	memoryTotal *prometheus.GaugeVec
	// disk
	diskUsed  *prometheus.GaugeVec
	diskTotal *prometheus.GaugeVec
	diskRead  *prometheus.GaugeVec
	diskWrite *prometheus.GaugeVec
	// network
	netin  *prometheus.GaugeVec
	netout *prometheus.GaugeVec
}

func newVmExporter(parent *nodeExporter) *vmExporter {
	exp := &vmExporter{parent: parent}
	exp.build()
	return exp
}

func (exp *vmExporter) build() {
	const namespace = "vm"
	constLabels := prometheus.Labels{"node_name": exp.parent.name}
	labels := []string{"vm_id", "vm_name", "vm_type", "vm_status"}

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
uptime: uptime
core: cpu cores
memory: max memory bytes
disk: max disk bytes`,
		ConstLabels: constLabels,
	}, append(labels, "uptime", "core", "memory", "disk"))
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
	// memory
	exp.memoryUsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "memory_used",
		Help:        "vm memory used bytes",
		ConstLabels: constLabels,
	}, labels)
	exp.memoryTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "memory_total",
		Help:        "vm memory total bytes",
		ConstLabels: constLabels,
	}, labels)
	// disk
	exp.diskUsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "disk_used",
		Help:        "vm disk used bytes",
		ConstLabels: constLabels,
	}, labels)
	exp.diskTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "disk_total",
		Help:        "vm disk total bytes",
		ConstLabels: constLabels,
	}, labels)
	exp.diskRead = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "disk_read",
		Help:        "vm disk readen bytes",
		ConstLabels: constLabels,
	}, labels)
	exp.diskWrite = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "disk_write",
		Help:        "vm disk written bytes",
		ConstLabels: constLabels,
	}, labels)
	// network
	exp.netin = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "netin",
		Help:        "vm network received bytes",
		ConstLabels: constLabels,
	}, labels)
	exp.netout = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "netout",
		Help:        "vm network sent bytes",
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
	// memory
	exp.memoryUsed.Describe(ch)
	exp.memoryTotal.Describe(ch)
	// disk
	exp.diskUsed.Describe(ch)
	exp.diskTotal.Describe(ch)
	exp.diskRead.Describe(ch)
	exp.diskWrite.Describe(ch)
	// network
	exp.netin.Describe(ch)
	exp.netout.Describe(ch)
}

func (exp *vmExporter) Collect(ch chan<- prometheus.Metric) {
	exp.updateStatus()

	// info
	exp.uptime.Collect(ch)
	exp.info.Collect(ch)
	// cpu
	exp.cpuUsage.Collect(ch)
	exp.cpuTotal.Collect(ch)
	// memory
	exp.memoryUsed.Collect(ch)
	exp.memoryTotal.Collect(ch)
	// disk
	exp.diskUsed.Collect(ch)
	exp.diskTotal.Collect(ch)
	exp.diskRead.Collect(ch)
	exp.diskWrite.Collect(ch)
	// network
	exp.netin.Collect(ch)
	exp.netout.Collect(ch)
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
	exp.info.Reset()
	for _, vm := range vms {
		labels := prometheus.Labels{
			"vm_id":     vm.ID,
			"vm_name":   vm.Name,
			"vm_type":   string(vm.Type),
			"vm_status": string(vm.Status),
		}
		// info
		exp.uptime.With(labels).Set(float64(vm.Uptime))
		exp.info.With(merge(labels, prometheus.Labels{
			"uptime": fmt.Sprintf("%d", vm.Uptime),
			"core":   fmt.Sprintf("%d", vm.MaxCpu),
			"memory": fmt.Sprintf("%d", vm.MaxMemory),
			"disk":   fmt.Sprintf("%d", vm.MaxDisk),
		})).Inc()
		// cpu
		exp.cpuUsage.With(labels).Set(vm.Cpu)
		exp.cpuTotal.With(labels).Set(float64(vm.MaxCpu))
		// memory
		exp.memoryUsed.With(labels).Set(float64(vm.Memory))
		exp.memoryTotal.With(labels).Set(float64(vm.MaxMemory))
		// disk
		exp.diskUsed.With(labels).Set(float64(vm.Disk))
		exp.diskTotal.With(labels).Set(float64(vm.MaxDisk))
		exp.diskRead.With(labels).Set(float64(vm.DiskRead))
		exp.diskWrite.With(labels).Set(float64(vm.DiskWrite))
		// network
		exp.netin.With(labels).Set(float64(vm.NetIn))
		exp.netout.With(labels).Set(float64(vm.NetOut))
	}
}
