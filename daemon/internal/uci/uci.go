package uci

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/valercha/OpenHapp/daemon/internal/config"
)

const defaultConfigPath = "/etc/config/openhapp"

// Store provides a minimal UCI-backed persistence layer for OpenHapp.
type Store struct {
	mu   sync.RWMutex
	path string
}

// New creates a new persistence store.
func New(path string) *Store {
	if strings.TrimSpace(path) == "" {
		path = defaultConfigPath
	}
	return &Store{path: path}
}

// Path returns the configured UCI file path.
func (s *Store) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.path
}

// Load reads the runtime config from UCI if the file exists.
func (s *Store) Load() (config.Config, error) {
	s.mu.RLock()
	path := s.path
	s.mu.RUnlock()

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config.Default(), nil
		}
		return config.Config{}, fmt.Errorf("open uci config: %w", err)
	}
	defer f.Close()

	cfg := config.Default()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		if fields[0] != "option" {
			continue
		}
		key := fields[1]
		value := strings.Trim(strings.Join(fields[2:], " "), "'\"")
		switch key {
		case "enabled":
			cfg.Enabled = value == "1" || strings.EqualFold(value, "true")
		case "autostart":
			cfg.Autostart = value == "1" || strings.EqualFold(value, "true")
		case "engine":
			cfg.Engine = value
		case "mode":
			cfg.Mode = value
		case "log_level":
			cfg.LogLevel = value
		case "listen":
			cfg.Listen = value
		case "subscription":
			cfg.Subscription = value
		}
	}
	if err := scanner.Err(); err != nil {
		return config.Config{}, fmt.Errorf("scan uci config: %w", err)
	}
	return cfg, nil
}

// Save persists the runtime config into a minimal UCI file.
func (s *Store) Save(cfg config.Config) error {
	s.mu.RLock()
	path := s.path
	s.mu.RUnlock()

	if err := os.MkdirAll("/etc/config", 0o755); err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}

	content := strings.Builder{}
	content.WriteString("config openhapp 'main'\n")
	content.WriteString(fmt.Sprintf("\toption enabled '%t'\n", cfg.Enabled))
	content.WriteString(fmt.Sprintf("\toption autostart '%t'\n", cfg.Autostart))
	content.WriteString(fmt.Sprintf("\toption engine '%s'\n", cfg.Engine))
	content.WriteString(fmt.Sprintf("\toption mode '%s'\n", cfg.Mode))
	content.WriteString(fmt.Sprintf("\toption log_level '%s'\n", cfg.LogLevel))
	content.WriteString(fmt.Sprintf("\toption listen '%s'\n", cfg.Listen))
	content.WriteString(fmt.Sprintf("\toption subscription '%s'\n", cfg.Subscription))

	return os.WriteFile(path, []byte(content.String()), 0o644)
}
