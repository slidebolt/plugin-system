package pkg_device

import (
	"context"
	 "github.com/slidebolt/plugin-system/pkg/pkg_logic"
	"github.com/slidebolt/plugin-sdk"
	"time"
)

type SystemAdapter struct {
	bundle sdk.Bundle
	id     sdk.UUID
}

func RegisterHost(b sdk.Bundle) {
	sid := sdk.SourceID("host-system")
	
	var dev sdk.Device
	if obj, ok := b.GetBySourceID(sid); ok {
		dev = obj.(sdk.Device)
	} else {
		dev, _ = b.CreateDevice()
		dev.UpdateMetadata("Host System", sid)
	}

	// Ensure entities exist
	ensureEntity(dev, "clock", "Clock", sdk.TYPE_SENSOR)
	ensureEntity(dev, "calendar", "Calendar", sdk.TYPE_SENSOR)
	ensureEntity(dev, "metrics", "System Metrics", sdk.TYPE_SENSOR)

	a := &SystemAdapter{bundle: b, id: dev.ID()}
	go a.run(context.Background())
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
		ent.UpdateRaw(data)
	}
}
