package profile

import "testing"

func TestProfileValidate(t *testing.T) {
	valid := Profile{
		ID:      "de-01",
		Name:    "Germany 01",
		Type:    "vless",
		Server:  "example.com",
		Port:    443,
		Enabled: true,
	}

	if err := valid.Validate(); err != nil {
		t.Fatalf("valid profile rejected: %v", err)
	}
}

func TestProfileValidateRejectsMissingFields(t *testing.T) {
	tests := []struct {
		name    string
		profile Profile
	}{
		{
			name: "missing id",
			profile: Profile{
				Name:   "Germany",
				Type:   "vless",
				Server: "example.com",
				Port:   443,
			},
		},
		{
			name: "missing name",
			profile: Profile{
				ID:     "de-01",
				Type:   "vless",
				Server: "example.com",
				Port:   443,
			},
		},
		{
			name: "missing type",
			profile: Profile{
				ID:     "de-01",
				Name:   "Germany",
				Server: "example.com",
				Port:   443,
			},
		},
		{
			name: "missing server",
			profile: Profile{
				ID:   "de-01",
				Name: "Germany",
				Type: "vless",
				Port: 443,
			},
		},
		{
			name: "invalid port",
			profile: Profile{
				ID:     "de-01",
				Name:   "Germany",
				Type:   "vless",
				Server: "example.com",
				Port:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.profile.Validate(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestProfileValidatePortBoundaries(t *testing.T) {
	valid := Profile{
		ID:     "de-01",
		Name:   "Germany",
		Type:   "vless",
		Server: "example.com",
	}

	t.Run("minimum valid port", func(t *testing.T) {
		profile := valid
		profile.Port = 1

		if err := profile.Validate(); err != nil {
			t.Fatalf("minimum valid port rejected: %v", err)
		}
	})

	t.Run("maximum valid port", func(t *testing.T) {
		profile := valid
		profile.Port = 65535

		if err := profile.Validate(); err != nil {
			t.Fatalf("maximum valid port rejected: %v", err)
		}
	})

	t.Run("port below minimum", func(t *testing.T) {
		profile := valid
		profile.Port = 0

		if err := profile.Validate(); err == nil {
			t.Fatal("expected validation error for port 0")
		}
	})

	t.Run("port above maximum", func(t *testing.T) {
		profile := valid
		profile.Port = 65536

		if err := profile.Validate(); err == nil {
			t.Fatal("expected validation error for port 65536")
		}
	})
}
