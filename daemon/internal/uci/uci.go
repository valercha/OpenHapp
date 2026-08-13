package uci

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/valercha/OpenHapp/daemon/internal/config"
)

const defaultConfigPath = "/etc/config/openhapp"

type Store struct {
	mu   sync.RWMutex
	path string
}

func New(path string) *Store {
	if strings.TrimSpace(path) == "" {
		path = defaultConfigPath
	}
	return &Store{path: path}
}

func (s *Store) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.path
}

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
	var section string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "config":
			if len(fields) >= 3 && fields[1] == "openhapp" {
				section = strings.Trim(fields[2], "'\"")
			}
		case "option":
			if section != "main" || len(fields) < 3 {
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
	}
	if err := scanner.Err(); err != nil {
		return config.Config{}, fmt.Errorf("scan uci config: %w", err)
	}
	return cfg, nil
}

func escapeUCIValue(v string) string {
	return strings.ReplaceAll(v, "'", "\\'")
}

func (s *Store) Save(cfg config.Config) error {
	s.mu.RLock()
	path := s.path
	s.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}

	var content strings.Builder
	content.WriteString("config openhapp 'main'\n")
	content.WriteString(fmt.Sprintf("\toption enabled '%t'\n", cfg.Enabled))
	content.WriteString(fmt.Sprintf("\toption autostart '%t'\n", cfg.Autostart))
	content.WriteString(fmt.Sprintf("\toption engine '%s'\n", escapeUCIValue(cfg.Engine)))
	content.WriteString(fmt.Sprintf("\toption mode '%s'\n", escapeUCIValue(cfg.Mode)))
	content.WriteString(fmt.Sprintf("\toption log_level '%s'\n", escapeUCIValue(cfg.LogLevel)))
	content.WriteString(fmt.Sprintf("\toption listen '%s'\n", escapeUCIValue(cfg.Listen)))
	content.WriteString(fmt.Sprintf("\toption subscription '%s'\n", escapeUCIValue(cfg.Subscription)))

	return os.WriteFile(path, []byte(content.String()), 0o644)
}
