package storage

import (
	"time"

	"github.com/getfider/fider/app/models"
)

// Base is a generic storage base interface
type Base interface {
	SetCurrentTenant(*models.Tenant)
	SetCurrentUser(*models.User)
}

// Idea contains read and write operations for ideas
type Idea interface {
	Base
	GetByID(ideaID int) (*models.Idea, error)
	GetBySlug(slug string) (*models.Idea, error)
	GetByNumber(number int) (*models.Idea, error)
	GetCommentsByIdea(idea *models.Idea) ([]*models.Comment, error)
	Search(query, filter, limit string, tags []string) ([]*models.Idea, error)
	GetAll() ([]*models.Idea, error)
	CountPerStatus() (map[int]int, error)
	Add(title, description string) (*models.Idea, error)
	Update(idea *models.Idea, title, description string) (*models.Idea, error)
	AddComment(idea *models.Idea, content string) (int, error)
	GetCommentByID(id int) (*models.Comment, error)
	UpdateComment(id int, content string) error
	AddSupporter(idea *models.Idea, user *models.User) error
	RemoveSupporter(idea *models.Idea, user *models.User) error
	AddSubscriber(idea *models.Idea, user *models.User) error
	RemoveSubscriber(idea *models.Idea, user *models.User) error
	GetActiveSubscribers(number int, channel models.NotificationChannel, event models.NotificationEvent) ([]*models.User, error)
	SetResponse(idea *models.Idea, text string, status int) error
	MarkAsDuplicate(idea *models.Idea, original *models.Idea) error
	IsReferenced(idea *models.Idea) (bool, error)
	SupportedBy() ([]int, error)
}

// User is used for user operations
type User interface {
	Base
	GetByID(userID int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByProvider(provider string, uid string) (*models.User, error)
	Register(user *models.User) error
	RegisterProvider(userID int, provider *models.UserProvider) error
	Update(settings *models.UpdateUserSettings) error
	Delete() error
	ChangeEmail(userID int, email string) error
	ChangeRole(userID int, role models.Role) error
	GetAll() ([]*models.User, error)
	GetUserSettings() (map[string]string, error)
	UpdateSettings(settings map[string]string) error
	HasSubscribedTo(ideaID int) (bool, error)
}

// Tenant contains read and write operations for tenants
type Tenant interface {
	Base
	Add(name string, subdomain string, status int) (*models.Tenant, error)
	First() (*models.Tenant, error)
	Activate(id int) error
	GetByDomain(domain string) (*models.Tenant, error)
	UpdateSettings(settings *models.UpdateTenantSettings) error
	UpdateAdvancedSettings(settings *models.UpdateTenantAdvancedSettings) error
	UpdatePrivacy(settings *models.UpdateTenantPrivacy) error
	IsSubdomainAvailable(subdomain string) (bool, error)
	IsCNAMEAvailable(cname string) (bool, error)
	SaveVerificationKey(key string, duration time.Duration, request models.NewEmailVerification) error
	FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error)
	SetKeyAsVerified(key string) error
	GetUpload(id int) (*models.Upload, error)
	SaveOAuthConfig(config *models.CreateEditOAuthConfig) error
	GetOAuthConfigByProvider(provider string) (*models.OAuthConfig, error)
	ListOAuthConfig() ([]*models.OAuthConfig, error)
}

// Tag contains read and write operations for tags
type Tag interface {
	Base
	Add(name, color string, isPublic bool) (*models.Tag, error)
	GetBySlug(slug string) (*models.Tag, error)
	Update(tag *models.Tag, name, color string, isPublic bool) (*models.Tag, error)
	Delete(tag *models.Tag) error
	GetAssigned(idea *models.Idea) ([]*models.Tag, error)
	AssignTag(tag *models.Tag, idea *models.Idea) error
	UnassignTag(tag *models.Tag, idea *models.Idea) error
	GetAll() ([]*models.Tag, error)
}

// Notification contains read and write operations for notifications
type Notification interface {
	Base
	Insert(user *models.User, title, link string, ideaID int) (*models.Notification, error)
	MarkAsRead(id int) error
	MarkAllAsRead() error
	TotalUnread() (int, error)
	GetActiveNotifications() ([]*models.Notification, error)
	GetNotification(id int) (*models.Notification, error)
}
