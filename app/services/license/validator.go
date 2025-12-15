package license

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

// DefaultPublicKey is the embedded public key for license validation
// This is the public key from the hosted Fider instance
// Self-hosted instances automatically use this - no configuration needed
// Can be overridden via LICENSE_PUBLIC_KEY environment variable
const DefaultPublicKey = "EpyoY4Fc3TroE7MIEWlLHU8OGaEiPkPhOy+RVwwC1zk="

// GenerateKey generates a commercial license key for a tenant using Ed25519 signatures
// Format: FIDER-COMMERCIAL-{tenantID}-{timestamp}-{signature}
// Panics if LICENSE_PRIVATE_KEY is not set
func GenerateKey(tenantID int) string {
	if env.Config.License.PrivateKey == "" {
		panic("LICENSE_PRIVATE_KEY environment variable must be set to generate license keys. This is required for hosted Fider instances that sell Pro subscriptions.")
	}

	privateKeyBytes, err := base64.StdEncoding.DecodeString(env.Config.License.PrivateKey)
	if err != nil {
		panic(fmt.Sprintf("Invalid LICENSE_PRIVATE_KEY: %v", err))
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		panic(fmt.Sprintf("Invalid LICENSE_PRIVATE_KEY length: expected %d bytes, got %d", ed25519.PrivateKeySize, len(privateKeyBytes)))
	}

	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%d-%d", tenantID, timestamp)
	signature := ed25519.Sign(ed25519.PrivateKey(privateKeyBytes), []byte(data))
	signatureHex := hex.EncodeToString(signature)

	return fmt.Sprintf("FIDER-COMMERCIAL-%d-%d-%s", tenantID, timestamp, signatureHex)
}

// ValidationResult contains the result of license key validation
type ValidationResult struct {
	IsValid  bool
	TenantID int
	Error    error
}

// ValidateKey validates a commercial license key using Ed25519 signature verification
// Uses embedded DefaultPublicKey unless LICENSE_PUBLIC_KEY environment variable is set
func ValidateKey(key string) *ValidationResult {
	if key == "" {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.New("license key is empty"),
		}
	}

	// Use environment variable if set, otherwise use embedded default
	publicKey := env.Config.License.PublicKey
	if publicKey == "" {
		publicKey = DefaultPublicKey
	}

	// Parse license key format
	parts := strings.Split(key, "-")
	if len(parts) != 5 || parts[0] != "FIDER" || parts[1] != "COMMERCIAL" {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.New("invalid key format: expected FIDER-COMMERCIAL-{tenantID}-{timestamp}-{signature}"),
		}
	}

	tenantID, err := strconv.Atoi(parts[2])
	if err != nil {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.Wrap(err, "invalid tenant ID in license key"),
		}
	}

	timestamp, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.Wrap(err, "invalid timestamp in license key"),
		}
	}

	providedSignatureHex := parts[4]

	// Decode and validate the public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.Wrap(err, "invalid public key configuration"),
		}
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    fmt.Errorf("invalid public key length: expected %d bytes, got %d", ed25519.PublicKeySize, len(publicKeyBytes)),
		}
	}

	// Decode the signature from the license key
	providedSignature, err := hex.DecodeString(providedSignatureHex)
	if err != nil {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.Wrap(err, "invalid signature format in license key"),
		}
	}

	// Verify the signature
	data := fmt.Sprintf("%d-%d", tenantID, timestamp)
	if !ed25519.Verify(ed25519.PublicKey(publicKeyBytes), []byte(data), providedSignature) {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.New("invalid signature: license key verification failed"),
		}
	}

	return &ValidationResult{
		IsValid:  true,
		TenantID: tenantID,
		Error:    nil,
	}
}
