package feedback

// IdeaService contains read and write operations for ideas
type IdeaService interface {
	GetAll(tenantID int64) ([]*Idea, error)
}
