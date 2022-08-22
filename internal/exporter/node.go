package exporter

type nodeExporter struct {
	name string
}

func newNodeExporter(name string) *nodeExporter {
	return &nodeExporter{
		name: name,
	}
}
