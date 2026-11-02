package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestURLFetchOptionsDefaultOff(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Options.BlockPrivateNetworks {
		t.Fatal("block_private_networks defaults to true")
	}
	if len(cfg.Options.AllowedSourceHosts) != 0 {
		t.Fatalf("allowed_source_hosts = %#v", cfg.Options.AllowedSourceHosts)
	}
}

func TestLoadURLFetchOptions(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	cfg, err := LoadFromContent(`{
	  "options": {
	    "block_private_networks": true,
	    "allowed_source_hosts": ["Images.Example.com", "10.0.0.8"]
	  }
	}`)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Options.BlockPrivateNetworks {
		t.Fatal("block_private_networks was not loaded")
	}
	if len(cfg.Options.AllowedSourceHosts) != 2 ||
		cfg.Options.AllowedSourceHosts[0] != "Images.Example.com" ||
		cfg.Options.AllowedSourceHosts[1] != "10.0.0.8" {
		t.Fatalf("allowed_source_hosts = %#v", cfg.Options.AllowedSourceHosts)
	}

	viper.Reset()
	cfg, err = LoadFromContent(`{"port": 3001}`)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Options == nil {
		t.Fatal("options were not defaulted")
	}
	if cfg.Options.BlockPrivateNetworks {
		t.Fatal("omitted block_private_networks loaded as true")
	}
	if len(cfg.Options.AllowedSourceHosts) != 0 {
		t.Fatalf("omitted allowed_source_hosts = %#v", cfg.Options.AllowedSourceHosts)
	}
}
