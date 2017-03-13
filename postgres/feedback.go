package postgres

import (
	"database/sql"

	"github.com/WeCanHearYou/wechy/feedback"
)

// IdeaService contains read and write operations for ideas
type IdeaService struct {
	DB *sql.DB
}

// GetAll returns all tenant ideas
func (svc IdeaService) GetAll(tenantID int64) ([]*feedback.Idea, error) {
	rows, err := svc.DB.Query("SELECT id, title, description FROM ideas WHERE tenant_id = $1", tenantID)
	if err != nil {
		return nil, err
	}

	var ideas []*feedback.Idea
	for rows.Next() {
		idea := &feedback.Idea{}
		rows.Scan(&idea.ID, &idea.Title, &idea.Description)
		ideas = append(ideas, idea)
	}

	return ideas, nil
}
