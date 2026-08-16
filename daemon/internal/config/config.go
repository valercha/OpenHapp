package config

import "time"

// Config holds the runtime configuration for openhappd.
type Config struct {
	Enabled      bool      `json:"enabled"`
	Engine       string    `json:"engine"`
	Mode         string    `json:"mode"`
	LogLevel     string    `json:"log_level"`
	Listen       string    `json:"listen"`
	Autostart    bool      `json:"autostart"`
	Subscription string    `json:"subscription"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Default returns a sane default configuration.
func Default() Config {
	return Config{
		Enabled:      true,
		Engine:       "sing-box",
		Mode:         "proxy",
		LogLevel:     "info",
		Listen:       "127.0.0.1:0",
		Autostart:    true,
		Subscription: "",
		UpdatedAt:    time.Now().UTC(),
	}
}
