package bundle

import (
	"context"

	"github.com/slidebolt/plugin-system/pkg/pkg_device"
	sdk "github.com/slidebolt/plugin-sdk"
)

type SystemPlugin struct {
	bundle sdk.Bundle
	cancel context.CancelFunc
	wait   func() // blocks until the poll goroutine has exited
}

func (p *SystemPlugin) Init(b sdk.Bundle) error {
	ctx, cancel := context.WithCancel(context.Background())
	p.bundle = b
	p.cancel = cancel

	b.UpdateMetadata("System Monitor")
	b.Log().Info("System Plugin Initializing...")

	p.wait = pkg_device.RegisterHost(ctx, b)
	return nil
}

// Shutdown cancels the poll goroutine and blocks until it has fully stopped.
func (p *SystemPlugin) Shutdown() {
	if p.cancel != nil {
		p.cancel()
		p.wait()
	}
}

// NewPlugin returns the concrete type so callers can reach Shutdown() without
// a type assertion.  *SystemPlugin still satisfies sdk.Plugin.
func NewPlugin() *SystemPlugin {
	return &SystemPlugin{}
}
