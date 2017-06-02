package storage

import "github.com/getfider/fider/app/models"

// Idea contains read and write operations for ideas
type Idea interface {
	GetByID(ideaID int) (*models.Idea, error)
	GetByNumber(number int) (*models.Idea, error)
	GetCommentsByIdea(number int) ([]*models.Comment, error)
	GetAll() ([]*models.Idea, error)
	Add(title, description string, userID int) (*models.Idea, error)
	AddComment(number int, content string, userID int) (int, error)
	AddSupporter(number, userID int) error
	RemoveSupporter(number, userID int) error
	SetResponse(number int, text string, userID, status int) error
	SupportedBy(userID int) ([]int, error)
}
