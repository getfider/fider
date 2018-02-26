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
	GetCommentsByIdea(number int) ([]*models.Comment, error)
	Search(query, filter string, tags []string) ([]*models.Idea, error)
	GetAll() ([]*models.Idea, error)
	CountPerStatus() (map[int]int, error)
	Add(title, description string, userID int) (*models.Idea, error)
	Update(number int, title, description string) (*models.Idea, error)
	AddComment(number int, content string, userID int) (int, error)
	AddSupporter(number, userID int) error
	RemoveSupporter(number, userID int) error
	AddSubscriber(number, userID int) error
	RemoveSubscriber(number, userID int) error
	GetActiveSubscribers(number int, channel models.NotificationChannel, event models.NotificationEvent) ([]*models.User, error)
	SetResponse(number int, text string, userID, status int) error
	MarkAsDuplicate(number, originalNumber, userID int) error
	SupportedBy(userID int) ([]int, error)
}

// User is used for user operations
type User interface {
	Base
	GetByID(userID int) (*models.User, error)
	GetByEmail(tenantID int, email string) (*models.User, error)
	GetByProvider(tenantID int, provider string, uid string) (*models.User, error)
	Register(user *models.User) error
	RegisterProvider(userID int, provider *models.UserProvider) error
	Update(userID int, settings *models.UpdateUserSettings) error
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
	IsSubdomainAvailable(subdomain string) (bool, error)
	IsCNAMEAvailable(cname string) (bool, error)
	SaveVerificationKey(key string, duration time.Duration, request models.NewEmailVerification) error
	FindVerificationByKey(kind models.EmailVerificationKind, key string) (*models.EmailVerification, error)
	SetKeyAsVerified(key string) error
}

// Tag contains read and write operations for tags
type Tag interface {
	Base
	Add(name, color string, isPublic bool) (*models.Tag, error)
	GetBySlug(slug string) (*models.Tag, error)
	Update(tagID int, name, color string, isPublic bool) (*models.Tag, error)
	Delete(tagID int) error
	GetAssigned(ideaID int) ([]*models.Tag, error)
	AssignTag(tagID, ideaID, userID int) error
	UnassignTag(tagID, ideaID int) error
	GetAll() ([]*models.Tag, error)
}

// Notification contains read and write operations for notifications
type Notification interface {
	Base
	Insert(user *models.User, title, link string, ideaID, authorID int) (*models.Notification, error)
	MarkAsRead(id int) error
	MarkAllAsRead() error
	TotalUnread() (int, error)
	GetActiveNotifications() ([]*models.Notification, error)
	GetNotification(id int) (*models.Notification, error)
}
