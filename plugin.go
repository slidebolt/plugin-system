package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	runner "github.com/slidebolt/sdk-runner"
	"github.com/slidebolt/sdk-types"
)

const (
	systemDeviceID = "system-device"
	timeEntityID   = "system-time"
	dateEntityID   = "system-date"
	cpuEntityID    = "system-cpu"
)

type SystemPlugin struct {
	eventSink runner.EventSink

	cpuMu     sync.Mutex
	prevTotal uint64
	prevIdle  uint64
	cpuInit   bool
}

func NewSystemPlugin() *SystemPlugin {
	return &SystemPlugin{}
}

func (p *SystemPlugin) OnInitialize(config runner.Config, state types.Storage) (types.Manifest, types.Storage) {
	p.eventSink = config.EventSink
	return types.Manifest{ID: "plugin-system", Name: "System Plugin", Version: "1.0.0"}, state
}

func (p *SystemPlugin) OnReady() {
	go p.tickLoop()
}

func (p *SystemPlugin) OnShutdown() {}

func (p *SystemPlugin) OnHealthCheck() (string, error) { return "perfect", nil }

func (p *SystemPlugin) OnStorageUpdate(current types.Storage) (types.Storage, error) {
	return current, nil
}

func (p *SystemPlugin) OnDeviceCreate(dev types.Device) (types.Device, error) { return dev, nil }
func (p *SystemPlugin) OnDeviceUpdate(dev types.Device) (types.Device, error) { return dev, nil }
func (p *SystemPlugin) OnDeviceDelete(id string) error                        { return nil }

func (p *SystemPlugin) OnDevicesList(current []types.Device) ([]types.Device, error) {
	byID := map[string]types.Device{}
	for _, d := range current {
		byID[d.ID] = d
	}
	byID[systemDeviceID] = runner.ReconcileDevice(byID[systemDeviceID], types.Device{
		ID:         systemDeviceID,
		SourceID:   systemDeviceID,
		SourceName: "System",
	})
	out := make([]types.Device, 0, len(byID))
	for _, d := range byID {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (p *SystemPlugin) OnDeviceSearch(q types.SearchQuery, results []types.Device) ([]types.Device, error) {
	return results, nil
}

func (p *SystemPlugin) OnEntityCreate(ent types.Entity) (types.Entity, error) { return ent, nil }
func (p *SystemPlugin) OnEntityUpdate(ent types.Entity) (types.Entity, error) { return ent, nil }
func (p *SystemPlugin) OnEntityDelete(deviceID, entityID string) error        { return nil }

func (p *SystemPlugin) OnEntitiesList(deviceID string, current []types.Entity) ([]types.Entity, error) {
	if deviceID != systemDeviceID {
		return current, nil
	}
	need := []types.Entity{
		{ID: timeEntityID, DeviceID: systemDeviceID, Domain: "sensor.time", LocalName: "Time"},
		{ID: dateEntityID, DeviceID: systemDeviceID, Domain: "sensor.date", LocalName: "Date"},
		{ID: cpuEntityID, DeviceID: systemDeviceID, Domain: "sensor.cpu", LocalName: "CPU"},
	}
	for _, ent := range need {
		if !entityExists(current, ent.ID) {
			current = append(current, ent)
		}
	}
	return current, nil
}

func (p *SystemPlugin) OnCommand(cmd types.Command, entity types.Entity) (types.Entity, error) {
	// System sensors are read-only and do not support writable commands by design.
	return entity, fmt.Errorf("commands are not supported for system sensors")
}

func (p *SystemPlugin) OnEvent(evt types.Event, entity types.Entity) (types.Entity, error) {
	entity.Data.Reported = evt.Payload
	entity.Data.Effective = evt.Payload
	entity.Data.SyncStatus = "in_sync"
	return entity, nil
}

func (p *SystemPlugin) tickLoop() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for t := range ticker.C {
		p.emitSystemEvents(t)
	}
}

func (p *SystemPlugin) emitSystemEvents(now time.Time) {
	if p.eventSink == nil {
		return
	}
	timePayload, _ := json.Marshal(map[string]any{
		"type":  "tick",
		"value": now.Format("15:04:05.000"),
		"ts":    now.UTC().Format(time.RFC3339Nano),
	})
	_ = p.eventSink.EmitEvent(types.InboundEvent{DeviceID: systemDeviceID, EntityID: timeEntityID, Payload: timePayload})

	datePayload, _ := json.Marshal(map[string]any{
		"type":  "tick",
		"value": now.Format("2006-01-02 15:04:05.000"),
		"ts":    now.UTC().Format(time.RFC3339Nano),
	})
	_ = p.eventSink.EmitEvent(types.InboundEvent{DeviceID: systemDeviceID, EntityID: dateEntityID, Payload: datePayload})

	cpuPayload, _ := json.Marshal(map[string]any{
		"type":    "tick",
		"percent": p.readCPUPercent(),
		"ts":      now.UTC().Format(time.RFC3339Nano),
	})
	_ = p.eventSink.EmitEvent(types.InboundEvent{DeviceID: systemDeviceID, EntityID: cpuEntityID, Payload: cpuPayload})
}

func entityExists(current []types.Entity, id string) bool {
	for _, e := range current {
		if e.ID == id {
			return true
		}
	}
	return false
}

func (p *SystemPlugin) readCPUPercent() float64 {
	total, idle, err := readProcStat()
	if err != nil {
		return 0
	}
	p.cpuMu.Lock()
	defer p.cpuMu.Unlock()
	if !p.cpuInit {
		p.prevTotal, p.prevIdle, p.cpuInit = total, idle, true
		return 0
	}
	dTotal := total - p.prevTotal
	dIdle := idle - p.prevIdle
	p.prevTotal, p.prevIdle = total, idle
	if dTotal == 0 {
		return 0
	}
	busy := float64(dTotal-dIdle) / float64(dTotal) * 100
	if busy < 0 {
		busy = 0
	}
	if busy > 100 {
		busy = 100
	}
	return busy
}

func readProcStat() (total uint64, idle uint64, err error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	if !s.Scan() {
		return 0, 0, fmt.Errorf("/proc/stat empty")
	}
	fields := strings.Fields(s.Text())
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, fmt.Errorf("unexpected /proc/stat format")
	}
	vals := make([]uint64, 0, len(fields)-1)
	for _, field := range fields[1:] {
		var v uint64
		_, scanErr := fmt.Sscanf(field, "%d", &v)
		if scanErr != nil {
			return 0, 0, scanErr
		}
		vals = append(vals, v)
		total += v
	}
	idle = vals[3]
	if len(vals) > 4 {
		idle += vals[4]
	}
	return total, idle, nil
}