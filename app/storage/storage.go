package storage

import "github.com/getfider/fider/app/models"

// Idea contains read and write operations for ideas
type Idea interface {
	GetByID(ideaID int) (*models.Idea, error)
	GetByNumber(number int) (*models.Idea, error)
	GetCommentsByIdea(number int) ([]*models.Comment, error)
	GetAll() ([]*models.Idea, error)
	Add(title, description string, userID int) (*models.Idea, error)
	Update(number int, title, description string) (*models.Idea, error)
	AddComment(number int, content string, userID int) (int, error)
	AddSupporter(number, userID int) error
	RemoveSupporter(number, userID int) error
	SetResponse(number int, text string, userID, status int) error
	SupportedBy(userID int) ([]int, error)
}

// User is used for user operations
type User interface {
	GetByID(userID int) (*models.User, error)
	GetByEmail(tenantID int, email string) (*models.User, error)
	GetByProvider(tenantID int, provider string, uid string) (*models.User, error)
	Register(user *models.User) error
	RegisterProvider(userID int, provider *models.UserProvider) error
	Update(userID int, settings *models.UpdateUserSettings) error
}

// Tenant contains read and write operations for tenants
type Tenant interface {
	Add(name string, subdomain string) (*models.Tenant, error)
	First() (*models.Tenant, error)
	GetByDomain(domain string) (*models.Tenant, error)
	UpdateSettings(settings *models.UpdateTenantSettings) error
	IsSubdomainAvailable(subdomain string) (bool, error)
	SaveVerificationKey(email, key string) error
	FindVerificationByKey(key string) (*models.SignInRequest, error)
	SetKeyAsVerified(key string) error
	SetCurrentTenant(*models.Tenant)
}
