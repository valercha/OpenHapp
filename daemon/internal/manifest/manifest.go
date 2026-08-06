package manifest

import (
	"encoding/json"
	"time"
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
func Default(version, engine, mode string) Manifest {
	if engine == "" {
		engine = "xray"
	}
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

// JSON returns a stable JSON representation of the manifest.
func (m Manifest) JSON() string {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}
