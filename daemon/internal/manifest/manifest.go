package manifest

import (
	"encoding/json"
	"time"

	"github.com/valercha/OpenHapp/daemon/internal/config"
)

// Manifest describes the runtime identity of OpenHapp.
type Manifest struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Engine      string    `json:"engine"`
	Mode        string    `json:"mode"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Default returns the default manifest for the current build.
func Default(version string, cfg config.Config) Manifest {
	engine := cfg.Engine
	if engine == "" {
		engine = "xray"
	}
	mode := cfg.Mode
	if mode == "" {
		mode = "proxy"
	}
	return Manifest{
		Name:        "OpenHapp",
		Version:     version,
		Engine:      engine,
		Mode:        mode,
		Description: "Modern OpenWrt VPN client with LuCI and ubus integration",
		UpdatedAt:   time.Now().UTC(),
	}
}

// FromConfig builds a manifest directly from a runtime config snapshot.
func FromConfig(version string, cfg config.Config) Manifest {
	return Default(version, cfg)
}

// WithTimestamp returns a copy of the manifest with UpdatedAt set to now.
func (m Manifest) WithTimestamp() Manifest {
	m.UpdatedAt = time.Now().UTC()
	return m
}

// JSON returns a stable JSON representation of the manifest.
func (m Manifest) JSON() string {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}
