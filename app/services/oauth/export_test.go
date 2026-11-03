package oauth

import "testing"

// SetSystemProviderStatus overrides the instance-level status of a built-in
// provider for the duration of the test (test-only hook).
func SetSystemProviderStatus(t *testing.T, provider string, status int) {
	for _, p := range systemProviders {
		if p.Provider == provider {
			original := p.Status
			p.Status = status
			t.Cleanup(func() { p.Status = original })
			return
		}
	}
	t.Fatalf("unknown system provider %s", provider)
}
