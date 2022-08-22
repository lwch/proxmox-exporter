package conf

import (
	"exporter/proxmox"
	"fmt"
	"os"
	"time"

	"github.com/lwch/runtime"
	"gopkg.in/yaml.v3"
)

// Configure configure struct
type Configure struct {
	Listen uint16
	Cli    *proxmox.Client
}

// Load load configure file
func Load(dir string, debug bool) *Configure {
	f, err := os.Open(dir)
	runtime.Assert(err)
	defer f.Close()

	var cfg struct {
		Listen uint16 `yaml:"listen"`
		Api    struct {
			User  string `yaml:"user"`
			Token string `yaml:"token"`
		} `yaml:"api"`
	}

	runtime.Assert(yaml.NewDecoder(f).Decode(&cfg))

	ip := "127.0.0.1"
	if debug {
		var ok bool
		ip, ok = os.LookupEnv("PROXMOX_IP")
		if !ok {
			ip = "127.0.0.1"
		}
	}

	return &Configure{
		Listen: cfg.Listen,
		Cli: proxmox.New(fmt.Sprintf("https://%s:8006", ip),
			cfg.Api.User, cfg.Api.Token, 5*time.Second),
	}
}
