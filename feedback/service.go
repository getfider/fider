package feedback

// IdeaService contains read and write operations for ideas
type IdeaService interface {
	GetAll(tenantID int64) ([]*Idea, error)
	Save(tenantID, userID int64, title, description string) (*Idea, error)
}
