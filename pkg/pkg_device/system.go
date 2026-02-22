package pkg_device

import (
	"context"
	"sync"
	"time"

	"github.com/slidebolt/plugin-system/pkg/pkg_logic"
	sdk "github.com/slidebolt/plugin-sdk"
)

type SystemAdapter struct {
	bundle sdk.Bundle
	id     sdk.UUID
}

// registerMu guards the check-then-create block in RegisterHost so concurrent
// callers (e.g. parallel tests hitting the same bundle) cannot each see "not
// found" and double-create the device or its entities.
var registerMu sync.Mutex

// RegisterHost creates (or reuses) the "host-system" device, attaches the
// three sensor entities, and starts the 1-second poll loop.  It returns a
// blocking wait function that resolves only after the goroutine has exited —
// call it after cancelling ctx to confirm a clean stop.
func RegisterHost(ctx context.Context, b sdk.Bundle) func() {
	registerMu.Lock()
	sid := sdk.SourceID("host-system")

	var dev sdk.Device
	if obj, ok := b.GetBySourceID(sid); ok {
		dev = obj.(sdk.Device)
	} else {
		dev, _ = b.CreateDevice()
		dev.UpdateMetadata("Host System", sid)
	}

	ensureEntity(dev, "clock", "Clock", sdk.TYPE_SENSOR)
	ensureEntity(dev, "calendar", "Calendar", sdk.TYPE_SENSOR)
	ensureEntity(dev, "metrics", "System Metrics", sdk.TYPE_SENSOR)
	registerMu.Unlock()

	a := &SystemAdapter{bundle: b, id: dev.ID()}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		a.run(ctx)
	}()

	return wg.Wait
}

func ensureEntity(dev sdk.Device, sid string, name string, kind sdk.EntityType) sdk.Entity {
	if obj, ok := dev.GetBySourceID(sdk.SourceID(sid)); ok {
		return obj.(sdk.Entity)
	}
	ent, _ := dev.CreateEntity(kind)
	ent.UpdateMetadata(name, sdk.SourceID(sid))
	return ent
}

func (a *SystemAdapter) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	dev, _ := a.bundle.GetDevice(a.id)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateState(dev, "clock", pkg_logic.GetClockState())
			updateState(dev, "calendar", pkg_logic.GetCalendarState())
			updateState(dev, "metrics", pkg_logic.GetMetricsState())
			a.bundle.Publish("system.tick", map[string]interface{}{
				"timestamp": time.Now().Unix(),
			})
		}
	}
}

func updateState(dev sdk.Device, sid string, data map[string]interface{}) {
	if obj, ok := dev.GetBySourceID(sdk.SourceID(sid)); ok {
		ent := obj.(sdk.Entity)
		// System metrics are purely functional, move to Properties.
		_ = ent.UpdateProperties(data)
	}
}
