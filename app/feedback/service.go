package feedback

// IdeaService contains read and write operations for ideas
type IdeaService interface {
	GetByID(tenantID, ideaID int64) (*Idea, error)
	GetCommentsByIdeaID(tenantID, ideaID int64) ([]*Comment, error)
	GetAll(tenantID int64) ([]*Idea, error)
	Save(tenantID, userID int64, title, description string) (*Idea, error)
	AddComment(userID, ideaID int64, content string) (int64, error)
}
