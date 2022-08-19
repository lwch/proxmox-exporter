package conf

import (
	"exporter/proxmox"
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
func Load(dir string) *Configure {
	f, err := os.Open(dir)
	runtime.Assert(err)
	defer f.Close()

	var cfg struct {
		Listen uint16 `yaml:"listen"`
		Api    struct {
			Url   string `yaml:"url"`
			User  string `yaml:"user"`
			Token string `yaml:"token"`
		} `yaml:"api"`
	}

	runtime.Assert(yaml.NewDecoder(f).Decode(&cfg))

	return &Configure{
		Listen: cfg.Listen,
		Cli: proxmox.New(cfg.Api.Url,
			cfg.Api.User, cfg.Api.Token, 5*time.Second),
	}
}
