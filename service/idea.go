package service

import (
	"database/sql"

	"github.com/WeCanHearYou/wchy/model"
)

// IdeaService contains read and write operations for ideas
type IdeaService interface {
	GetAll(tenantID int64) ([]*model.Idea, error)
}

// PostgresIdeaService contains read and write operations for ideas
type PostgresIdeaService struct {
	DB *sql.DB
}

// GetAll returns all tenant ideas
func (svc PostgresIdeaService) GetAll(tenantID int64) ([]*model.Idea, error) {
	rows, err := svc.DB.Query("SELECT id, title, description FROM ideas WHERE tenant_id = $1", tenantID)
	if err != nil {
		return nil, err
	}

	var ideas []*model.Idea
	for rows.Next() {
		idea := &model.Idea{}
		rows.Scan(&idea.ID, &idea.Title, &idea.Description)
		ideas = append(ideas, idea)
	}

	return ideas, nil
}
