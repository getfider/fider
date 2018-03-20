package postgres

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"database/sql"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

type dbIdea struct {
	ID               int            `db:"id"`
	Number           int            `db:"number"`
	Title            string         `db:"title"`
	Slug             string         `db:"slug"`
	Description      string         `db:"description"`
	CreatedOn        time.Time      `db:"created_on"`
	User             *dbUser        `db:"user"`
	ViewerSupported  bool           `db:"viewer_supported"`
	TotalSupporters  int            `db:"supporters"`
	TotalComments    int            `db:"comments"`
	RecentSupporters int            `db:"recent_supporters"`
	RecentComments   int            `db:"recent_comments"`
	Status           int            `db:"status"`
	Response         sql.NullString `db:"response"`
	RespondedOn      dbx.NullTime   `db:"response_date"`
	ResponseUser     *dbUser        `db:"response_user"`
	OriginalNumber   sql.NullInt64  `db:"original_number"`
	OriginalTitle    sql.NullString `db:"original_title"`
	OriginalSlug     sql.NullString `db:"original_slug"`
	OriginalStatus   sql.NullInt64  `db:"original_status"`
	Tags             []string       `db:"tags"`
}

func (i *dbIdea) toModel() *models.Idea {
	idea := &models.Idea{
		ID:              i.ID,
		Number:          i.Number,
		Title:           i.Title,
		Slug:            i.Slug,
		Description:     i.Description,
		CreatedOn:       i.CreatedOn,
		ViewerSupported: i.ViewerSupported,
		TotalSupporters: i.TotalSupporters,
		TotalComments:   i.TotalComments,
		Status:          i.Status,
		User:            i.User.toModel(),
		Tags:            i.Tags,
	}

	if i.Response.Valid {
		idea.Response = &models.IdeaResponse{
			Text:        i.Response.String,
			RespondedOn: i.RespondedOn.Time,
			User:        i.ResponseUser.toModel(),
		}
		if idea.Status == models.IdeaDuplicate && i.OriginalNumber.Valid {
			idea.Response.Original = &models.OriginalIdea{
				Number: int(i.OriginalNumber.Int64),
				Slug:   i.OriginalSlug.String,
				Title:  i.OriginalTitle.String,
				Status: int(i.OriginalStatus.Int64),
			}
		}
	}
	return idea
}

type dbComment struct {
	ID        int          `db:"id"`
	Content   string       `db:"content"`
	CreatedOn time.Time    `db:"created_on"`
	User      *dbUser      `db:"user"`
	EditedOn  dbx.NullTime `db:"edited_on"`
	EditedBy  *dbUser      `db:"edited_by"`
}

func (c *dbComment) toModel() *models.Comment {
	comment := &models.Comment{
		ID:        c.ID,
		Content:   c.Content,
		CreatedOn: c.CreatedOn,
		User:      c.User.toModel(),
	}
	if c.EditedOn.Valid {
		comment.EditedBy = c.EditedBy.toModel()
		comment.EditedOn = &c.EditedOn.Time
	}
	return comment
}

type dbStatusCount struct {
	Status int `db:"status"`
	Count  int `db:"count"`
}

// IdeaStorage contains read and write operations for ideas
type IdeaStorage struct {
	trx    *dbx.Trx
	tenant *models.Tenant
	user   *models.User
}

// NewIdeaStorage creates a new IdeaStorage
func NewIdeaStorage(trx *dbx.Trx) *IdeaStorage {
	return &IdeaStorage{
		trx: trx,
	}
}

// SetCurrentTenant to current context
func (s *IdeaStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// SetCurrentUser to current context
func (s *IdeaStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

var (
	sqlSelectIdeasWhere = `	WITH 
													agg_comments AS (
															SELECT 
																	idea_id, 
																	COUNT(CASE WHEN comments.created_on > CURRENT_DATE - INTERVAL '30 days' THEN 1 END) as recent,
																	COUNT(*) as all
															FROM comments 
															INNER JOIN ideas
															ON ideas.id = comments.idea_id
															AND ideas.tenant_id = comments.tenant_id
															WHERE ideas.tenant_id = $1
															GROUP BY idea_id
													),
													agg_supporters AS (
															SELECT 
																	idea_id, 
																	COUNT(*) as recent
															FROM idea_supporters 
															INNER JOIN ideas
															ON ideas.id = idea_supporters.idea_id
															AND ideas.tenant_id = idea_supporters.tenant_id
															WHERE ideas.tenant_id = $1
															AND idea_supporters.created_on > CURRENT_DATE - INTERVAL '30 days' 
															GROUP BY idea_id
													)
													SELECT i.id, 
																i.number, 
																i.title, 
																i.slug, 
																i.description, 
																i.created_on,
																i.supporters,
																COALESCE(agg_c.all, 0) as comments,
																COALESCE(agg_s.recent, 0) AS recent_supporters,
																COALESCE(agg_c.recent, 0) AS recent_comments,																
																i.status, 
																u.id AS user_id, 
																u.name AS user_name, 
																u.email AS user_email,
																u.role AS user_role,
																i.response,
																i.response_date,
																r.id AS response_user_id, 
																r.name AS response_user_name, 
																r.email AS response_user_email, 
																r.role AS response_user_role,
																d.number AS original_number,
																d.title AS original_title,
																d.slug AS original_slug,
																d.status AS original_status,
																array_remove(array_agg(t.slug), NULL) AS tags,
																COALESCE(%s, false) AS viewer_supported
													FROM ideas i
													INNER JOIN users u
													ON u.id = i.user_id
													LEFT JOIN users r
													ON r.id = i.response_user_id
													LEFT JOIN idea_tags it
													ON it.idea_id = i.id
													LEFT JOIN ideas d
													ON d.id = i.original_id
													LEFT JOIN tags t
													ON t.id = it.tag_id
													%s
													LEFT JOIN agg_comments agg_c
													ON agg_c.idea_id = i.id
													LEFT JOIN agg_supporters agg_s
													ON agg_s.idea_id = i.id
													WHERE i.status != ` + strconv.Itoa(models.IdeaDeleted) + ` AND %s
													GROUP BY i.id, u.id, r.id, d.id, agg_c.all, agg_c.recent, agg_s.recent`
)

func (s *IdeaStorage) getIdeaQuery(filter string) string {
	viewerSupportedSubQuery := "null"
	if s.user != nil {
		viewerSupportedSubQuery = fmt.Sprintf("(SELECT true FROM idea_supporters WHERE idea_id = i.id AND user_id = %d)", s.user.ID)
	}
	tagCondition := `AND t.is_public = true`
	if s.user != nil && s.user.IsCollaborator() {
		tagCondition = ``
	}
	return fmt.Sprintf(sqlSelectIdeasWhere, viewerSupportedSubQuery, tagCondition, filter)
}

func (s *IdeaStorage) getSingle(query string, args ...interface{}) (*models.Idea, error) {
	idea := dbIdea{}

	if err := s.trx.Get(&idea, query, args...); err != nil {
		return nil, err
	}

	return idea.toModel(), nil
}

// GetByID returns idea by given id
func (s *IdeaStorage) GetByID(ideaID int) (*models.Idea, error) {
	return s.getSingle(s.getIdeaQuery("i.tenant_id = $1 AND i.id = $2"), s.tenant.ID, ideaID)
}

// GetBySlug returns idea by tenant and slug
func (s *IdeaStorage) GetBySlug(slug string) (*models.Idea, error) {
	return s.getSingle(s.getIdeaQuery("i.tenant_id = $1 AND i.slug = $2"), s.tenant.ID, slug)
}

// GetByNumber returns idea by tenant and number
func (s *IdeaStorage) GetByNumber(number int) (*models.Idea, error) {
	return s.getSingle(s.getIdeaQuery("i.tenant_id = $1 AND i.number = $2"), s.tenant.ID, number)
}

// GetAll returns all tenant ideas
func (s *IdeaStorage) GetAll() ([]*models.Idea, error) {
	return s.Search("", "all", []string{})
}

// CountPerStatus returns total number of ideas per status
func (s *IdeaStorage) CountPerStatus() (map[int]int, error) {
	stats := []*dbStatusCount{}
	err := s.trx.Select(&stats, "SELECT status, COUNT(*) AS count FROM ideas WHERE tenant_id = $1 GROUP BY status", s.tenant.ID)
	if err != nil {
		return nil, err
	}
	result := make(map[int]int, len(stats))
	for _, v := range stats {
		result[v.Status] = v.Count
	}
	return result, nil
}

// Search existing ideas based on input
func (s *IdeaStorage) Search(query, filter string, tags []string) ([]*models.Idea, error) {
	innerQuery := s.getIdeaQuery("i.tenant_id = $1 AND i.status = ANY($2)")

	var (
		ideas []*dbIdea
		err   error
	)
	if query != "" {
		scoreField := "ts_rank(setweight(to_tsvector(title), 'A') || setweight(to_tsvector(description), 'B'), to_tsquery('english', $3)) + similarity(title, $4) + similarity(description, $4)"
		sql := fmt.Sprintf(`
			SELECT * FROM (%s) AS q 
			WHERE %s > 0.1 
			ORDER BY %s DESC
		`, innerQuery, scoreField, scoreField)
		err = s.trx.Select(&ideas, sql, s.tenant.ID, pq.Array([]int{
			models.IdeaOpen,
			models.IdeaStarted,
			models.IdeaPlanned,
			models.IdeaCompleted,
			models.IdeaDeclined,
		}), ToTSQuery(query), query)
	} else {
		statuses, sort := getFilterData(filter)
		sql := fmt.Sprintf(`
			SELECT * FROM (%s) AS q 
			WHERE tags @> $3
			ORDER BY %s DESC
		`, innerQuery, sort)
		err = s.trx.Select(&ideas, sql, s.tenant.ID, pq.Array(statuses), pq.Array(tags))
	}

	if err != nil {
		return nil, err
	}

	var result = make([]*models.Idea, len(ideas))
	for i, idea := range ideas {
		result[i] = idea.toModel()
	}
	return result, nil
}

// GetCommentsByIdea returns all comments from given idea
func (s *IdeaStorage) GetCommentsByIdea(number int) ([]*models.Comment, error) {
	comments := []*dbComment{}
	err := s.trx.Select(&comments,
		`SELECT c.id, 
				c.content, 
				c.created_on, 
				c.edited_on, 
				u.id AS user_id, 
				u.name AS user_name,
				u.email AS user_email,
				u.role AS user_role, 
				e.id AS edited_by_id, 
				e.name AS edited_by_name,
				e.email AS edited_by_email,
				e.role AS edited_by_role
		FROM comments c
		INNER JOIN ideas i
		ON i.id = c.idea_id
		AND i.tenant_id = c.tenant_id
		INNER JOIN users u
		ON u.id = c.user_id
		AND u.tenant_id = c.tenant_id
		LEFT JOIN users e
		ON e.id = c.edited_by_id
		AND e.tenant_id = c.tenant_id
		WHERE i.number = $1
		AND i.tenant_id = $2
		ORDER BY c.created_on ASC`, number, s.tenant.ID)
	if err != nil {
		return nil, err
	}

	var result = make([]*models.Comment, len(comments))
	for i, comment := range comments {
		result[i] = comment.toModel()
	}
	return result, nil
}

// Update given idea
func (s *IdeaStorage) Update(number int, title, description string) (*models.Idea, error) {
	err := s.trx.Execute(`UPDATE ideas SET title = $1, slug = $2, description = $3 
												WHERE number = $4 AND tenant_id = $5`, title, slug.Make(title), description, number, s.tenant.ID)
	if err != nil {
		return nil, err
	}

	return s.GetByNumber(number)
}

// Add a new idea in the database
func (s *IdeaStorage) Add(title, description string, userID int) (*models.Idea, error) {
	var id int
	err := s.trx.Get(&id,
		`INSERT INTO ideas (title, slug, number, description, tenant_id, user_id, created_on, supporters, status) 
		 VALUES ($1, $2, (SELECT COALESCE(MAX(number), 0) + 1 FROM ideas i WHERE i.tenant_id = $4), $3, $4, $5, $6, 0, 0) 
		 RETURNING id`, title, slug.Make(title), description, s.tenant.ID, userID, time.Now())
	if err != nil {
		return nil, err
	}

	idea, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.internalAddSubscriber(idea.Number, userID, false); err != nil {
		return nil, err
	}

	return idea, nil
}

// AddComment places a new comment on an idea
func (s *IdeaStorage) AddComment(number int, content string, userID int) (int, error) {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return 0, err
	}

	var id int
	if err := s.trx.Get(&id,
		"INSERT INTO comments (tenant_id, idea_id, content, user_id, created_on) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		s.tenant.ID, idea.ID, content, userID, time.Now()); err != nil {
		return 0, err
	}

	if err := s.internalAddSubscriber(number, userID, false); err != nil {
		return 0, err
	}

	return id, nil
}

// GetCommentByID returns a comment by given ID
func (s *IdeaStorage) GetCommentByID(id int) (*models.Comment, error) {
	comment := dbComment{}
	err := s.trx.Get(&comment,
		`SELECT c.id, 
						c.content, 
						c.created_on, 
						c.edited_on, 
						u.id AS user_id, 
						u.name AS user_name,
						u.email AS user_email,
						u.role AS user_role, 
						e.id AS edited_by_id, 
						e.name AS edited_by_name,
						e.email AS edited_by_email,
						e.role AS edited_by_role
		FROM comments c
		INNER JOIN users u
		ON u.id = c.user_id
		AND u.tenant_id = c.tenant_id
		LEFT JOIN users e
		ON e.id = c.edited_by_id
		AND e.tenant_id = c.tenant_id
		WHERE c.id = $1
		AND c.tenant_id = $2`, id, s.tenant.ID)

	if err != nil {
		return nil, err
	}

	return comment.toModel(), nil
}

// UpdateComment with given ID and content
func (s *IdeaStorage) UpdateComment(id int, content string) error {
	return s.trx.Execute(`UPDATE comments SET content = $1, edited_on = $2, edited_by_id = $3 
												WHERE id = $4 AND tenant_id = $5`, content, time.Now(), s.user.ID, id, s.tenant.ID)
}

// AddSupporter adds user to idea list of supporters
func (s *IdeaStorage) AddSupporter(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}

	if !idea.CanBeSupported() {
		return nil
	}

	alreadySupported, err := s.trx.Exists("SELECT 1 FROM idea_supporters WHERE user_id = $1 AND idea_id = $2 AND tenant_id = $3", userID, idea.ID, s.tenant.ID)
	if err != nil {
		return err
	}

	if alreadySupported {
		return nil
	}

	if err := s.trx.Execute(`UPDATE ideas SET supporters = supporters + 1 WHERE id = $1 AND tenant_id = $2`, idea.ID, s.tenant.ID); err != nil {
		return err
	}

	if err := s.trx.Execute(
		`INSERT INTO idea_supporters (tenant_id, user_id, idea_id, created_on) VALUES ($1, $2, $3, $4)  ON CONFLICT DO NOTHING`,
		s.tenant.ID, userID, idea.ID, time.Now()); err != nil {
		return err
	}

	return s.internalAddSubscriber(number, userID, false)
}

// RemoveSupporter removes user from idea list of supporters
func (s *IdeaStorage) RemoveSupporter(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}

	if !idea.CanBeSupported() {
		return nil
	}

	didSupport, err := s.trx.Exists("SELECT 1 FROM idea_supporters WHERE user_id = $1 AND idea_id = $2 AND tenant_id = $3", userID, idea.ID, s.tenant.ID)
	if err != nil {
		return err
	}

	if !didSupport {
		return nil
	}

	if err := s.trx.Execute(`UPDATE ideas SET supporters = supporters - 1 WHERE id = $1 AND tenant_id = $2`, idea.ID, s.tenant.ID); err != nil {
		return err
	}

	return s.trx.Execute(`DELETE FROM idea_supporters WHERE user_id = $1 AND idea_id = $2 AND tenant_id = $3`, userID, idea.ID, s.tenant.ID)
}

// AddSubscriber adds user to the idea list of subscribers
func (s *IdeaStorage) AddSubscriber(number, userID int) error {
	return s.internalAddSubscriber(number, userID, true)
}

func (s *IdeaStorage) internalAddSubscriber(number, userID int, force bool) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}

	conflict := " DO NOTHING"
	if force {
		conflict = "(user_id, idea_id) DO UPDATE SET status = $5, updated_on = $4"
	}

	return s.trx.Execute(fmt.Sprintf(`
	INSERT INTO idea_subscribers (tenant_id, user_id, idea_id, created_on, updated_on, status)
	VALUES ($1, $2, $3, $4, $4, $5)  ON CONFLICT %s`, conflict),
		s.tenant.ID, userID, idea.ID, time.Now(), models.SubscriberActive,
	)
}

// RemoveSubscriber removes user from idea list of subscribers
func (s *IdeaStorage) RemoveSubscriber(number, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}

	return s.trx.Execute(`
		INSERT INTO idea_subscribers (tenant_id, user_id, idea_id, created_on, updated_on, status)
		VALUES ($1, $2, $3, $4, $4, 0) ON CONFLICT (user_id, idea_id)
		DO UPDATE SET status = 0, updated_on = $4`,
		s.tenant.ID, userID, idea.ID, time.Now(),
	)
}

// GetActiveSubscribers based on input and settings
func (s *IdeaStorage) GetActiveSubscribers(number int, channel models.NotificationChannel, event models.NotificationEvent) ([]*models.User, error) {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return make([]*models.User, 0), err
	}

	var users []*dbUser

	if len(event.RequiresSubscripionUserRoles) == 0 {
		err = s.trx.Select(&users, `
			SELECT DISTINCT u.id, u.name, u.email, u.tenant_id, u.role
			FROM users u
			LEFT JOIN user_settings set
			ON set.user_id = u.id
			AND set.tenant_id = u.tenant_id
			AND set.key = $1
			WHERE u.tenant_id = $2
			AND (
				(set.value IS NULL AND u.role = ANY($3))
				OR CAST(set.value AS integer) & $4 > 0
			)`,
			event.UserSettingsKeyName,
			s.tenant.ID,
			pq.Array(event.DefaultEnabledUserRoles),
			channel,
		)
	} else {
		err = s.trx.Select(&users, `
			SELECT DISTINCT u.id, u.name, u.email, u.tenant_id, u.role
			FROM users u
			LEFT JOIN idea_subscribers sub
			ON sub.user_id = u.id
			AND sub.idea_id = $1
			AND sub.tenant_id = u.tenant_id
			LEFT JOIN user_settings set
			ON set.user_id = u.id
			AND set.key = $3
			AND set.tenant_id = u.tenant_id
			WHERE u.tenant_id = $4
			AND ( sub.status = $2 OR (sub.status IS NULL AND NOT u.role = ANY($7)) )
			AND (
				(set.value IS NULL AND u.role = ANY($5))
				OR CAST(set.value AS integer) & $6 > 0
			)`,
			idea.ID,
			models.SubscriberActive,
			event.UserSettingsKeyName,
			s.tenant.ID,
			pq.Array(event.DefaultEnabledUserRoles),
			channel,
			pq.Array(event.RequiresSubscripionUserRoles),
		)
	}

	if err != nil {
		return nil, err
	}

	var result = make([]*models.User, len(users))
	for i, user := range users {
		result[i] = user.toModel()
	}
	return result, nil
}

// SetResponse changes current idea response
func (s *IdeaStorage) SetResponse(number int, text string, userID, status int) error {
	if status == models.IdeaDuplicate {
		return errors.New("Use MarkAsDuplicate to change an idea status to Duplicate")
	}

	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}

	respondedOn := time.Now()
	if idea.Status == status && idea.Response != nil {
		respondedOn = idea.Response.RespondedOn
	}
	return s.trx.Execute(`
	UPDATE ideas 
	SET response = $3, original_id = NULL, response_date = $4, response_user_id = $5, status = $6 
	WHERE id = $1 and tenant_id = $2
	`, idea.ID, s.tenant.ID, text, respondedOn, userID, status)
}

// MarkAsDuplicate set idea as a duplicate of another idea
func (s *IdeaStorage) MarkAsDuplicate(number, originalNumber, userID int) error {
	idea, err := s.GetByNumber(number)
	if err != nil {
		return err
	}
	original, err := s.GetByNumber(originalNumber)
	if err != nil {
		return err
	}

	respondedOn := time.Now()
	if idea.Status == models.IdeaDuplicate && idea.Response != nil {
		respondedOn = idea.Response.RespondedOn
	}

	users, err := s.trx.QueryIntArray("SELECT user_id FROM idea_supporters WHERE idea_id = $1 AND tenant_id = $2", idea.ID, s.tenant.ID)
	if err != nil {
		return err
	}

	for _, u := range users {
		if err := s.AddSupporter(original.Number, u); err != nil {
			return err
		}
	}

	return s.trx.Execute(`
	UPDATE ideas 
	SET response = '', original_id = $3, response_date = $4, response_user_id = $5, status = $6 
	WHERE id = $1 and tenant_id = $2
	`, idea.ID, s.tenant.ID, original.ID, respondedOn, userID, models.IdeaDuplicate)
}

// IsReferenced returns true if another idea is referencing given idea
func (s *IdeaStorage) IsReferenced(number int) (bool, error) {
	return s.trx.Exists(`
		SELECT 1 FROM ideas i 
		INNER JOIN ideas o
		ON o.tenant_id = i.tenant_id
		AND o.id = i.original_id
		WHERE i.tenant_id = $1
		AND o.number = $2`, s.tenant.ID, number)
}

// SupportedBy returns a list of Idea ID supported by given user
func (s *IdeaStorage) SupportedBy(userID int) ([]int, error) {
	return s.trx.QueryIntArray("SELECT idea_id FROM idea_supporters WHERE user_id = $1 AND tenant_id = $2", userID, s.tenant.ID)
}
