package exporter

import (
	"bufio"
	"exporter/proxmox"
	"exporter/sensors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/anatol/smart.go"
	"github.com/jaypipes/ghw"
	"github.com/lwch/logging"
	"github.com/prometheus/client_golang/prometheus"
)

type nodeExporter struct {
	parent *Exporter
	name   string
	vm     *vmExporter

	// stats
	uptime            prometheus.Gauge
	info              *prometheus.GaugeVec
	cpuUsage          prometheus.Gauge
	cpuLoadAverage    *prometheus.GaugeVec
	cpuFrequency      *prometheus.GaugeVec
	memoryUsed        prometheus.Gauge
	memoryFree        prometheus.Gauge
	memoryTotal       prometheus.Gauge
	swapUsed          prometheus.Gauge
	swapFree          prometheus.Gauge
	swapTotal         prometheus.Gauge
	rootfsUsed        prometheus.Gauge
	rootfsFree        prometheus.Gauge
	rootfsTotal       prometheus.Gauge
	ioWait            prometheus.Gauge
	storageInfo       *prometheus.GaugeVec
	storageUsed       *prometheus.GaugeVec
	storageFree       *prometheus.GaugeVec
	storageTotal      *prometheus.GaugeVec
	storageUsage      *prometheus.GaugeVec
	sensors           *prometheus.GaugeVec
	netin             prometheus.Gauge
	netout            prometheus.Gauge
	smartTemperature  *prometheus.GaugeVec
	smartUsedPercent  *prometheus.GaugeVec
	smartReaden       *prometheus.GaugeVec
	smartWritten      *prometheus.GaugeVec
	smartPowerCycles  *prometheus.GaugeVec
	smartPowerOnHours *prometheus.GaugeVec
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
	exp.uptime = prometheus.NewGauge(prometheus.GaugeOpts{
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
	exp.cpuFrequency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "cpu_frequency",
		Help:        "node cpu frequency of each core",
		ConstLabels: labels,
	}, []string{"processor"})

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
		"storage",
		"type",
		"content_vztmpl",
		"content_iso",
		"content_backup",
		"content_snippets",
		"content_rootdir",
		"content_images",
	})
	exp.storageUsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "storage_used",
		Help:        "node storage used bytes",
		ConstLabels: labels,
	}, []string{"storage_name"})
	exp.storageFree = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "storage_free",
		Help:        "node storage free bytes",
		ConstLabels: labels,
	}, []string{"storage_name"})
	exp.storageTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "storage_total",
		Help:        "node storage total bytes",
		ConstLabels: labels,
	}, []string{"storage_name"})
	exp.storageUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "storage_usage",
		Help:        "node storage usage ratio(precent)",
		ConstLabels: labels,
	}, []string{"storage_name"})
	// sensors
	exp.sensors = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "sensors",
		Help:        "use sensors command to get device temperature and cpu fan speed",
		ConstLabels: labels,
	}, []string{"chip_name", "label_name", "feature_name"})
	// network
	exp.netin = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "netin",
		Help:        "node received bytes",
		ConstLabels: labels,
	})
	exp.netout = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "netout",
		Help:        "node sent bytes",
		ConstLabels: labels,
	})
	// smart
	smartLabels := []string{"device", "type"}
	exp.smartTemperature = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "smart_temperature",
		Help:        "temperature of smart data",
		ConstLabels: labels,
	}, smartLabels)
	exp.smartWritten = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "smart_written",
		Help:        "written bytes of smart data(lba 512 bytes padding)",
		ConstLabels: labels,
	}, smartLabels)
	exp.smartReaden = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "smart_readden",
		Help:        "readden bytes of smart data(lba 512 bytes padding)",
		ConstLabels: labels,
	}, smartLabels)
	exp.smartUsedPercent = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "smart_used_percent",
		Help:        "used percent of smart data(nvme)",
		ConstLabels: labels,
	}, smartLabels)
	exp.smartPowerOnHours = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "smart_poweron_hours",
		Help:        "poweron hours of smart data",
		ConstLabels: labels,
	}, smartLabels)
	exp.smartPowerCycles = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "smart_power_cycles",
		Help:        "power cycles of smart data",
		ConstLabels: labels,
	}, smartLabels)
}

func (exp *nodeExporter) Describe(ch chan<- *prometheus.Desc) {
	// online
	exp.uptime.Describe(ch)
	exp.info.Describe(ch)
	// cpu
	exp.cpuUsage.Describe(ch)
	exp.cpuLoadAverage.Describe(ch)
	exp.cpuFrequency.Describe(ch)
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
	exp.storageUsed.Describe(ch)
	exp.storageFree.Describe(ch)
	exp.storageTotal.Describe(ch)
	exp.storageUsage.Describe(ch)
	// sensors
	exp.sensors.Describe(ch)
	// network
	exp.netin.Describe(ch)
	exp.netout.Describe(ch)
	// smart
	exp.smartTemperature.Describe(ch)
	exp.smartUsedPercent.Describe(ch)
	exp.smartReaden.Describe(ch)
	exp.smartWritten.Describe(ch)
	exp.smartPowerOnHours.Describe(ch)
	exp.smartPowerCycles.Describe(ch)

	// vm describe
	exp.vm.Describe(ch)
}

func (exp *nodeExporter) Collect(ch chan<- prometheus.Metric) {
	// collect values
	exp.updateStatus()

	// online
	exp.uptime.Collect(ch)
	exp.info.Collect(ch)
	// cpu
	exp.cpuUsage.Collect(ch)
	exp.cpuLoadAverage.Collect(ch)
	exp.cpuFrequency.Collect(ch)
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
	exp.storageUsed.Collect(ch)
	exp.storageFree.Collect(ch)
	exp.storageTotal.Collect(ch)
	exp.storageUsage.Collect(ch)
	// sensors
	exp.sensors.Collect(ch)
	// network
	exp.netin.Collect(ch)
	exp.netout.Collect(ch)
	// smart
	exp.smartTemperature.Collect(ch)
	exp.smartUsedPercent.Collect(ch)
	exp.smartReaden.Collect(ch)
	exp.smartWritten.Collect(ch)
	exp.smartPowerOnHours.Collect(ch)
	exp.smartPowerCycles.Collect(ch)

	// vm collect
	exp.vm.Collect(ch)
}

func (exp *nodeExporter) updateStatus() {
	status, err := exp.parent.cli.NodeStatus(exp.name)
	if err != nil {
		logging.Error("get node status: %v", err)
		return
	}
	exp.updateInfo(status)
	exp.updateCpu(status)
	exp.updateMemory(status)
	exp.updateDisk(status)
	exp.updateSensors()
	exp.updateNetwork()
}

func (exp *nodeExporter) updateInfo(status proxmox.NodeStatus) {
	exp.uptime.Set(float64(status.Uptime))
	exp.info.With(prometheus.Labels{
		"model":          status.CpuInfo.Model,
		"sockets":        fmt.Sprintf("%d", status.CpuInfo.Sockets),
		"cores":          fmt.Sprintf("%d", status.CpuInfo.Cores),
		"threads":        fmt.Sprintf("%d", status.CpuInfo.Threads),
		"mhz":            status.CpuInfo.Frequency,
		"kernel_version": status.KernelVersion,
		"pve_version":    status.PveVersion,
	}).Inc()
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

	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		logging.Error("get cpu info: %v", err)
		return
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var processor string
	var mhz float64
	for s.Scan() {
		fields := strings.Split(s.Text(), ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "processor":
			if len(processor) > 0 {
				label := prometheus.Labels{"processor": processor}
				exp.cpuFrequency.With(label).Set(mhz)
			}
			processor = value
		case "cpu MHz", "clock":
			if t, err := strconv.ParseFloat(strings.Replace(value, "MHz", "", 1), 64); err == nil {
				mhz = t
			}
		}
	}
	if len(processor) > 0 {
		label := prometheus.Labels{"processor": processor}
		exp.cpuFrequency.With(label).Set(mhz)
	}
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
	exp.updateSmart()
}

func (exp *nodeExporter) updateStorage() {
	storages, err := exp.parent.cli.NodeStorage(exp.name)
	if err != nil {
		logging.Error("get node storage: %v", err)
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
		exp.storageInfo.With(labels).Inc()

		labels = prometheus.Labels{"storage_name": storage.Storage}
		exp.storageUsed.With(labels).Set(float64(storage.Used))
		exp.storageFree.With(labels).Set(float64(storage.Available))
		exp.storageTotal.With(labels).Set(float64(storage.Total))
		exp.storageUsage.With(labels).Set(storage.Ratio * 100.)
	}
}

func (exp *nodeExporter) updateSensors() {
	sensors, err := sensors.Get()
	if err != nil {
		logging.Error("get sensors: %v", err)
		return
	}
	for _, sensor := range sensors {
		labels := prometheus.Labels{"chip_name": sensor.Chip}
		for label, value := range sensor.Features {
			labels["label_name"] = label
			for feature, value := range value {
				labels["feature_name"] = feature
				exp.sensors.With(labels).Set(value)
			}
		}
	}
}

func (exp *nodeExporter) updateNetwork() {
	datas, err := exp.parent.cli.NodeRrdData(exp.name)
	if err != nil {
		logging.Error("get node rrddata: %v", err)
		return
	}
	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Time > datas[j].Time
	})
	if len(datas) > 0 {
		exp.netin.Set(datas[0].NetIn)
		exp.netout.Set(datas[0].NetOut)
	}
}

func (exp *nodeExporter) updateSmart() {
	block, err := ghw.Block()
	if err != nil {
		logging.Error("get disks: %v", err)
		return
	}
	for _, disk := range block.Disks {
		if disk.DriveType != ghw.DRIVE_TYPE_HDD &&
			disk.DriveType != ghw.DRIVE_TYPE_SSD {
			continue
		}

		dev, err := smart.Open("/dev/" + disk.Name)
		if err != nil {
			logging.Error("open smart [%s]: %v", disk.Name, err)
			continue
		}
		defer dev.Close()

		labels := prometheus.Labels{"device": disk.Name}

		switch sm := dev.(type) {
		case *smart.NVMeDevice:
			log, err := sm.ReadSMART()
			if err != nil {
				logging.Error("read smart [%s]: %v", disk.Name, err)
				continue
			}
			labels["type"] = "nvme"
			exp.smartTemperature.With(labels).Set(float64(log.Temperature) - 273.1)
			exp.smartUsedPercent.With(labels).Set(float64(log.PercentUsed))
			exp.smartReaden.With(labels).Set(float64(log.DataUnitsRead.Val[0]))
			exp.smartWritten.With(labels).Set(float64(log.DataUnitsWritten.Val[0]))
			exp.smartPowerOnHours.With(labels).Set(float64(log.PowerOnHours.Val[0]))
			exp.smartPowerCycles.With(labels).Set(float64(log.PowerCycles.Val[0]))
		case *smart.SataDevice:
			page, err := sm.ReadSMARTData()
			if err != nil {
				logging.Error("read smart [%s]: %v", disk.Name, err)
				continue
			}
			labels["type"] = "sata"
			for _, attr := range page.Attrs {
				switch attr.Name {
				case "Airflow_Temperature_Cel":
					current, _, _, _, err := attr.ParseAsTemperature()
					if err != nil {
						logging.Error("get temperature [%s]: %v", disk.Name, err)
						continue
					}
					exp.smartTemperature.With(labels).Set(float64(current))
				case "Total_LBAs_Read":
					exp.smartReaden.With(labels).Set(float64(attr.ValueRaw))
				case "Total_LBAs_Written":
					exp.smartWritten.With(labels).Set(float64(attr.ValueRaw))
				case "Power_On_Hours":
					exp.smartPowerOnHours.With(labels).Set(float64(attr.ValueRaw))
				case "Power_Cycle_Count":
					exp.smartPowerCycles.With(labels).Set(float64(attr.ValueRaw))
				}
			}
		}
	}
}
