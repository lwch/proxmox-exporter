package main

import (
	"exporter/internal/conf"
	"exporter/proxmox"
	"flag"
	"fmt"
	"os"

	"github.com/lwch/runtime"
)

func main() {
	cf := flag.String("conf", "", "configure file dir")
	flag.Parse()

	if len(*cf) == 0 {
		fmt.Println("Missing -conf param")
		os.Exit(1)
	}

	cfg := conf.Load(*cf)

	list, err := cfg.Cli.Resources(proxmox.ResourceVM)
	runtime.Assert(err)
	fmt.Println(list)
}
