package main

import (
	"exporter/internal/conf"
	"exporter/internal/exporter"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/lwch/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type app struct {
	cfg *conf.Configure
}

func (app *app) Start(s service.Service) error {
	go func() {
		// cfg.Cli.SetDebug(true)
		// fmt.Println(cfg.Cli.ClusterResources(proxmox.ResourceVM))

		exp := exporter.New(app.cfg.Cli)

		reg := prometheus.NewRegistry()
		reg.MustRegister(exp)
		http.Handle("/metrics", promhttp.HandlerFor(
			reg, promhttp.HandlerOpts{Registry: reg},
		))
		http.ListenAndServe(fmt.Sprintf(":%d", app.cfg.Listen), nil)
	}()
	return nil
}

func (app *app) Stop(s service.Service) error {
	return nil
}

func main() {
	cf := flag.String("conf", "", "configure file dir")
	debug := flag.Bool("debug", false, "run in debug modes")
	act := flag.String("action", "", "install or uninstall")
	flag.Parse()

	if len(*cf) == 0 {
		fmt.Println("Missing -conf param")
		os.Exit(1)
	}

	cfg := conf.Load(*cf, *debug)

	app := app{cfg: cfg}

	cfgDir := *cf
	if !filepath.IsAbs(cfgDir) {
		var err error
		cfgDir, err = filepath.Abs(cfgDir)
		runtime.Assert(err)
	}

	svc, err := service.New(&app, &service.Config{
		Name:         "proxmox-exporter",
		DisplayName:  "proxmox-exporter",
		Description:  "proxmox prometheus exporter",
		UserName:     "root",
		Arguments:    []string{"-conf", cfgDir},
		Dependencies: []string{"After=network.target"},
	})
	runtime.Assert(err)

	switch *act {
	case "install":
		runtime.Assert(svc.Install())
	case "uninstall":
		runtime.Assert(svc.Uninstall())
	default:
		runtime.Assert(svc.Run())
	}
}
