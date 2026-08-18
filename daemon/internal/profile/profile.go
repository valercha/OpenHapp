package profile

import (
	"fmt"
	"strings"
)

// Profile represents a normalized proxy endpoint.
// It is intentionally independent from UCI and transport-specific URL formats.
type Profile struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Server     string            `json:"server"`
	Port       int               `json:"port"`
	Enabled    bool              `json:"enabled"`
	Properties map[string]string `json:"properties,omitempty"`
}

// Validate checks the minimum fields required for a usable profile.
func (p Profile) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("profile id is empty")
	}

	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("profile name is empty")
	}

	if strings.TrimSpace(p.Type) == "" {
		return fmt.Errorf("profile type is empty")
	}

	if strings.TrimSpace(p.Server) == "" {
		return fmt.Errorf("profile server is empty")
	}

	if p.Port < 1 || p.Port > 65535 {
		return fmt.Errorf("profile port must be between 1 and 65535")
	}

	return nil
}
