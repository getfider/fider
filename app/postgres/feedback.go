package postgres

import (
	"database/sql"

	"time"

	"github.com/WeCanHearYou/wechy/app/feedback"
)

// IdeaService contains read and write operations for ideas
type IdeaService struct {
	DB *sql.DB
}

// GetAll returns all tenant ideas
func (svc IdeaService) GetAll(tenantID int64) ([]*feedback.Idea, error) {
	rows, err := svc.DB.Query("SELECT id, title, description, created_on FROM ideas WHERE tenant_id = $1", tenantID)
	if err != nil {
		return nil, err
	}

	var ideas []*feedback.Idea
	for rows.Next() {
		idea := &feedback.Idea{}
		rows.Scan(&idea.ID, &idea.Title, &idea.Description, &idea.CreatedOn)
		ideas = append(ideas, idea)
	}

	return ideas, nil
}

// Save a new idea in the database
func (svc IdeaService) Save(tenantID, userID int64, title, description string) (*feedback.Idea, error) {
	tx, err := svc.DB.Begin()
	if err != nil {
		return nil, err
	}

	idea := new(feedback.Idea)
	idea.Title = title
	idea.Description = description

	if err = tx.QueryRow("INSERT INTO ideas (title, description, tenant_id, user_id, created_on) VALUES ($1, $2, $3, $4, $5) RETURNING id", title, description, tenantID, userID, time.Now()).Scan(&idea.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	return idea, tx.Commit()
}
