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

// Post contains read and write operations for posts
type Post interface {
	Base
	GetByID(postID int) (*models.Post, error)
	GetBySlug(slug string) (*models.Post, error)
	GetByNumber(number int) (*models.Post, error)
	GetCommentsByPost(post *models.Post) ([]*models.Comment, error)
	Search(query, view, limit string, tags []string) ([]*models.Post, error)
	GetAll() ([]*models.Post, error)
	CountPerStatus() (map[models.PostStatus]int, error)
	Add(title, description string) (*models.Post, error)
	Update(post *models.Post, title, description string) (*models.Post, error)
	AddComment(post *models.Post, content string) (int, error)
	GetCommentByID(id int) (*models.Comment, error)
	UpdateComment(id int, content string) error
	DeleteComment(id int) error
	AddVote(post *models.Post, user *models.User) error
	RemoveVote(post *models.Post, user *models.User) error
	AddSubscriber(post *models.Post, user *models.User) error
	RemoveSubscriber(post *models.Post, user *models.User) error
	GetActiveSubscribers(number int, channel models.NotificationChannel, event models.NotificationEvent) ([]*models.User, error)
	SetResponse(post *models.Post, text string, status models.PostStatus) error
	MarkAsDuplicate(post *models.Post, original *models.Post) error
	IsReferenced(post *models.Post) (bool, error)
	VotedBy() ([]int, error)
	ListVotes(post *models.Post, limit int) ([]*models.Vote, error)
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
	HasSubscribedTo(postID int) (bool, error)
	GetByAPIKey(apiKey string) (*models.User, error)
	RegenerateAPIKey() (string, error)
	Block(userID int) error
	Unblock(userID int) error
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
	GetAssigned(post *models.Post) ([]*models.Tag, error)
	AssignTag(tag *models.Tag, post *models.Post) error
	UnassignTag(tag *models.Tag, post *models.Post) error
	GetAll() ([]*models.Tag, error)
}

// Notification contains read and write operations for notifications
type Notification interface {
	Base
	Insert(user *models.User, title, link string, postID int) (*models.Notification, error)
	MarkAsRead(id int) error
	MarkAllAsRead() error
	TotalUnread() (int, error)
	GetActiveNotifications() ([]*models.Notification, error)
	GetNotification(id int) (*models.Notification, error)
}

// Event contains read and write operations for Audit Events
type Event interface {
	Base
	Add(clientIP, name string) (*models.Event, error)
	GetByID(id int) (*models.Event, error)
}
