package config

import "time"

// Config holds the runtime configuration for openhappd.
type Config struct {
	Enabled      bool
	Engine       string
	Mode         string
	LogLevel     string
	Listen       string
	Autostart    bool
	Subscription  string
	UpdatedAt    time.Time
}

// Default returns a sane default configuration.
func Default() Config {
	return Config{
		Enabled:    true,
		Engine:     "xray",
		Mode:       "proxy",
		LogLevel:   "info",
		Listen:     "127.0.0.1:0",
		Autostart:  true,
		Subscription: "",
		UpdatedAt:  time.Now().UTC(),
	}
}
