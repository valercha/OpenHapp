package uci

import (
	"os"
	"testing"

	"github.com/valercha/OpenHapp/daemon/internal/config"
)

func TestStoreRoundTrip(t *testing.T) {
	path := t.TempDir() + "/openhapp"
	store := New(path)

	want := config.Default()
	want.Enabled = false
	want.Autostart = true
	want.Engine = "sing-box"
	want.Mode = "tun"
	want.LogLevel = "debug"
	want.Listen = "127.0.0.1:1080"
	want.Subscription = "https://example.test/sub?a=1&b=2"

	if err := store.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if got.Enabled != want.Enabled ||
		got.Autostart != want.Autostart ||
		got.Engine != want.Engine ||
		got.Mode != want.Mode ||
		got.LogLevel != want.LogLevel ||
		got.Listen != want.Listen ||
		got.Subscription != want.Subscription {
		t.Fatalf("round-trip mismatch: got=%+v want=%+v", got, want)
	}
}

func TestLoadIgnoresOtherSections(t *testing.T) {
	path := t.TempDir() + "/openhapp"
	content := "config openhapp 'main'\n" +
		"\toption engine 'xray'\n" +
		"config other 'main'\n" +
		"\toption engine 'wrong'\n"

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write test config: %v", err)
	}

	got, err := New(path).Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if got.Engine != "xray" {
		t.Fatalf("unexpected engine: %q", got.Engine)
	}
}
