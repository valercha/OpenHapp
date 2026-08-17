package uci

import (
	"os"
	"strings"
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
	want.Ownership = "external"
	want.Mode = "tun"
	want.LogLevel = "debug"
	want.Listen = "127.0.0.1:1080"
	want.Subscription = "https://example.test/sub?a=1&b=2"

	if err := store.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved config: %v", err)
	}

	saved := string(content)

	if !strings.Contains(saved, "option enabled '0'") {
		t.Fatalf("expected UCI enabled flag to be 0, got:\n%s", saved)
	}

	if !strings.Contains(saved, "option autostart '1'") {
		t.Fatalf("expected UCI autostart flag to be 1, got:\n%s", saved)
	}

	if !strings.Contains(saved, "option ownership 'external'") {
		t.Fatalf("expected UCI ownership to be external, got:\n%s", saved)
	}

	got, err := store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if got.Enabled != want.Enabled ||
		got.Autostart != want.Autostart ||
		got.Engine != want.Engine ||
		got.Ownership != want.Ownership ||
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

func TestStoreSaveAddsMainWithoutRemovingOtherSections(t *testing.T) {
	path := t.TempDir() + "/openhapp"
	store := New(path)

	initial := `config profile 'de-01'
	option name 'Germany'
	option server 'example.com'

config subscription 'sub1'
	option name 'Example'
	option url 'https://example.test/sub'
`

	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial config: %v", err)
	}

	cfg := config.Default()
	cfg.LogLevel = "debug"

	if err := store.Save(cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved config: %v", err)
	}

	saved := string(content)

	for _, expected := range []string{
		"config openhapp 'main'",
		"option log_level 'debug'",
		"config profile 'de-01'",
		"option name 'Germany'",
		"config subscription 'sub1'",
		"option url 'https://example.test/sub'",
	} {
		if !strings.Contains(saved, expected) {
			t.Fatalf("expected content %q, got:\n%s", expected, saved)
		}
	}
}

func TestStoreSavePreservesOtherSections(t *testing.T) {
	path := t.TempDir() + "/openhapp"
	store := New(path)

	initial := `config openhapp 'main'
	option enabled '1'
	option autostart '1'
	option engine 'sing-box'
	option ownership 'external'
	option mode 'proxy'
	option log_level 'info'
	option listen '127.0.0.1:0'
	option subscription ''

config profile 'de-01'
	option name 'Germany'
	option server 'example.com'
	option port '443'

config subscription 'sub1'
	option name 'Example'
	option url 'https://example.test/sub'
`

	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial config: %v", err)
	}

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	cfg.LogLevel = "debug"

	if err := store.Save(cfg); err != nil {
		t.Fatalf("save: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved config: %v", err)
	}

	saved := string(content)

	for _, expected := range []string{
		"config profile 'de-01'",
		"option name 'Germany'",
		"option server 'example.com'",
		"option port '443'",
		"config subscription 'sub1'",
		"option name 'Example'",
		"option url 'https://example.test/sub'",
	} {
		if !strings.Contains(saved, expected) {
			t.Fatalf("expected preserved content %q, got:\n%s", expected, saved)
		}
	}

	if !strings.Contains(saved, "option log_level 'debug'") {
		t.Fatalf("expected updated main config, got:\n%s", saved)
	}
}
