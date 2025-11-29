package commercial

import (
	"context"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/services"
	"github.com/getfider/fider/app/services/license"
	commercialHandlers "github.com/getfider/fider/commercial/handlers"
	commercialApiv1 "github.com/getfider/fider/commercial/handlers/apiv1"
)

// Commercial license service implementation
type commercialLicenseService struct {
	isValid  bool
	tenantID int
}

func init() {
	ctx := context.Background()
	var svc *commercialLicenseService

	if env.IsSingleHostMode() && env.Config.License.Key != "" {
		// Self-hosted with license key: validate at startup
		result := license.ValidateKey(env.Config.License.Key)
		if !result.IsValid {
			panic(result.Error)
		}

		svc = &commercialLicenseService{
			isValid:  true,
			tenantID: result.TenantID,
		}

		log.Infof(ctx, "Commercial license validated for tenant @{TenantID}", dto.Props{
			"TenantID": result.TenantID,
		})
	} else {
		// Multi-tenant hosted OR self-hosted without key
		svc = &commercialLicenseService{isValid: env.IsMultiHostMode()}
	}

	// Register commercial license service
	services.License = svc

	// Always register commercial handlers (they check licensing internally)
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

func (s *commercialLicenseService) IsCommercialFeatureEnabled(feature string) bool {
	// In self-hosted mode, check license validity
	if env.IsSingleHostMode() {
		return s.isValid && feature == services.FeatureContentModeration
	}

	// In multi-tenant mode, this is a global check
	// Per-tenant checking happens via tenant.IsCommercial()
	return true
}

func (s *commercialLicenseService) GetLicenseInfo() *services.LicenseInfo {
	if s.isValid {
		return &services.LicenseInfo{
			IsValid:       true,
			Features:      []string{services.FeatureContentModeration},
			LicenseHolder: "Licensed",
		}
	}

	return &services.LicenseInfo{
		IsValid:       false,
		Features:      []string{},
		LicenseHolder: "Unlicensed",
	}
}
