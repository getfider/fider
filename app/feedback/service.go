package feedback

// IdeaService contains read and write operations for ideas
type IdeaService interface {
	GetByID(tenantID, ideaID int) (*Idea, error)
	GetCommentsByIdeaID(tenantID, ideaID int) ([]*Comment, error)
	GetAll(tenantID int) ([]*Idea, error)
	Save(tenantID, userID int, title, description string) (*Idea, error)
	AddComment(userID, ideaID int, content string) (int, error)
}
