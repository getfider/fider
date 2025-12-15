package license_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/license"
)

func TestGenerateAndValidateKey_Ed25519(t *testing.T) {
	// Generate a test key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Set up environment with Ed25519 keys
	env.Config.License.PrivateKey = base64.StdEncoding.EncodeToString(privateKey)
	env.Config.License.PublicKey = base64.StdEncoding.EncodeToString(publicKey)
	defer func() {
		env.Config.License.PrivateKey = ""
		env.Config.License.PublicKey = ""
	}()

	// Generate a license key
	tenantID := 42
	key := license.GenerateKey(tenantID)

	// Validate the license key
	result := license.ValidateKey(key)

	if !result.IsValid {
		t.Errorf("Expected license key to be valid, but got error: %v", result.Error)
	}

	if result.TenantID != tenantID {
		t.Errorf("Expected tenant ID %d, got %d", tenantID, result.TenantID)
	}
}

func TestValidateKey_Ed25519_InvalidSignature(t *testing.T) {
	// Generate a test key pair
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	env.Config.License.PrivateKey = base64.StdEncoding.EncodeToString(privateKey)
	env.Config.License.PublicKey = base64.StdEncoding.EncodeToString(publicKey)
	defer func() {
		env.Config.License.PrivateKey = ""
		env.Config.License.PublicKey = ""
	}()

	// Generate a valid license key
	key := license.GenerateKey(42)

	// Tamper with the signature (last part of the key)
	tamperedKey := key[:len(key)-4] + "XXXX"

	// Validate the tampered key
	result := license.ValidateKey(tamperedKey)

	if result.IsValid {
		t.Error("Expected tampered license key to be invalid")
	}

	if result.Error == nil {
		t.Error("Expected validation error for tampered key")
	}
}

func TestValidateKey_Ed25519_WrongPublicKey(t *testing.T) {
	// Generate two different key pairs
	_, privateKey1, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	publicKey2, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	// Generate key with one private key
	env.Config.License.PrivateKey = base64.StdEncoding.EncodeToString(privateKey1)
	env.Config.License.PublicKey = ""
	defer func() {
		env.Config.License.PrivateKey = ""
		env.Config.License.PublicKey = ""
	}()

	key := license.GenerateKey(42)

	// Try to validate with a different public key
	env.Config.License.PrivateKey = ""
	env.Config.License.PublicKey = base64.StdEncoding.EncodeToString(publicKey2)

	result := license.ValidateKey(key)

	if result.IsValid {
		t.Error("Expected license key to be invalid when validated with wrong public key")
	}
}

func TestValidateKey_EmptyKey(t *testing.T) {
	result := license.ValidateKey("")

	if result.IsValid {
		t.Error("Expected empty key to be invalid")
	}

	if result.Error == nil {
		t.Error("Expected validation error for empty key")
	}
}

func TestValidateKey_InvalidFormat(t *testing.T) {
	testCases := []string{
		"INVALID-KEY",
		"FIDER-INVALID-123-456-789",
		"FIDER-COMMERCIAL-123", // Too few parts
		"FIDER-COMMERCIAL-abc-456-789", // Invalid tenant ID
		"FIDER-COMMERCIAL-123-abc-789", // Invalid timestamp
	}

	for _, key := range testCases {
		result := license.ValidateKey(key)

		if result.IsValid {
			t.Errorf("Expected key '%s' to be invalid", key)
		}

		if result.Error == nil {
			t.Errorf("Expected validation error for key '%s'", key)
		}
	}
}

func TestValidateKey_UsesEmbeddedDefault(t *testing.T) {
	// Clear LICENSE_PUBLIC_KEY to test embedded default
	originalPrivateKey := env.Config.License.PrivateKey
	originalPublicKey := env.Config.License.PublicKey

	// Set private key to match the embedded public key (lPHzCZhOBBihIusKWs5lzXCgGxZEKBpCiplkmZSjGpU=)
	// This is the corresponding private key for the DefaultPublicKey
	env.Config.License.PrivateKey = "B0Srq+DtsBI6MwY4sljCywXNLWAizfCdiJ6emGx+xIqU8fMJmE4EGKEi6wpazmXNcKAbFkQoGkKKmWSZlKMalQ=="
	env.Config.License.PublicKey = "" // Leave empty to test embedded default

	defer func() {
		env.Config.License.PrivateKey = originalPrivateKey
		env.Config.License.PublicKey = originalPublicKey
	}()

	// Generate a key with the private key
	key := license.GenerateKey(123)

	// Validate without setting LICENSE_PUBLIC_KEY - should use embedded default
	result := license.ValidateKey(key)

	if !result.IsValid {
		t.Errorf("Expected license key to be valid using embedded default public key, but got error: %v", result.Error)
	}

	if result.TenantID != 123 {
		t.Errorf("Expected tenant ID 123, got %d", result.TenantID)
	}
}
