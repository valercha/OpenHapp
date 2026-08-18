package uci

import (
	"strings"
	"testing"
)

func TestParseDocumentPreservesSections(t *testing.T) {
	input := `# OpenHapp configuration
config openhapp 'main'
	option engine 'sing-box'
	option mode 'proxy'

config profile 'de-01'
	option name 'Germany'
	option server 'example.com'

config subscription 'sub1'
	option url 'https://example.test/sub'
`

	doc, err := ParseDocument([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if len(doc.Sections) != 3 {
		t.Fatalf("expected 3 sections, got %d", len(doc.Sections))
	}

	if doc.Sections[0].Type != "openhapp" || doc.Sections[0].Name != "main" {
		t.Fatalf("unexpected main section: %+v", doc.Sections[0])
	}

	if doc.Sections[1].Type != "profile" || doc.Sections[1].Name != "de-01" {
		t.Fatalf("unexpected profile section: %+v", doc.Sections[1])
	}

	name, ok := doc.Sections[1].Option("name")
	if !ok || name != "Germany" {
		t.Fatalf("unexpected profile name: %q", name)
	}

	server, ok := doc.Sections[1].Option("server")
	if !ok || server != "example.com" {
		t.Fatalf("unexpected profile server: %q", server)
	}
}

func TestDocumentSetAndDeleteOption(t *testing.T) {
	section := NewSection("profile", "de-01")

	section.SetOption("name", "Germany")
	section.SetOption("server", "example.com")

	name, ok := section.Option("name")
	if !ok || name != "Germany" {
		t.Fatalf("unexpected name: %q", name)
	}

	section.SetOption("name", "Germany Updated")

	name, ok = section.Option("name")
	if !ok || name != "Germany Updated" {
		t.Fatalf("unexpected updated name: %q", name)
	}

	section.DeleteOption("server")

	if _, ok := section.Option("server"); ok {
		t.Fatal("server option should have been deleted")
	}
}

func TestDocumentRenderRoundTrip(t *testing.T) {
	input := `# header

config profile 'de-01'
	option name 'Germany'
	option server 'example.com'

config subscription 'sub1'
	option url 'https://example.test/sub'
`

	doc, err := ParseDocument([]byte(input))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	rendered := string(doc.Render())

	for _, expected := range []string{
		"config profile 'de-01'",
		"option name 'Germany'",
		"option server 'example.com'",
		"config subscription 'sub1'",
		"option url 'https://example.test/sub'",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("missing %q in rendered document:\n%s", expected, rendered)
		}
	}
}
