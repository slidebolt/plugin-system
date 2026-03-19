package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	contract "github.com/slidebolt/sb-contract"
	domain "github.com/slidebolt/sb-domain"
	messenger "github.com/slidebolt/sb-messenger-sdk"
	storage "github.com/slidebolt/sb-storage-sdk"
)

const PluginID = "plugin-system"

type Time struct {
	Timestamp int64 `json:"timestamp"`
	Hour      int   `json:"hour"`
	Minute    int   `json:"minute"`
	Day       int   `json:"day"`
	Week      int   `json:"week"`
	Month     int   `json:"month"`
	Year      int   `json:"year"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude,omitempty"`
	City      string  `json:"city,omitempty"`
	Country   string  `json:"country,omitempty"`
	Timezone  string  `json:"timezone,omitempty"`
}

func init() {
	domain.Register("time", Time{})
	domain.Register("location", Location{})
}

type App struct {
	msg        messenger.Messenger
	store      storage.Storage
	cmds       *messenger.Commands
	subs       []messenger.Subscription
	timeTicker *time.Ticker
	stopCh     chan struct{}
}

func New() *App { return &App{} }

func (a *App) Hello() contract.HelloResponse {
	return contract.HelloResponse{
		ID:              PluginID,
		Kind:            contract.KindPlugin,
		ContractVersion: contract.ContractVersion,
		DependsOn:       []string{"messenger", "storage"},
	}
}

func (a *App) OnStart(deps map[string]json.RawMessage) (json.RawMessage, error) {
	msg, err := messenger.Connect(deps)
	if err != nil {
		return nil, fmt.Errorf("connect messenger: %w", err)
	}
	a.msg = msg

	storeClient, err := storage.Connect(deps)
	if err != nil {
		return nil, fmt.Errorf("connect storage: %w", err)
	}
	a.store = storeClient

	a.cmds = messenger.NewCommands(msg, domain.LookupCommand)
	sub, err := a.cmds.Receive(PluginID+".>", a.handleCommand)
	if err != nil {
		return nil, fmt.Errorf("subscribe commands: %w", err)
	}
	a.subs = append(a.subs, sub)

	if err := a.SeedSystemDevices(); err != nil {
		return nil, fmt.Errorf("seed system devices: %w", err)
	}

	a.stopCh = make(chan struct{})
	a.startTimeUpdates()

	log.Println("plugin-system: started")
	return nil, nil
}

func (a *App) OnShutdown() error {
	if a.stopCh != nil {
		close(a.stopCh)
	}
	if a.timeTicker != nil {
		a.timeTicker.Stop()
	}
	for _, sub := range a.subs {
		sub.Unsubscribe()
	}
	if a.store != nil {
		a.store.Close()
	}
	if a.msg != nil {
		a.msg.Close()
	}
	return nil
}

func (a *App) handleCommand(addr messenger.Address, cmd any) {
	switch c := cmd.(type) {
	case domain.LightTurnOn:
		log.Printf("plugin-system: light %s turn_on", addr.Key())
	case domain.LightTurnOff:
		log.Printf("plugin-system: light %s turn_off transition=%v", addr.Key(), c.Transition)
	case domain.LightSetBrightness:
		log.Printf("plugin-system: light %s set_brightness brightness=%d", addr.Key(), c.Brightness)
	case domain.LightSetColorTemp:
		log.Printf("plugin-system: light %s set_color_temp mireds=%d", addr.Key(), c.Mireds)
	case domain.LightSetRGB:
		log.Printf("plugin-system: light %s set_rgb r=%d g=%d b=%d", addr.Key(), c.R, c.G, c.B)
	case domain.LightSetRGBW:
		log.Printf("plugin-system: light %s set_rgbw r=%d g=%d b=%d w=%d", addr.Key(), c.R, c.G, c.B, c.W)
	case domain.LightSetRGBWW:
		log.Printf("plugin-system: light %s set_rgbww r=%d g=%d b=%d cw=%d ww=%d", addr.Key(), c.R, c.G, c.B, c.CW, c.WW)
	case domain.LightSetHS:
		log.Printf("plugin-system: light %s set_hs hue=%.1f sat=%.1f", addr.Key(), c.Hue, c.Saturation)
	case domain.LightSetXY:
		log.Printf("plugin-system: light %s set_xy x=%.4f y=%.4f", addr.Key(), c.X, c.Y)
	case domain.LightSetWhite:
		log.Printf("plugin-system: light %s set_white white=%d", addr.Key(), c.White)
	case domain.LightSetEffect:
		log.Printf("plugin-system: light %s set_effect effect=%s", addr.Key(), c.Effect)
	case domain.SwitchTurnOn:
		log.Printf("plugin-system: switch %s turn_on", addr.Key())
	case domain.SwitchTurnOff:
		log.Printf("plugin-system: switch %s turn_off", addr.Key())
	case domain.SwitchToggle:
		log.Printf("plugin-system: switch %s toggle", addr.Key())
	case domain.FanTurnOn:
		log.Printf("plugin-system: fan %s turn_on", addr.Key())
	case domain.FanTurnOff:
		log.Printf("plugin-system: fan %s turn_off", addr.Key())
	case domain.FanSetSpeed:
		log.Printf("plugin-system: fan %s set_speed percentage=%d", addr.Key(), c.Percentage)
	case domain.CoverOpen:
		log.Printf("plugin-system: cover %s open", addr.Key())
	case domain.CoverClose:
		log.Printf("plugin-system: cover %s close", addr.Key())
	case domain.CoverSetPosition:
		log.Printf("plugin-system: cover %s set_position pos=%d", addr.Key(), c.Position)
	case domain.LockLock:
		log.Printf("plugin-system: lock %s lock", addr.Key())
	case domain.LockUnlock:
		log.Printf("plugin-system: lock %s unlock", addr.Key())
	case domain.ButtonPress:
		log.Printf("plugin-system: button %s press", addr.Key())
	case domain.NumberSetValue:
		log.Printf("plugin-system: number %s set_value value=%v", addr.Key(), c.Value)
	case domain.SelectOption:
		log.Printf("plugin-system: select %s set_option option=%s", addr.Key(), c.Option)
	case domain.TextSetValue:
		log.Printf("plugin-system: text %s set_value value=%s", addr.Key(), c.Value)
	case domain.ClimateSetMode:
		log.Printf("plugin-system: climate %s set_mode mode=%s", addr.Key(), c.HVACMode)
	case domain.ClimateSetTemperature:
		log.Printf("plugin-system: climate %s set_temperature temp=%v", addr.Key(), c.Temperature)
	default:
		log.Printf("plugin-system: unknown command %T for %s", cmd, addr.Key())
	}
}

func IsSystemDevice(deviceID string) bool {
	return deviceID == "time" || deviceID == "location"
}

func (a *App) entityExists(plugin, device, id string) bool {
	key := domain.EntityKey{Plugin: plugin, DeviceID: device, ID: id}
	_, err := a.store.Get(key)
	return err == nil
}

func (a *App) SeedSystemDevices() error {
	now := time.Now()
	_, week := now.ISOWeek()

	timeEntities := []domain.Entity{
		{ID: "timestamp", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Timestamp", State: Time{Timestamp: now.Unix()}},
		{ID: "hour", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Hour", State: Time{Hour: now.Hour()}},
		{ID: "minute", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Minute", State: Time{Minute: now.Minute()}},
		{ID: "day", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Day", State: Time{Day: now.Day()}},
		{ID: "week", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Week", State: Time{Week: week}},
		{ID: "month", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Month", State: Time{Month: int(now.Month())}},
		{ID: "year", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Year", State: Time{Year: now.Year()}},
	}
	for _, e := range timeEntities {
		if !a.entityExists(e.Plugin, e.DeviceID, e.ID) {
			if err := a.store.Save(e); err != nil {
				return fmt.Errorf("save time entity %s: %w", e.ID, err)
			}
		}
	}

	if !a.entityExists(PluginID, "location", "position") {
		if err := a.store.Save(domain.Entity{
			ID: "position", Plugin: PluginID, DeviceID: "location", Type: "location", Name: "Position",
			State: Location{Latitude: 0.0, Longitude: 0.0, Timezone: "UTC"},
		}); err != nil {
			return fmt.Errorf("save location entity: %w", err)
		}
	}

	return nil
}

func (a *App) startTimeUpdates() {
	a.timeTicker = time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-a.timeTicker.C:
				a.updateTimeEntities()
			case <-a.stopCh:
				return
			}
		}
	}()
}

func (a *App) updateTimeEntities() {
	now := time.Now()
	_, week := now.ISOWeek()

	_ = a.store.Save(domain.Entity{ID: "timestamp", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Timestamp", State: Time{Timestamp: now.Unix()}})

	if now.Second() == 0 {
		_ = a.store.Save(domain.Entity{ID: "hour", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Hour", State: Time{Hour: now.Hour()}})
		_ = a.store.Save(domain.Entity{ID: "minute", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Minute", State: Time{Minute: now.Minute()}})
	}

	if now.Hour() == 0 && now.Minute() == 0 && now.Second() == 0 {
		_ = a.store.Save(domain.Entity{ID: "day", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Day", State: Time{Day: now.Day()}})
		_ = a.store.Save(domain.Entity{ID: "week", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Week", State: Time{Week: week}})
		_ = a.store.Save(domain.Entity{ID: "month", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Month", State: Time{Month: int(now.Month())}})
		_ = a.store.Save(domain.Entity{ID: "year", Plugin: PluginID, DeviceID: "time", Type: "time", Name: "Year", State: Time{Year: now.Year()}})
	}
}
