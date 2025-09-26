package services

// LicenseService provides license validation for commercial features
type LicenseService interface {
	// IsCommercialFeatureEnabled checks if a specific commercial feature is licensed
	IsCommercialFeatureEnabled(feature string) bool

	// GetLicenseInfo returns information about the current license
	GetLicenseInfo() *LicenseInfo
}

// LicenseInfo contains information about the current license
type LicenseInfo struct {
	IsValid        bool     `json:"isValid"`
	Features       []string `json:"features"`
	ExpiresAt      *string  `json:"expiresAt,omitempty"`
	LicenseHolder  string   `json:"licenseHolder"`
}

// Commercial feature constants
const (
	FeatureContentModeration = "content-moderation"
)

// Default implementation that denies all commercial features
type defaultLicenseService struct{}

func (s *defaultLicenseService) IsCommercialFeatureEnabled(feature string) bool {
	return false
}

func (s *defaultLicenseService) GetLicenseInfo() *LicenseInfo {
	return &LicenseInfo{
		IsValid:       false,
		Features:      []string{},
		LicenseHolder: "Open Source",
	}
}

// Global license service instance - can be overridden by commercial code
var License LicenseService = &defaultLicenseService{}

// Helper function for easy access
func IsCommercialFeatureEnabled(feature string) bool {
	return License.IsCommercialFeatureEnabled(feature)
}