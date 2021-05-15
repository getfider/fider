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
	IsPrivate bool
}

type UpdateTenantSettings struct {
	Logo           *models.ImageUpload
	Title          string
	Invitation     string
	WelcomeMessage string
	CNAME          string
}

type UpdateTenantAdvancedSettings struct {
	CustomCSS string
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
