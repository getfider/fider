package commercial

import (
	"github.com/getfider/fider/app/services"
)

// This file will contain commercial license validation and service registration
// It will be populated in later phases

func init() {
	// TODO: Phase 4 - Register commercial license service
	// TODO: Phase 4 - Register commercial route overrides
	// TODO: Phase 3 - Register commercial service handlers

	// Example (will be implemented later):
	// if isValidCommercialLicense() {
	//     services.License = &commercialLicenseService{}
	//     registerCommercialRoutes()
	//     registerCommercialServices()
	// }
}

// Placeholder for commercial license service (to be implemented)
type commercialLicenseService struct{}

func (s *commercialLicenseService) IsCommercialFeatureEnabled(feature string) bool {
	// TODO: Implement actual license validation
	// For now, return true to enable commercial features during development
	return feature == services.FeatureContentModeration
}

func (s *commercialLicenseService) GetLicenseInfo() *services.LicenseInfo {
	return &services.LicenseInfo{
		IsValid:       true,
		Features:      []string{services.FeatureContentModeration},
		LicenseHolder: "Commercial License Holder",
	}
}