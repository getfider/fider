package cmd

import (
	"time"

	"github.com/getfider/fider/app/models"
)

type CreateTenant struct {
	Name      string
	Subdomain string
	Status    int

	Result *models.Tenant
}

type UpdateTenantPrivacySettings struct {
	Settings *models.UpdateTenantPrivacy
}

type UpdateTenantSettings struct {
	Settings *models.UpdateTenantSettings
}

type UpdateTenantAdvancedSettings struct {
	Settings *models.UpdateTenantAdvancedSettings
}

type ActivateTenant struct {
	TenantID int
}

type SaveVerificationKey struct {
	Key      string
	Duration time.Duration
	Request  models.NewEmailVerification
}

type SetKeyAsVerified struct {
	Key string
}
