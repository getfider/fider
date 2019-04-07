package cmd

import "github.com/getfider/fider/app/models"

type UpdateTenantPrivacySettings struct {
	Settings *models.UpdateTenantPrivacy
}

type UpdateTenantSettings struct {
	Settings *models.UpdateTenantSettings
}

type UpdateTenantBillingSettings struct {
	Settings *models.TenantBilling
}

type UpdateTenantAdvancedSettings struct {
	Settings *models.UpdateTenantAdvancedSettings
}

type ActivateTenant struct {
	TenantID int
}
