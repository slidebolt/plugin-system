package bundle

import (
	 "github.com/slidebolt/plugin-system/pkg/pkg_device"
	"github.com/slidebolt/plugin-sdk"
)

type SystemPlugin struct {
	bundle sdk.Bundle
}

func (p *SystemPlugin) Init(b sdk.Bundle) error {
	p.bundle = b
	b.UpdateMetadata("System Monitor")
	b.Log().Info("System Plugin Initializing...")

	// Register the local host metrics device
	pkg_device.RegisterHost(b)

	return nil
}

func NewPlugin() sdk.Plugin {
	return &SystemPlugin{}
}
