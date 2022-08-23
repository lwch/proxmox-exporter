package exporter

import (
	"exporter/proxmox"
	"fmt"
	"strconv"

	"github.com/lwch/logging"
	"github.com/prometheus/client_golang/prometheus"
)

type nodeExporter struct {
	parent *Exporter
	name   string
	vm     *vmExporter

	// stats
	upTime         prometheus.Gauge
	info           *prometheus.GaugeVec
	cpuUsage       prometheus.Gauge
	cpuLoadAverage *prometheus.GaugeVec
	memoryUsed     prometheus.Gauge
	memoryFree     prometheus.Gauge
	memoryTotal    prometheus.Gauge
	swapUsed       prometheus.Gauge
	swapFree       prometheus.Gauge
	swapTotal      prometheus.Gauge
	rootfsUsed     prometheus.Gauge
	rootfsFree     prometheus.Gauge
	rootfsTotal    prometheus.Gauge
	ioWait         prometheus.Gauge
	storageInfo    *prometheus.GaugeVec
}

func newNodeExporter(parent *Exporter, name string) *nodeExporter {
	exp := &nodeExporter{
		parent: parent,
		name:   name,
	}
	exp.vm = newVmExporter(exp)
	exp.build()
	return exp
}

func (exp *nodeExporter) build() {
	const namespace = "node"
	labels := prometheus.Labels{"node_name": exp.name}

	// online
	exp.upTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "uptime",
		Help:        "node uptime",
		ConstLabels: labels,
	})
	exp.info = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "info",
		Help: `node info, labels:
model: cpu model
sockets: cpu sockets count
cores: cpu cores
threads: cpu threads
mhz: cpu frequency
kernel_version: linux kernel version
pve_version: proxmox version`,
		ConstLabels: labels,
	}, []string{"model", "sockets", "cores", "threads", "mhz",
		"kernel_version", "pve_version"})

	// cpu
	exp.cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "cpu_usage",
		Help:        "node cpu usage ratio(precent)",
		ConstLabels: labels,
	})
	exp.cpuLoadAverage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "cpu_loadavg",
		Help:        "node cpu load average",
		ConstLabels: labels,
	}, []string{"minute"})

	// memory
	exp.memoryUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "memory_used",
		Help:        "used memory bytes of this node",
		ConstLabels: labels,
	})
	exp.memoryFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "memory_free",
		Help:        "free memory bytes of this node",
		ConstLabels: labels,
	})
	exp.memoryTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "memory_total",
		Help:        "total memory bytes of this node",
		ConstLabels: labels,
	})

	// swap
	exp.swapUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "swap_used",
		Help:        "used swap bytes of this node",
		ConstLabels: labels,
	})
	exp.swapFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "swap_free",
		Help:        "free swap bytes of this node",
		ConstLabels: labels,
	})
	exp.swapTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "swap_total",
		Help:        "total swap bytes of this node",
		ConstLabels: labels,
	})

	// rootfs
	exp.rootfsUsed = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "rootfs_used",
		Help:        "used rootfs bytes of this node",
		ConstLabels: labels,
	})
	exp.rootfsFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "rootfs_free",
		Help:        "free rootfs bytes of this node",
		ConstLabels: labels,
	})
	exp.rootfsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "rootfs_total",
		Help:        "total rootfs bytes of this node",
		ConstLabels: labels,
	})
	// disk
	exp.ioWait = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "iowait",
		Help:        "node iowait ratio(precent)",
		ConstLabels: labels,
	})
	exp.storageInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "storage_info",
		Help: `node storage info, labels:
content_*: allowed content type
storage: storage name
type: storage type`,
		ConstLabels: labels,
	}, []string{
		"content_vztmpl",
		"content_iso",
		"content_backup",
		"content_snippets",
		"content_rootdir",
		"content_images",
		"storage",
		"type",
	})
}

func (exp *nodeExporter) Describe(ch chan<- *prometheus.Desc) {
	// online
	exp.upTime.Describe(ch)
	exp.info.Describe(ch)
	// cpu
	exp.cpuUsage.Describe(ch)
	exp.cpuLoadAverage.Describe(ch)
	// memory
	exp.memoryUsed.Describe(ch)
	exp.memoryFree.Describe(ch)
	exp.memoryTotal.Describe(ch)
	// swap
	exp.swapUsed.Describe(ch)
	exp.swapFree.Describe(ch)
	exp.swapTotal.Describe(ch)
	// rootfs
	exp.rootfsUsed.Describe(ch)
	exp.rootfsFree.Describe(ch)
	exp.rootfsTotal.Describe(ch)
	// disk
	exp.ioWait.Describe(ch)
	exp.storageInfo.Describe(ch)

	// vm describe
	exp.vm.Describe(ch)
}

func (exp *nodeExporter) Collect(ch chan<- prometheus.Metric) {
	// collect values
	exp.updateStatus()

	// online
	exp.upTime.Collect(ch)
	exp.info.Collect(ch)
	// cpu
	exp.cpuUsage.Collect(ch)
	exp.cpuLoadAverage.Collect(ch)
	// memory
	exp.memoryUsed.Collect(ch)
	exp.memoryFree.Collect(ch)
	exp.memoryTotal.Collect(ch)
	// swap
	exp.swapUsed.Collect(ch)
	exp.swapFree.Collect(ch)
	exp.swapTotal.Collect(ch)
	// rootfs
	exp.rootfsUsed.Collect(ch)
	exp.rootfsFree.Collect(ch)
	exp.rootfsTotal.Collect(ch)
	// disk
	exp.ioWait.Collect(ch)
	exp.storageInfo.Collect(ch)

	// vm collect
	exp.vm.Collect(ch)
}

func (exp *nodeExporter) updateStatus() {
	status, err := exp.parent.cli.NodeStatus(exp.name)
	if err != nil {
		logging.Error("get node [%s] status: %v", exp.name, err)
		return
	}
	exp.updateInfo(status)
	exp.updateCpu(status)
	exp.updateMemory(status)
	exp.updateDisk(status)
}

func (exp *nodeExporter) updateInfo(status proxmox.NodeStatus) {
	exp.upTime.Set(float64(status.Uptime))
	exp.info.With(prometheus.Labels{
		"model":          status.CpuInfo.Model,
		"sockets":        fmt.Sprintf("%d", status.CpuInfo.Sockets),
		"cores":          fmt.Sprintf("%d", status.CpuInfo.Cores),
		"threads":        fmt.Sprintf("%d", status.CpuInfo.Threads),
		"mhz":            status.CpuInfo.Frequency,
		"kernel_version": status.KernelVersion,
		"pve_version":    status.PveVersion,
	}).Set(1)
}

func (exp *nodeExporter) updateCpu(status proxmox.NodeStatus) {
	exp.cpuUsage.Set(status.Cpu * 100.)

	loadAvg := make([]float64, len(status.LoadAverage))
	for i, v := range status.LoadAverage {
		n, _ := strconv.ParseFloat(v, 64)
		loadAvg[i] = n
	}
	exp.cpuLoadAverage.With(prometheus.Labels{"minute": "1"}).Set(loadAvg[0])
	exp.cpuLoadAverage.With(prometheus.Labels{"minute": "5"}).Set(loadAvg[1])
	exp.cpuLoadAverage.With(prometheus.Labels{"minute": "15"}).Set(loadAvg[2])
}

func (exp *nodeExporter) updateMemory(status proxmox.NodeStatus) {
	// memory
	exp.memoryUsed.Set(float64(status.Memory.Used))
	exp.memoryFree.Set(float64(status.Memory.Free))
	exp.memoryTotal.Set(float64(status.Memory.Total))
	// swap
	exp.swapUsed.Set(float64(status.Swap.Used))
	exp.swapFree.Set(float64(status.Swap.Free))
	exp.swapTotal.Set(float64(status.Swap.Total))
}

func (exp *nodeExporter) updateDisk(status proxmox.NodeStatus) {
	// rootfs
	exp.rootfsUsed.Set(float64(status.RootFs.Used))
	exp.rootfsFree.Set(float64(status.RootFs.Free))
	exp.rootfsTotal.Set(float64(status.RootFs.Total))
	// disk
	exp.ioWait.Set(status.Wait * 100.)
	exp.updateStorage()
}

func (exp *nodeExporter) updateStorage() {
	storages, err := exp.parent.cli.NodeStorage(exp.name)
	if err != nil {
		logging.Error("get node [%s] storage: %v", exp.name, err)
		return
	}
	for _, storage := range storages {
		labels := make(prometheus.Labels)
		labels["content_vztmpl"] = "false"
		labels["content_iso"] = "false"
		labels["content_backup"] = "false"
		labels["content_snippets"] = "false"
		labels["content_rootdir"] = "false"
		labels["content_images"] = "false"
		for _, content := range storage.Content {
			switch content {
			case proxmox.ContentTemplate:
				labels["content_vztmpl"] = "true"
			case proxmox.ContentIso:
				labels["content_iso"] = "true"
			case proxmox.ContentBackup:
				labels["content_backup"] = "true"
			case proxmox.ContentSnippets:
				labels["content_snippets"] = "true"
			case proxmox.ContentRootDir:
				labels["content_rootdir"] = "true"
			case proxmox.ContentImages:
				labels["content_images"] = "true"
			}
		}
		labels["storage"] = storage.Storage
		labels["type"] = storage.Type
		exp.storageInfo.With(labels).Set(1)
	}
}
