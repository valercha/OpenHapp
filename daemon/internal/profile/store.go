package profile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/valercha/OpenHapp/daemon/internal/uci"
)

const defaultConfigPath = "/etc/config/openhapp"

const profileSectionType = "profile"

// Store persists normalized profiles as UCI profile sections.
type Store struct {
	mu   sync.RWMutex
	path string
}

// NewStore creates a profile store backed by an OpenHapp UCI config.
func NewStore(path string) *Store {
	if strings.TrimSpace(path) == "" {
		path = defaultConfigPath
	}

	return &Store{path: path}
}

// Path returns the backing UCI config path.
func (s *Store) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.path
}

// List returns all stored profiles sorted by ID.
func (s *Store) List() ([]Profile, error) {
	doc, err := s.loadDocument()
	if err != nil {
		return nil, err
	}

	profiles := make([]Profile, 0)

	for _, section := range doc.Sections {
		if section.Type != profileSectionType {
			continue
		}

		profile, err := profileFromSection(section)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, profile)
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].ID < profiles[j].ID
	})

	return profiles, nil
}

// Get returns a profile by ID.
func (s *Store) Get(id string) (Profile, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Profile{}, fmt.Errorf("profile id is empty")
	}

	doc, err := s.loadDocument()
	if err != nil {
		return Profile{}, err
	}

	for _, section := range doc.Sections {
		if section.Type != profileSectionType || section.Name != id {
			continue
		}

		return profileFromSection(section)
	}

	return Profile{}, fmt.Errorf("profile %q not found", id)
}

// Save creates or replaces a profile.
func (s *Store) Save(profile Profile) error {
	if err := profile.Validate(); err != nil {
		return err
	}

	doc, err := s.loadDocument()
	if err != nil {
		return err
	}

	found := false

	for i := range doc.Sections {
		section := &doc.Sections[i]

		if section.Type != profileSectionType || section.Name != profile.ID {
			continue
		}

		*section = sectionFromProfile(profile)
		found = true
		break
	}

	if !found {
		doc.Sections = append(doc.Sections, sectionFromProfile(profile))
	}

	return s.saveDocument(doc)
}

// Delete removes a profile by ID.
func (s *Store) Delete(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("profile id is empty")
	}

	doc, err := s.loadDocument()
	if err != nil {
		return err
	}

	found := false
	filtered := doc.Sections[:0]

	for _, section := range doc.Sections {
		if section.Type == profileSectionType && section.Name == id {
			found = true
			continue
		}

		filtered = append(filtered, section)
	}

	if !found {
		return fmt.Errorf("profile %q not found", id)
	}

	doc.Sections = filtered

	return s.saveDocument(doc)
}

func profileFromSection(section uci.Section) (Profile, error) {
	profile := Profile{
		ID:         section.Name,
		Properties: map[string]string{},
	}

	var ok bool

	if profile.Name, ok = section.Option("name"); !ok {
		return Profile{}, fmt.Errorf("profile %q missing name", section.Name)
	}

	if profile.Type, ok = section.Option("type"); !ok {
		return Profile{}, fmt.Errorf("profile %q missing type", section.Name)
	}

	if profile.Server, ok = section.Option("server"); !ok {
		return Profile{}, fmt.Errorf("profile %q missing server", section.Name)
	}

	portValue, ok := section.Option("port")
	if !ok {
		return Profile{}, fmt.Errorf("profile %q missing port", section.Name)
	}

	var port int
	if _, err := fmt.Sscanf(portValue, "%d", &port); err != nil {
		return Profile{}, fmt.Errorf("profile %q has invalid port %q", section.Name, portValue)
	}

	profile.Port = port

	enabledValue, ok := section.Option("enabled")
	if ok {
		profile.Enabled = enabledValue == "1" ||
			strings.EqualFold(enabledValue, "true")
	} else {
		profile.Enabled = true
	}

	for _, line := range section.Lines {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 3 || fields[0] != "option" {
			continue
		}

		key := fields[1]
		if key == "name" ||
			key == "type" ||
			key == "server" ||
			key == "port" ||
			key == "enabled" {
			continue
		}

		value := strings.Trim(strings.Join(fields[2:], " "), "'\"")
		profile.Properties[key] = value
	}

	if err := profile.Validate(); err != nil {
		return Profile{}, err
	}

	if len(profile.Properties) == 0 {
		profile.Properties = nil
	}

	return profile, nil
}

func sectionFromProfile(profile Profile) uci.Section {
	section := uci.NewSection(profileSectionType, profile.ID)

	section.SetOption("name", profile.Name)
	section.SetOption("type", profile.Type)
	section.SetOption("server", profile.Server)
	section.SetOption("port", fmt.Sprintf("%d", profile.Port))

	if profile.Enabled {
		section.SetOption("enabled", "1")
	} else {
		section.SetOption("enabled", "0")
	}

	keys := make([]string, 0, len(profile.Properties))
	for key := range profile.Properties {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		if key == "" ||
			key == "name" ||
			key == "type" ||
			key == "server" ||
			key == "port" ||
			key == "enabled" {
			continue
		}

		section.SetOption(key, profile.Properties[key])
	}

	return section
}

func (s *Store) loadDocument() (uci.Document, error) {
	s.mu.RLock()
	path := s.path
	s.mu.RUnlock()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return uci.Document{}, nil
		}

		return uci.Document{}, fmt.Errorf("read profile config: %w", err)
	}

	doc, err := uci.ParseDocument(data)
	if err != nil {
		return uci.Document{}, fmt.Errorf("parse profile config: %w", err)
	}

	return doc, nil
}

func (s *Store) saveDocument(doc uci.Document) error {
	s.mu.RLock()
	path := s.path
	s.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}

	if err := os.WriteFile(path, doc.Render(), 0o644); err != nil {
		return fmt.Errorf("write profile config: %w", err)
	}

	return nil
}
