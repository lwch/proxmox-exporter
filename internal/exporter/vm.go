package exporter

import "github.com/prometheus/client_golang/prometheus"

type vmExporter struct {
	parent *nodeExporter
}

func newVmExporter(parent *nodeExporter) *vmExporter {
	return &vmExporter{parent: parent}
}

func (exp *vmExporter) Describe(ch chan<- *prometheus.Desc) {
}

func (exp *vmExporter) Collect(ch chan<- prometheus.Metric) {
}
