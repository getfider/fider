package commercial

import (
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/services"
	commercialHandlers "github.com/getfider/fider/commercial/handlers"
	commercialApiv1 "github.com/getfider/fider/commercial/handlers/apiv1"
)

func init() {
	// Register commercial license service
	services.License = &commercialLicenseService{}

	// Register commercial HTTP handlers
	handlers.RegisterModerationHandlers(
		commercialHandlers.ModerationPage,
		commercialHandlers.GetModerationItems,
		commercialHandlers.GetModerationCount,
	)

	apiv1.RegisterModerationHandlers(
		commercialApiv1.ApprovePost,
		commercialApiv1.DeclinePost,
		commercialApiv1.ApproveComment,
		commercialApiv1.DeclineComment,
		commercialApiv1.DeclinePostAndBlock,
		commercialApiv1.DeclineCommentAndBlock,
		commercialApiv1.ApprovePostAndVerify,
		commercialApiv1.ApproveCommentAndVerify,
	)
}

// Commercial license service implementation
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