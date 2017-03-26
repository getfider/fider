package postgres

import (
	"database/sql"

	"time"

	"github.com/WeCanHearYou/wechy/app"
	"github.com/WeCanHearYou/wechy/app/feedback"
)

// IdeaService contains read and write operations for ideas
type IdeaService struct {
	DB *sql.DB
}

// GetAll returns all tenant ideas
func (svc IdeaService) GetAll(tenantID int) ([]*feedback.Idea, error) {
	rows, err := svc.DB.Query(`SELECT i.id, i.title, i.description, i.created_on, u.id, u.name, u.email
														 FROM ideas i
														 INNER JOIN users u
														 ON u.id = i.user_id
														 WHERE i.tenant_id = $1
														 ORDER BY i.created_on DESC`, tenantID)
	if err != nil {
		return nil, err
	}

	var ideas []*feedback.Idea
	for rows.Next() {
		idea := &feedback.Idea{}
		rows.Scan(&idea.ID, &idea.Title, &idea.Description, &idea.CreatedOn, &idea.User.ID, &idea.User.Name, &idea.User.Email)
		ideas = append(ideas, idea)
	}

	return ideas, nil
}

// GetByID returns idea by given id
func (svc IdeaService) GetByID(tenantID, ideaID int) (*feedback.Idea, error) {
	rows, err := svc.DB.Query(`SELECT i.id, i.title, i.description, i.created_on, u.id, u.name, u.email
														 FROM ideas i
														 INNER JOIN users u
														 ON u.id = i.user_id
														 WHERE i.tenant_id = $1
														 AND i.id = $2
														 ORDER BY i.created_on DESC`, tenantID, ideaID)
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		idea := &feedback.Idea{}
		rows.Scan(&idea.ID, &idea.Title, &idea.Description, &idea.CreatedOn, &idea.User.ID, &idea.User.Name, &idea.User.Email)
		return idea, nil
	}
	return nil, app.ErrNotFound
}

// GetCommentsByIdeaID returns all coments from given idea
func (svc IdeaService) GetCommentsByIdeaID(tenantID, ideaID int) ([]*feedback.Comment, error) {
	rows, err := svc.DB.Query(`SELECT c.id, c.content, c.created_on, u.id, u.name, u.email
														 FROM comments c
														 INNER JOIN ideas i
														 ON i.id = c.idea_id
														 INNER JOIN users u
														 ON u.id = c.user_id
														 WHERE i.id = $1
														 AND i.tenant_id = $2
														 ORDER BY c.created_on DESC`, ideaID, tenantID)
	if err != nil {
		return nil, err
	}

	var comments []*feedback.Comment
	for rows.Next() {
		c := &feedback.Comment{}
		rows.Scan(&c.ID, &c.Content, &c.CreatedOn, &c.User.ID, &c.User.Name, &c.User.Email)
		comments = append(comments, c)
	}

	return comments, nil
}

// Save a new idea in the database
func (svc IdeaService) Save(tenantID, userID int, title, description string) (*feedback.Idea, error) {
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

// AddComment places a new comment on an idea
func (svc IdeaService) AddComment(userID, ideaID int, content string) (int, error) {
	tx, err := svc.DB.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	if err = tx.QueryRow("INSERT INTO comments (idea_id, content, user_id, created_on) VALUES ($1, $2, $3, $4) RETURNING id", ideaID, content, userID, time.Now()).Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
