package postgres

import (
	"time"

	"database/sql"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/gosimple/slug"
)

// IdeaStorage contains read and write operations for ideas
type IdeaStorage struct {
	DB *dbx.Database
}

var (
	sqlSelectIdeasWhere = `SELECT i.id, 
								  i.number, 
								  i.title, 
								  i.slug, 
								  i.description, 
								  i.created_on,
								  i.supporters, 
								  i.status, 
								  u.id AS user_id, 
								  u.name AS user_name, 
								  u.email AS user_email
							FROM ideas i
							INNER JOIN users u
							ON u.id = i.user_id
							WHERE`
)

// GetAll returns all tenant ideas
func (s *IdeaStorage) GetAll(tenantID int) ([]*models.Idea, error) {
	var ideas []*models.Idea
	err := s.DB.Select(&ideas, sqlSelectIdeasWhere+" i.tenant_id = $1 ORDER BY i.created_on DESC", tenantID)
	if err != nil {
		return nil, err
	}

	return ideas, nil
}

// GetByID returns idea by given id
func (s *IdeaStorage) GetByID(tenantID, ideaID int) (*models.Idea, error) {
	idea := models.Idea{}

	err := s.DB.Get(&idea, sqlSelectIdeasWhere+" i.tenant_id = $1 AND i.id = $2 ORDER BY i.created_on DESC", tenantID, ideaID)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &idea, nil
}

// GetByNumber returns idea by tenant and number
func (s *IdeaStorage) GetByNumber(tenantID, number int) (*models.Idea, error) {
	idea := models.Idea{}

	err := s.DB.Get(&idea, sqlSelectIdeasWhere+" i.tenant_id = $1 AND i.number = $2 ORDER BY i.created_on DESC", tenantID, number)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &idea, nil
}

// GetCommentsByIdeaID returns all coments from given idea
func (s *IdeaStorage) GetCommentsByIdeaID(tenantID, ideaID int) ([]*models.Comment, error) {
	comments := []*models.Comment{}
	err := s.DB.Select(&comments,
		`SELECT c.id, 
				c.content, 
				c.created_on, 
				u.id AS user_id, 
				u.name AS user_email,
				u.email AS user_email
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

	return comments, nil
}

// Save a new idea in the database
func (s *IdeaStorage) Save(tenantID, userID int, title, description string) (*models.Idea, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	idea := new(models.Idea)
	idea.Title = title
	idea.Description = description

	row := tx.QueryRow(`INSERT INTO ideas (title, slug, number, description, tenant_id, user_id, created_on, supporters, status) 
						VALUES ($1, $2, (SELECT COALESCE(MAX(number), 0) + 1 FROM ideas i WHERE i.tenant_id = $4), $3, $4, $5, $6, 0, 0) 
						RETURNING id`, title, slug.Make(title), description, tenantID, userID, time.Now())
	if err = row.Scan(&idea.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	return idea, tx.Commit()
}

// AddComment places a new comment on an idea
func (s *IdeaStorage) AddComment(userID, ideaID int, content string) (int, error) {
	tx, err := s.DB.Begin()
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

// AddSupporter adds user to idea list of supporters
func (s *IdeaStorage) AddSupporter(userID, ideaID int) error {
	alreadySupported, err := s.DB.Exists("SELECT 1 FROM idea_supporters WHERE user_id = $1 AND idea_id = $2", userID, ideaID)
	if err != nil {
		return err
	}

	if alreadySupported {
		return nil
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`UPDATE ideas SET supporters = supporters + 1 WHERE id = $1`, ideaID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`INSERT INTO idea_supporters (user_id, idea_id, created_on) VALUES ($1, $2, $3)`, userID, ideaID, time.Now()); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// RemoveSupporter removes user from idea list of supporters
func (s *IdeaStorage) RemoveSupporter(userID, ideaID int) error {
	didSupport, err := s.DB.Exists("SELECT 1 FROM idea_supporters WHERE user_id = $1 AND idea_id = $2", userID, ideaID)
	if err != nil {
		return err
	}

	if !didSupport {
		return nil
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`UPDATE ideas SET supporters = supporters - 1 WHERE id = $1`, ideaID); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`DELETE FROM idea_supporters WHERE user_id = $1 AND idea_id = $2`, userID, ideaID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
