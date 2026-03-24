package app_test

import (
	"encoding/json"
	"testing"

	systemapp "github.com/slidebolt/plugin-system/app"
	domain "github.com/slidebolt/sb-domain"
	testkit "github.com/slidebolt/sb-testkit"
)

func TestStorageContract_OnStartSeedsSystemEntities(t *testing.T) {
	env := testkit.NewTestEnv(t)
	env.Start("messenger")
	env.Start("storage")

	app := systemapp.New()
	if _, err := app.OnStart(map[string]json.RawMessage{"messenger": env.MessengerPayload()}); err != nil {
		t.Fatalf("OnStart: %v", err)
	}
	t.Cleanup(func() { _ = app.OnShutdown() })

	raw, err := env.Storage().Get(domain.EntityKey{Plugin: systemapp.PluginID, DeviceID: "time", ID: "timestamp"})
	if err != nil {
		t.Fatalf("get timestamp entity: %v", err)
	}
	var got domain.Entity
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	state, ok := got.State.(systemapp.Time)
	if !ok {
		t.Fatalf("state type = %T", got.State)
	}
	if state.Timestamp == 0 {
		t.Fatalf("state = %+v", state)
	}
}
