package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

// GenerateKey generates a commercial license key for a tenant
// Format: FIDER-COMMERCIAL-{tenantID}-{timestamp}-{hmacHex}
// Panics if LICENSE_MASTER_SECRET is not set
func GenerateKey(tenantID int) string {
	if env.Config.License.MasterSecret == "" {
		panic("LICENSE_MASTER_SECRET environment variable must be set to generate license keys. This is required for hosted Fider instances that sell Pro subscriptions.")
	}
	timestamp := time.Now().Unix()
	data := fmt.Sprintf("%d-%d", tenantID, timestamp)
	mac := hmac.New(sha256.New, []byte(env.Config.License.MasterSecret))
	mac.Write([]byte(data))
	hash := hex.EncodeToString(mac.Sum(nil))
	return fmt.Sprintf("FIDER-COMMERCIAL-%d-%d-%s", tenantID, timestamp, hash)
}

// ValidationResult contains the result of license key validation
type ValidationResult struct {
	IsValid  bool
	TenantID int
	Error    error
}

// ValidateKey validates a commercial license key and returns the tenant ID
// Requires LICENSE_MASTER_SECRET to be set for validation
func ValidateKey(key string) *ValidationResult {
	if key == "" {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.New("license key is empty"),
		}
	}

	if env.Config.License.MasterSecret == "" {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.New("LICENSE_MASTER_SECRET environment variable must be set to validate license keys"),
		}
	}

	parts := strings.Split(key, "-")
	if len(parts) != 5 || parts[0] != "FIDER" || parts[1] != "COMMERCIAL" {
		return &ValidationResult{
			IsValid:  false,
			TenantID: 0,
			Error:    errors.New("invalid key format: expected FIDER-COMMERCIAL-{tenantID}-{timestamp}-{hmac}"),
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

	providedHash := parts[4]

	// Recompute HMAC to verify signature
	data := fmt.Sprintf("%d-%d", tenantID, timestamp)
	mac := hmac.New(sha256.New, []byte(env.Config.License.MasterSecret))
	mac.Write([]byte(data))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	// Constant-time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(expectedHash), []byte(providedHash)) != 1 {
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
