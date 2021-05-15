package cmd

import (
	"time"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entities"
	"github.com/getfider/fider/app/models/enum"
)

type CreateTenant struct {
	Name      string
	Subdomain string
	Status    int

	Result *entities.Tenant
}

type UpdateTenantPrivacySettings struct {
	IsPrivate bool
}

type UpdateTenantSettings struct {
	Logo           *dto.ImageUpload
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
	Request  NewEmailVerification
}

//NewEmailVerification is used to register a new email verification process
type NewEmailVerification interface {
	GetEmail() string
	GetName() string
	GetUser() *entities.User
	GetKind() enum.EmailVerificationKind
}

type SetKeyAsVerified struct {
	Key string
}
