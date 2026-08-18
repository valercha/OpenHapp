package profile

import (
	"os"
	"strings"
	"testing"
)

func TestStoreCRUD(t *testing.T) {
	path := t.TempDir() + "/openhapp"
	store := NewStore(path)

	profile := Profile{
		ID:      "de-01",
		Name:    "Germany 01",
		Type:    "vless",
		Server:  "example.com",
		Port:    443,
		Enabled: true,
		Properties: map[string]string{
			"uuid": "test-uuid",
			"sni":  "example.com",
		},
	}

	if err := store.Save(profile); err != nil {
		t.Fatalf("save: %v", err)
	}

	got, err := store.Get("de-01")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got.ID != profile.ID ||
		got.Name != profile.Name ||
		got.Type != profile.Type ||
		got.Server != profile.Server ||
		got.Port != profile.Port ||
		got.Enabled != profile.Enabled {
		t.Fatalf("unexpected profile: got=%+v want=%+v", got, profile)
	}

	if got.Properties["uuid"] != "test-uuid" ||
		got.Properties["sni"] != "example.com" {
		t.Fatalf("unexpected properties: %+v", got.Properties)
	}

	profiles, err := store.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}

	profile.Name = "Germany Updated"
	profile.Port = 8443

	if err := store.Save(profile); err != nil {
		t.Fatalf("update: %v", err)
	}

	updated, err := store.Get("de-01")
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}

	if updated.Name != "Germany Updated" || updated.Port != 8443 {
		t.Fatalf("unexpected updated profile: %+v", updated)
	}

	if err := store.Delete("de-01"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	if _, err := store.Get("de-01"); err == nil {
		t.Fatal("expected deleted profile to be unavailable")
	}
}

func TestStorePreservesOtherSections(t *testing.T) {
	path := t.TempDir() + "/openhapp"

	initial := `config openhapp 'main'
	option engine 'sing-box'
	option mode 'proxy'

config subscription 'sub1'
	option name 'Example'
	option url 'https://example.test/sub'

config profile 'old'
	option name 'Old'
	option type 'vless'
	option server 'old.example.com'
	option port '443'
`

	if err := os.WriteFile(path, []byte(initial), 0o644); err != nil {
		t.Fatalf("write initial config: %v", err)
	}

	store := NewStore(path)

	profile := Profile{
		ID:      "de-01",
		Name:    "Germany",
		Type:    "vless",
		Server:  "example.com",
		Port:    443,
		Enabled: true,
	}

	if err := store.Save(profile); err != nil {
		t.Fatalf("save: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	saved := string(content)

	for _, expected := range []string{
		"config openhapp 'main'",
		"option engine 'sing-box'",
		"config subscription 'sub1'",
		"option url 'https://example.test/sub'",
		"config profile 'old'",
		"config profile 'de-01'",
	} {
		if !strings.Contains(saved, expected) {
			t.Fatalf("expected preserved content %q, got:\n%s", expected, saved)
		}
	}
}

func TestStoreListSortsByID(t *testing.T) {
	path := t.TempDir() + "/openhapp"
	store := NewStore(path)

	profiles := []Profile{
		{
			ID:      "de-02",
			Name:    "Germany 02",
			Type:    "vless",
			Server:  "two.example.com",
			Port:    443,
			Enabled: true,
		},
		{
			ID:      "de-01",
			Name:    "Germany 01",
			Type:    "vless",
			Server:  "one.example.com",
			Port:    443,
			Enabled: true,
		},
	}

	for _, profile := range profiles {
		if err := store.Save(profile); err != nil {
			t.Fatalf("save %s: %v", profile.ID, err)
		}
	}

	got, err := store.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(got))
	}

	if got[0].ID != "de-01" || got[1].ID != "de-02" {
		t.Fatalf("unexpected order: %+v", got)
	}
}
