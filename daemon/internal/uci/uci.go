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
			section = ""
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
			case "ownership":
				cfg.Ownership = value
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

func mainSectionHeader(line string) bool {
	line = strings.TrimSpace(line)
	return line == "config openhapp 'main'" ||
		line == `config openhapp "main"`
}

func renderMainConfig(cfg config.Config) string {
	var content strings.Builder

	content.WriteString("config openhapp 'main'\n")

	enabled := "0"
	if cfg.Enabled {
		enabled = "1"
	}

	autostart := "0"
	if cfg.Autostart {
		autostart = "1"
	}

	content.WriteString(fmt.Sprintf("\toption enabled '%s'\n", enabled))
	content.WriteString(fmt.Sprintf("\toption autostart '%s'\n", autostart))
	content.WriteString(fmt.Sprintf("\toption engine '%s'\n", escapeUCIValue(cfg.Engine)))
	content.WriteString(fmt.Sprintf("\toption ownership '%s'\n", escapeUCIValue(cfg.Ownership)))
	content.WriteString(fmt.Sprintf("\toption mode '%s'\n", escapeUCIValue(cfg.Mode)))
	content.WriteString(fmt.Sprintf("\toption log_level '%s'\n", escapeUCIValue(cfg.LogLevel)))
	content.WriteString(fmt.Sprintf("\toption listen '%s'\n", escapeUCIValue(cfg.Listen)))
	content.WriteString(fmt.Sprintf("\toption subscription '%s'\n", escapeUCIValue(cfg.Subscription)))

	return content.String()
}

func (s *Store) Save(cfg config.Config) error {
	s.mu.RLock()
	path := s.path
	s.mu.RUnlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}

	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read uci config: %w", err)
	}

	mainConfig := renderMainConfig(cfg)

	var content strings.Builder
	lines := strings.Split(string(existing), "\n")

	inMain := false
	wroteMain := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "config ") {
			inMain = mainSectionHeader(trimmed)

			if inMain {
				if !wroteMain {
					content.WriteString(mainConfig)
					wroteMain = true
				}
				continue
			}
		}

		if inMain {
			continue
		}

		if line != "" || content.Len() > 0 {
			content.WriteString(line)
			content.WriteByte('\n')
		}
	}

	if !wroteMain {
		if content.Len() > 0 && !strings.HasSuffix(content.String(), "\n\n") {
			content.WriteByte('\n')
		}
		content.WriteString(mainConfig)
	}

	return os.WriteFile(path, []byte(content.String()), 0o644)
}
