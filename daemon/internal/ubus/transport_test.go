package ubus

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/valercha/OpenHapp/daemon/internal/config"
	"github.com/valercha/OpenHapp/daemon/internal/engine"
	"github.com/valercha/OpenHapp/daemon/internal/manifest"
	"github.com/valercha/OpenHapp/daemon/internal/profile"
	"github.com/valercha/OpenHapp/daemon/internal/service"
	"github.com/valercha/OpenHapp/daemon/internal/state"
)

func newTestDispatcher(t *testing.T) *Dispatcher {
	t.Helper()

	cfg := config.Default()
	st := state.New("test")
	svc := service.New(cfg, st)
	profiles := profile.NewStore(t.TempDir() + "/openhapp")

	return NewDispatcher(New(
		svc,
		st,
		cfg,
		manifest.FromConfig("test", cfg),
		profiles,
	))
}

func TestHandleJSON(t *testing.T) {
	d := newTestDispatcher(t)

	payload := []byte(`{"method":"status"}`)
	response, err := d.HandleJSON(context.Background(), payload)
	if err != nil {
		t.Fatalf("handle json: %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(response, &decoded); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if decoded.Error != "" {
		t.Fatalf("unexpected error: %s", decoded.Error)
	}
}

func TestHandleJSONRejectsMalformedRequest(t *testing.T) {
	d := newTestDispatcher(t)

	if _, err := d.HandleJSON(
		context.Background(),
		[]byte(`{"method":`),
	); err == nil {
		t.Fatal("expected malformed JSON error")
	}
}

func TestHandleJSONEngineInfo(t *testing.T) {
	cfg := config.Default()
	cfg.Engine = "sing-box"

	st := state.New("test")
	svc := service.New(cfg, st)

	eng := engine.New("sing-box")
	svc.SetEngine(eng)

	profiles := profile.NewStore(t.TempDir() + "/openhapp")

	d := NewDispatcher(New(
		svc,
		st,
		cfg,
		manifest.FromConfig("test", cfg),
		profiles,
	))

	payload := []byte(`{"method":"engine_info"}`)
	response, err := d.HandleJSON(context.Background(), payload)
	if err != nil {
		t.Fatalf("handle engine_info: %v", err)
	}

	var decoded Response
	if err := json.Unmarshal(response, &decoded); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if decoded.Error != "" {
		t.Fatalf("unexpected error: %s", decoded.Error)
	}

	raw, err := json.Marshal(decoded.Result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}

	var info map[string]any
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("decode engine info: %v", err)
	}

	if info["name"] != "sing-box" {
		t.Fatalf("unexpected engine name: %#v", info["name"])
	}

	if info["running"] != false {
		t.Fatalf("unexpected running state: %#v", info["running"])
	}
}

func TestHandleJSONProfileCRUD(t *testing.T) {
	d := newTestDispatcher(t)
	ctx := context.Background()

	addPayload := []byte(`{
		"method":"profile_add",
		"params":{
			"id":"de-01",
			"name":"Germany 01",
			"type":"vless",
			"server":"example.com",
			"port":443,
			"enabled":true,
			"properties":{
				"uuid":"test-uuid",
				"sni":"example.com"
			}
		}
	}`)

	response, err := d.HandleJSON(ctx, addPayload)
	if err != nil {
		t.Fatalf("profile_add: %v", err)
	}

	var addResponse Response
	if err := json.Unmarshal(response, &addResponse); err != nil {
		t.Fatalf("decode profile_add response: %v", err)
	}

	if addResponse.Error != "" {
		t.Fatalf("profile_add error: %s", addResponse.Error)
	}

	listResponse, err := d.HandleJSON(
		ctx,
		[]byte(`{"method":"profile_list"}`),
	)
	if err != nil {
		t.Fatalf("profile_list: %v", err)
	}

	var decodedList Response
	if err := json.Unmarshal(listResponse, &decodedList); err != nil {
		t.Fatalf("decode profile_list response: %v", err)
	}

	if decodedList.Error != "" {
		t.Fatalf("profile_list error: %s", decodedList.Error)
	}

	rawList, err := json.Marshal(decodedList.Result)
	if err != nil {
		t.Fatalf("marshal profile_list result: %v", err)
	}

	var profiles []profile.Profile
	if err := json.Unmarshal(rawList, &profiles); err != nil {
		t.Fatalf("decode profile_list result: %v", err)
	}

	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}

	if profiles[0].ID != "de-01" {
		t.Fatalf("unexpected profile id: %q", profiles[0].ID)
	}

	getResponse, err := d.HandleJSON(
		ctx,
		[]byte(`{"method":"profile_get","params":{"id":"de-01"}}`),
	)
	if err != nil {
		t.Fatalf("profile_get: %v", err)
	}

	var decodedGet Response
	if err := json.Unmarshal(getResponse, &decodedGet); err != nil {
		t.Fatalf("decode profile_get response: %v", err)
	}

	if decodedGet.Error != "" {
		t.Fatalf("profile_get error: %s", decodedGet.Error)
	}

	rawProfile, err := json.Marshal(decodedGet.Result)
	if err != nil {
		t.Fatalf("marshal profile_get result: %v", err)
	}

	var got profile.Profile
	if err := json.Unmarshal(rawProfile, &got); err != nil {
		t.Fatalf("decode profile_get result: %v", err)
	}

	if got.Name != "Germany 01" {
		t.Fatalf("unexpected profile name: %q", got.Name)
	}

	updatePayload := []byte(`{
		"method":"profile_update",
		"params":{
			"id":"de-01",
			"name":"Germany Updated",
			"type":"vless",
			"server":"example.com",
			"port":8443,
			"enabled":true
		}
	}`)

	if _, err := d.HandleJSON(ctx, updatePayload); err != nil {
		t.Fatalf("profile_update: %v", err)
	}

	deleteResponse, err := d.HandleJSON(
		ctx,
		[]byte(`{"method":"profile_delete","params":{"id":"de-01"}}`),
	)
	if err != nil {
		t.Fatalf("profile_delete: %v", err)
	}

	var decodedDelete Response
	if err := json.Unmarshal(deleteResponse, &decodedDelete); err != nil {
		t.Fatalf("decode profile_delete response: %v", err)
	}

	if decodedDelete.Error != "" {
		t.Fatalf("profile_delete error: %s", decodedDelete.Error)
	}

	finalListResponse, err := d.HandleJSON(
		ctx,
		[]byte(`{"method":"profile_list"}`),
	)
	if err != nil {
		t.Fatalf("profile_list after delete: %v", err)
	}

	var decodedFinalList Response
	if err := json.Unmarshal(finalListResponse, &decodedFinalList); err != nil {
		t.Fatalf("decode final profile_list response: %v", err)
	}

	rawFinalList, err := json.Marshal(decodedFinalList.Result)
	if err != nil {
		t.Fatalf("marshal final profile_list result: %v", err)
	}

	var finalProfiles []profile.Profile
	if err := json.Unmarshal(rawFinalList, &finalProfiles); err != nil {
		t.Fatalf("decode final profile_list result: %v", err)
	}

	if len(finalProfiles) != 0 {
		t.Fatalf("expected empty final profile list, got %d", len(finalProfiles))
	}
}
