package strategy

import "testing"

func TestGetDNS(t *testing.T) {
	tests := []struct {
		name       string
		endpoint   OpenObserveHttpEndpoint
		expected   string
		shouldFail bool
	}{
		{
			name:     "Valid URL with HTTPS",
			endpoint: OpenObserveHttpEndpoint("https://api.openobserve.ai/api/ignacio_antonio_organization_31417_mOfaFHQDTo7kIuX"),
			expected: "api.openobserve.ai",
		},
		{
			name:     "Valid URL with HTTP",
			endpoint: OpenObserveHttpEndpoint("http://api.openobserve.ai/api/ignacio_antonio_organization_31417_mOfaFHQDTo7kIuX"),
			expected: "api.openobserve.ai",
		},
		{
			name:     "URL with no path",
			endpoint: OpenObserveHttpEndpoint("https://api.openobserve.ai"),
			expected: "api.openobserve.ai",
		},
		{
			name:     "URL without scheme",
			endpoint: OpenObserveHttpEndpoint("api.openobserve.ai/api/ignacio_antonio_organization_31417_mOfaFHQDTo7kIuX"),
			expected: "api.openobserve.ai",
		},
		{
			name:       "Invalid URL",
			endpoint:   OpenObserveHttpEndpoint("://api.openobserve.ai"),
			expected:   "",
			shouldFail: true,
		},
		{
			name:       "Empty URL",
			endpoint:   OpenObserveHttpEndpoint(""),
			expected:   "",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dns := tt.endpoint.GetDNS()
			if !tt.shouldFail && dns != tt.expected {
				t.Errorf("GetDNS() = %v, want %v", dns, tt.expected)
			}
			if tt.shouldFail && dns != "" {
				t.Errorf("GetDNS() = %v, want empty string for failure", dns)
			}
		})
	}
}
