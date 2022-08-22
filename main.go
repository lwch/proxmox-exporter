package main

import (
	"exporter/internal/conf"
	"exporter/internal/exporter"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cf := flag.String("conf", "", "configure file dir")
	flag.Parse()

	if len(*cf) == 0 {
		fmt.Println("Missing -conf param")
		os.Exit(1)
	}

	cfg := conf.Load(*cf)

	// cfg.Cli.SetDebug(true)
	// fmt.Println(cfg.Cli.ClusterTasks())

	exp := exporter.New(cfg.Cli)

	reg := prometheus.NewRegistry()
	reg.MustRegister(exp)
	http.Handle("/metrics", promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{Registry: reg},
	))
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Listen), nil)
}
