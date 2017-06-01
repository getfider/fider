package postgres

import (
	"time"

	"database/sql"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/gosimple/slug"
)

type dbIdea struct {
	ID              int            `db:"id"`
	Number          int            `db:"number"`
	Title           string         `db:"title"`
	Slug            string         `db:"slug"`
	Description     string         `db:"description"`
	CreatedOn       time.Time      `db:"created_on"`
	User            *dbUser        `db:"user"`
	TotalSupporters int            `db:"supporters"`
	Status          int            `db:"status"`
	Response        sql.NullString `db:"response"`
	RespondedOn     dbx.NullTime   `db:"response_date"`
	ResponseUser    *dbUser        `db:"response_user"`
}

func (i *dbIdea) toModel() *models.Idea {
	idea := &models.Idea{
		ID:              i.ID,
		Number:          i.Number,
		Title:           i.Title,
		Slug:            i.Slug,
		Description:     i.Description,
		CreatedOn:       i.CreatedOn,
		TotalSupporters: i.TotalSupporters,
		Status:          i.Status,
		User:            i.User.toModel(),
	}

	if i.Response.Valid {
		idea.Response = &models.IdeaResponse{
			Text:      i.Response.String,
			CreatedOn: i.RespondedOn.Time,
			User:      i.ResponseUser.toModel(),
		}
	}
	return idea
}

type dbComment struct {
	ID        int       `db:"id"`
	Content   string    `db:"content"`
	CreatedOn time.Time `db:"created_on"`
	User      *dbUser   `db:"user"`
}

func (c *dbComment) toModel() *models.Comment {
	return &models.Comment{
		ID:        c.ID,
		Content:   c.Content,
		CreatedOn: c.CreatedOn,
		User:      c.User.toModel(),
	}
}

// IdeaStorage contains read and write operations for ideas
type IdeaStorage struct {
	Trx    *dbx.Trx
	Tenant *models.Tenant
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
																u.email AS user_email,
																i.response,
																i.response_date,
																r.id AS response_user_id, 
																r.name AS response_user_name, 
																r.email AS response_user_email
													FROM ideas i
													INNER JOIN users u
													ON u.id = i.user_id
													LEFT JOIN users r
													ON r.id = i.response_user_id
													WHERE`
)

// GetAll returns all tenant ideas
func (s *IdeaStorage) GetAll() ([]*models.Idea, error) {
	var ideas []*dbIdea
	err := s.Trx.Select(&ideas, sqlSelectIdeasWhere+" i.tenant_id = $1 ORDER BY i.created_on DESC", s.Tenant.ID)
	if err != nil {
		return nil, err
	}

	var result = make([]*models.Idea, len(ideas))
	for i, idea := range ideas {
		result[i] = idea.toModel()
	}
	return result, nil
}

// GetByID returns idea by given id
func (s *IdeaStorage) GetByID(ideaID int) (*models.Idea, error) {
	idea := dbIdea{}

	err := s.Trx.Get(&idea, sqlSelectIdeasWhere+" i.tenant_id = $1 AND i.id = $2 ORDER BY i.created_on DESC", s.Tenant.ID, ideaID)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return idea.toModel(), nil
}

// GetByNumber returns idea by tenant and number
func (s *IdeaStorage) GetByNumber(tenantID, number int) (*models.Idea, error) {
	idea := dbIdea{}

	err := s.Trx.Get(&idea, sqlSelectIdeasWhere+" i.tenant_id = $1 AND i.number = $2 ORDER BY i.created_on DESC", tenantID, number)
	if err == sql.ErrNoRows {
		return nil, app.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return idea.toModel(), nil
}

// GetCommentsByIdeaID returns all coments from given idea
func (s *IdeaStorage) GetCommentsByIdeaID(tenantID, ideaID int) ([]*models.Comment, error) {
	comments := []*dbComment{}
	err := s.Trx.Select(&comments,
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

	var result = make([]*models.Comment, len(comments))
	for i, comment := range comments {
		result[i] = comment.toModel()
	}
	return result, nil
}

// Save a new idea in the database
func (s *IdeaStorage) Save(tenantID, userID int, title, description string) (*models.Idea, error) {
	idea := new(models.Idea)
	idea.Title = title
	idea.Description = description

	row := s.Trx.QueryRow(`INSERT INTO ideas (title, slug, number, description, tenant_id, user_id, created_on, supporters, status) 
						VALUES ($1, $2, (SELECT COALESCE(MAX(number), 0) + 1 FROM ideas i WHERE i.tenant_id = $4), $3, $4, $5, $6, 0, 0) 
						RETURNING id`, title, slug.Make(title), description, tenantID, userID, time.Now())
	if err := row.Scan(&idea.ID); err != nil {
		return nil, err
	}

	return idea, nil
}

// AddComment places a new comment on an idea
func (s *IdeaStorage) AddComment(userID, ideaID int, content string) (int, error) {
	var id int
	if err := s.Trx.QueryRow("INSERT INTO comments (idea_id, content, user_id, created_on) VALUES ($1, $2, $3, $4) RETURNING id", ideaID, content, userID, time.Now()).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// AddSupporter adds user to idea list of supporters
func (s *IdeaStorage) AddSupporter(tenantID, userID, ideaID int) error {
	idea, err := s.GetByID(ideaID)
	if err != nil {
		return err
	}
	if idea.Status >= models.IdeaCompleted {
		return nil
	}

	alreadySupported, err := s.Trx.Exists("SELECT 1 FROM idea_supporters WHERE user_id = $1 AND idea_id = $2", userID, ideaID)
	if err != nil {
		return err
	}

	if alreadySupported {
		return nil
	}

	if err := s.Trx.Execute(`UPDATE ideas SET supporters = supporters + 1 WHERE id = $1`, ideaID); err != nil {
		return err
	}

	if err := s.Trx.Execute(`INSERT INTO idea_supporters (user_id, idea_id, created_on) VALUES ($1, $2, $3)`, userID, ideaID, time.Now()); err != nil {
		return err
	}

	return nil
}

// RemoveSupporter removes user from idea list of supporters
func (s *IdeaStorage) RemoveSupporter(tenantID, userID, ideaID int) error {
	idea, err := s.GetByID(ideaID)
	if err != nil {
		return err
	}

	if idea.Status >= models.IdeaCompleted {
		return nil
	}

	didSupport, err := s.Trx.Exists("SELECT 1 FROM idea_supporters WHERE user_id = $1 AND idea_id = $2", userID, ideaID)
	if err != nil {
		return err
	}

	if !didSupport {
		return nil
	}

	if err := s.Trx.Execute(`UPDATE ideas SET supporters = supporters - 1 WHERE id = $1`, ideaID); err != nil {
		return err
	}

	if err := s.Trx.Execute(`DELETE FROM idea_supporters WHERE user_id = $1 AND idea_id = $2`, userID, ideaID); err != nil {
		return err
	}

	return nil
}

// SetResponse changes current idea response
func (s *IdeaStorage) SetResponse(tenantID, ideaID int, text string, userID, status int) error {
	if err := s.Trx.Execute(`UPDATE ideas SET response = $3, response_date = $4, response_user_id = $5, status = $6 WHERE id = $1 and tenant_id = $2`, ideaID, tenantID, text, time.Now(), userID, status); err != nil {
		return err
	}
	return nil
}
