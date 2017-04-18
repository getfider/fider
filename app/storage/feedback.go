package storage

import "github.com/WeCanHearYou/wechy/app/models"

// Idea contains read and write operations for ideas
type Idea interface {
	GetByID(tenantID, ideaID int) (*models.Idea, error)
	GetByNumber(tenantID, number int) (*models.Idea, error)
	GetCommentsByIdeaID(tenantID, ideaID int) ([]*models.Comment, error)
	GetAll(tenantID int) ([]*models.Idea, error)
	Save(tenantID, userID int, title, description string) (*models.Idea, error)
	AddComment(userID, ideaID int, content string) (int, error)
	AddSupport(userID, ideaID int) error
	RemoveSupport(userID, ideaID int) error
}
