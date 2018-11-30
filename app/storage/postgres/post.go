package postgres

import (
	"fmt"
	"strconv"
	"time"

	"database/sql"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/gosimple/slug"
	"github.com/lib/pq"
)

type dbPost struct {
	ID             int            `db:"id"`
	Number         int            `db:"number"`
	Title          string         `db:"title"`
	Slug           string         `db:"slug"`
	Description    string         `db:"description"`
	CreatedAt      time.Time      `db:"created_at"`
	User           *dbUser        `db:"user"`
	HasVoted       bool           `db:"has_voted"`
	VotesCount     int            `db:"votes_count"`
	CommentsCount  int            `db:"comments_count"`
	RecentVotes    int            `db:"recent_votes_count"`
	RecentComments int            `db:"recent_comments_count"`
	Status         int            `db:"status"`
	Response       sql.NullString `db:"response"`
	RespondedAt    dbx.NullTime   `db:"response_date"`
	ResponseUser   *dbUser        `db:"response_user"`
	OriginalNumber sql.NullInt64  `db:"original_number"`
	OriginalTitle  sql.NullString `db:"original_title"`
	OriginalSlug   sql.NullString `db:"original_slug"`
	OriginalStatus sql.NullInt64  `db:"original_status"`
	Tags           []string       `db:"tags"`
}

func (i *dbPost) toModel() *models.Post {
	post := &models.Post{
		ID:            i.ID,
		Number:        i.Number,
		Title:         i.Title,
		Slug:          i.Slug,
		Description:   i.Description,
		CreatedAt:     i.CreatedAt,
		HasVoted:      i.HasVoted,
		VotesCount:    i.VotesCount,
		CommentsCount: i.CommentsCount,
		Status:        models.PostStatus(i.Status),
		User:          i.User.toModel(),
		Tags:          i.Tags,
	}

	if i.Response.Valid {
		post.Response = &models.PostResponse{
			Text:        i.Response.String,
			RespondedAt: i.RespondedAt.Time,
			User:        i.ResponseUser.toModel(),
		}
		if post.Status == models.PostDuplicate && i.OriginalNumber.Valid {
			post.Response.Original = &models.OriginalPost{
				Number: int(i.OriginalNumber.Int64),
				Slug:   i.OriginalSlug.String,
				Title:  i.OriginalTitle.String,
				Status: models.PostStatus(i.OriginalStatus.Int64),
			}
		}
	}
	return post
}

type dbComment struct {
	ID        int          `db:"id"`
	Content   string       `db:"content"`
	CreatedAt time.Time    `db:"created_at"`
	User      *dbUser      `db:"user"`
	EditedAt  dbx.NullTime `db:"edited_at"`
	EditedBy  *dbUser      `db:"edited_by"`
}

func (c *dbComment) toModel() *models.Comment {
	comment := &models.Comment{
		ID:        c.ID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		User:      c.User.toModel(),
	}
	if c.EditedAt.Valid {
		comment.EditedBy = c.EditedBy.toModel()
		comment.EditedAt = &c.EditedAt.Time
	}
	return comment
}

type dbStatusCount struct {
	Status models.PostStatus `db:"status"`
	Count  int               `db:"count"`
}

// PostStorage contains read and write operations for posts
type PostStorage struct {
	trx    *dbx.Trx
	tenant *models.Tenant
	user   *models.User
}

// NewPostStorage creates a new PostStorage
func NewPostStorage(trx *dbx.Trx) *PostStorage {
	return &PostStorage{
		trx: trx,
	}
}

// SetCurrentTenant to current context
func (s *PostStorage) SetCurrentTenant(tenant *models.Tenant) {
	s.tenant = tenant
}

// SetCurrentUser to current context
func (s *PostStorage) SetCurrentUser(user *models.User) {
	s.user = user
}

var (
	sqlSelectPostsWhere = `	WITH 
													agg_comments AS (
															SELECT 
																	post_id, 
																	COUNT(CASE WHEN comments.created_at > CURRENT_DATE - INTERVAL '30 days' THEN 1 END) as recent,
																	COUNT(*) as all
															FROM comments 
															INNER JOIN posts
															ON posts.id = comments.post_id
															AND posts.tenant_id = comments.tenant_id
															WHERE posts.tenant_id = $1
															AND comments.deleted_at IS NULL
															GROUP BY post_id
													),
													agg_votes AS (
															SELECT 
															post_id, 
																	COUNT(CASE WHEN post_votes.created_at > CURRENT_DATE - INTERVAL '30 days'  THEN 1 END) as recent,
																	COUNT(*) as all
															FROM post_votes 
															INNER JOIN posts
															ON posts.id = post_votes.post_id
															AND posts.tenant_id = post_votes.tenant_id
															WHERE posts.tenant_id = $1
															GROUP BY post_id
													)
													SELECT p.id, 
																p.number, 
																p.title, 
																p.slug, 
																p.description, 
																p.created_at,
																COALESCE(agg_s.all, 0) as votes_count,
																COALESCE(agg_c.all, 0) as comments_count,
																COALESCE(agg_s.recent, 0) AS recent_votes_count,
																COALESCE(agg_c.recent, 0) AS recent_comments_count,																
																p.status, 
																u.id AS user_id, 
																u.name AS user_name, 
																u.email AS user_email,
																u.role AS user_role,
																u.status AS user_status,
																p.response,
																p.response_date,
																r.id AS response_user_id, 
																r.name AS response_user_name, 
																r.email AS response_user_email, 
																r.role AS response_user_role,
																r.status AS response_user_status,
																d.number AS original_number,
																d.title AS original_title,
																d.slug AS original_slug,
																d.status AS original_status,
																array_remove(array_agg(t.slug), NULL) AS tags,
																COALESCE(%s, false) AS has_voted
													FROM posts p
													INNER JOIN users u
													ON u.id = p.user_id
													LEFT JOIN users r
													ON r.id = p.response_user_id
													LEFT JOIN post_tags pt
													ON pt.post_id = p.id
													LEFT JOIN posts d
													ON d.id = p.original_id
													LEFT JOIN tags t
													ON t.id = pt.tag_id
													%s
													LEFT JOIN agg_comments agg_c
													ON agg_c.post_id = p.id
													LEFT JOIN agg_votes agg_s
													ON agg_s.post_id = p.id
													WHERE p.status != ` + strconv.Itoa(int(models.PostDeleted)) + ` AND %s
													GROUP BY p.id, u.id, r.id, d.id, agg_s.all, agg_c.all, agg_c.recent, agg_s.recent`
)

func (s *PostStorage) getPostQuery(filter string) string {
	hasVotedSubQuery := "null"
	if s.user != nil {
		hasVotedSubQuery = fmt.Sprintf("(SELECT true FROM post_votes WHERE post_id = p.id AND user_id = %d)", s.user.ID)
	}
	tagCondition := `AND t.is_public = true`
	if s.user != nil && s.user.IsCollaborator() {
		tagCondition = ``
	}
	return fmt.Sprintf(sqlSelectPostsWhere, hasVotedSubQuery, tagCondition, filter)
}

func (s *PostStorage) getSingle(query string, args ...interface{}) (*models.Post, error) {
	post := dbPost{}

	if err := s.trx.Get(&post, query, args...); err != nil {
		return nil, err
	}

	return post.toModel(), nil
}

// GetByID returns post by given id
func (s *PostStorage) GetByID(postID int) (*models.Post, error) {
	post, err := s.getSingle(s.getPostQuery("p.tenant_id = $1 AND p.id = $2"), s.tenant.ID, postID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get post with id '%d'", postID)
	}
	return post, nil
}

// GetBySlug returns post by tenant and slug
func (s *PostStorage) GetBySlug(slug string) (*models.Post, error) {
	post, err := s.getSingle(s.getPostQuery("p.tenant_id = $1 AND p.slug = $2"), s.tenant.ID, slug)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get post with slug '%s'", slug)
	}
	return post, nil
}

// GetByNumber returns post by tenant and number
func (s *PostStorage) GetByNumber(number int) (*models.Post, error) {
	post, err := s.getSingle(s.getPostQuery("p.tenant_id = $1 AND p.number = $2"), s.tenant.ID, number)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get post with number '%d'", number)
	}
	return post, nil
}

// GetAll returns all tenant posts
func (s *PostStorage) GetAll() ([]*models.Post, error) {
	posts, err := s.Search("", "all", "all", []string{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all posts")
	}
	return posts, nil
}

// CountPerStatus returns total number of posts per status
func (s *PostStorage) CountPerStatus() (map[models.PostStatus]int, error) {
	stats := []*dbStatusCount{}
	err := s.trx.Select(&stats, "SELECT status, COUNT(*) AS count FROM posts WHERE tenant_id = $1 GROUP BY status", s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count posts per status")
	}
	result := make(map[models.PostStatus]int, len(stats))
	for _, v := range stats {
		result[v.Status] = v.Count
	}
	return result, nil
}

// Search existing posts based on input
func (s *PostStorage) Search(query, view, limit string, tags []string) ([]*models.Post, error) {
	innerQuery := s.getPostQuery("p.tenant_id = $1 AND p.status = ANY($2)")

	if limit != "all" {
		if _, err := strconv.Atoi(limit); err != nil {
			limit = "30"
		}
	}

	var (
		posts []*dbPost
		err   error
	)
	if query != "" {
		scoreField := "ts_rank(setweight(to_tsvector(title), 'A') || setweight(to_tsvector(description), 'B'), to_tsquery('english', $3)) + similarity(title, $4) + similarity(description, $4)"
		sql := fmt.Sprintf(`
			SELECT * FROM (%s) AS q 
			WHERE %s > 0.1
			ORDER BY %s DESC
			LIMIT %s
		`, innerQuery, scoreField, scoreField, limit)
		err = s.trx.Select(&posts, sql, s.tenant.ID, pq.Array([]models.PostStatus{
			models.PostOpen,
			models.PostStarted,
			models.PostPlanned,
			models.PostCompleted,
			models.PostDeclined,
		}), ToTSQuery(query), query)
	} else {
		condition, statuses, sort := getViewData(view)
		sql := fmt.Sprintf(`
			SELECT * FROM (%s) AS q 
			WHERE tags @> $3 %s
			ORDER BY %s DESC
			LIMIT %s
		`, innerQuery, condition, sort, limit)
		err = s.trx.Select(&posts, sql, s.tenant.ID, pq.Array(statuses), pq.Array(tags))
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to search posts")
	}

	var result = make([]*models.Post, len(posts))
	for i, post := range posts {
		result[i] = post.toModel()
	}
	return result, nil
}

// GetCommentsByPost returns all comments from given post
func (s *PostStorage) GetCommentsByPost(post *models.Post) ([]*models.Comment, error) {
	comments := []*dbComment{}
	err := s.trx.Select(&comments,
		`SELECT c.id, 
				c.content, 
				c.created_at, 
				c.edited_at, 
				u.id AS user_id, 
				u.name AS user_name,
				u.email AS user_email,
				u.role AS user_role, 
				u.status AS user_status, 
				e.id AS edited_by_id, 
				e.name AS edited_by_name,
				e.email AS edited_by_email,
				e.role AS edited_by_role,
				e.status AS edited_by_status
		FROM comments c
		INNER JOIN posts p
		ON p.id = c.post_id
		AND p.tenant_id = c.tenant_id
		INNER JOIN users u
		ON u.id = c.user_id
		AND u.tenant_id = c.tenant_id
		LEFT JOIN users e
		ON e.id = c.edited_by_id
		AND e.tenant_id = c.tenant_id
		WHERE p.id = $1
		AND p.tenant_id = $2
		AND c.deleted_at IS NULL
		ORDER BY c.created_at ASC`, post.ID, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed get comments of post with id '%d'", post.ID)
	}

	var result = make([]*models.Comment, len(comments))
	for i, comment := range comments {
		result[i] = comment.toModel()
	}
	return result, nil
}

// Update given post
func (s *PostStorage) Update(post *models.Post, title, description string) (*models.Post, error) {
	_, err := s.trx.Execute(`UPDATE posts SET title = $1, slug = $2, description = $3 
													 WHERE id = $4 AND tenant_id = $5`, title, slug.Make(title), description, post.ID, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed update post")
	}

	post.Slug = slug.Make(title)
	post.Title = title
	post.Description = description

	return post, nil
}

// Add a new post in the database
func (s *PostStorage) Add(title, description string) (*models.Post, error) {
	var id int
	err := s.trx.Get(&id,
		`INSERT INTO posts (title, slug, number, description, tenant_id, user_id, created_at, status) 
		 VALUES ($1, $2, (SELECT COALESCE(MAX(number), 0) + 1 FROM posts p WHERE p.tenant_id = $4), $3, $4, $5, $6, 0) 
		 RETURNING id`, title, slug.Make(title), description, s.tenant.ID, s.user.ID, time.Now())
	if err != nil {
		return nil, errors.Wrap(err, "failed add new post")
	}

	post, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := s.internalAddSubscriber(post, s.user, false); err != nil {
		return nil, err
	}

	return post, nil
}

// AddComment places a new comment on an post
func (s *PostStorage) AddComment(post *models.Post, content string) (int, error) {
	var id int
	if err := s.trx.Get(&id, `
		INSERT INTO comments (tenant_id, post_id, content, user_id, created_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`, s.tenant.ID, post.ID, content, s.user.ID, time.Now()); err != nil {
		return 0, errors.Wrap(err, "failed add new comment")
	}

	if err := s.internalAddSubscriber(post, s.user, false); err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteComment by its id
func (s *PostStorage) DeleteComment(id int) error {
	if _, err := s.trx.Execute(
		"UPDATE comments SET deleted_at = $1, deleted_by_id = $2 WHERE id = $3 AND tenant_id = $4",
		time.Now(), s.user.ID, id, s.tenant.ID,
	); err != nil {
		return errors.Wrap(err, "failed delete comment")
	}
	return nil
}

// GetCommentByID returns a comment by given ID
func (s *PostStorage) GetCommentByID(id int) (*models.Comment, error) {
	comment := dbComment{}
	err := s.trx.Get(&comment,
		`SELECT c.id, 
						c.content, 
						c.created_at, 
						c.edited_at, 
						u.id AS user_id, 
						u.name AS user_name,
						u.email AS user_email,
						u.role AS user_role, 
						u.status AS user_status, 
						e.id AS edited_by_id, 
						e.name AS edited_by_name,
						e.email AS edited_by_email,
						e.role AS edited_by_role,
						e.status AS edited_by_status
		FROM comments c
		INNER JOIN users u
		ON u.id = c.user_id
		AND u.tenant_id = c.tenant_id
		LEFT JOIN users e
		ON e.id = c.edited_by_id
		AND e.tenant_id = c.tenant_id
		WHERE c.id = $1
		AND c.tenant_id = $2
		AND c.deleted_at IS NULL`, id, s.tenant.ID)

	if err != nil {
		return nil, err
	}

	return comment.toModel(), nil
}

// UpdateComment with given ID and content
func (s *PostStorage) UpdateComment(id int, content string) error {
	_, err := s.trx.Execute(`
		UPDATE comments SET content = $1, edited_at = $2, edited_by_id = $3 
		WHERE id = $4 AND tenant_id = $5`, content, time.Now(), s.user.ID, id, s.tenant.ID)
	if err != nil {
		return errors.Wrap(err, "failed update comment")
	}
	return nil
}

// AddVote adds user to post list of votes
func (s *PostStorage) AddVote(post *models.Post, user *models.User) error {
	if !post.CanBeVoted() {
		return nil
	}

	_, err := s.trx.Execute(
		`INSERT INTO post_votes (tenant_id, user_id, post_id, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`,
		s.tenant.ID, user.ID, post.ID, time.Now())

	if err != nil {
		return errors.Wrap(err, "failed add vote to post")
	}

	return s.internalAddSubscriber(post, user, false)
}

// RemoveVote removes user from post list of votes
func (s *PostStorage) RemoveVote(post *models.Post, user *models.User) error {
	if !post.CanBeVoted() {
		return nil
	}

	_, err := s.trx.Execute(`DELETE FROM post_votes WHERE user_id = $1 AND post_id = $2 AND tenant_id = $3`, user.ID, post.ID, s.tenant.ID)
	if err != nil {
		return errors.Wrap(err, "failed to remove vote from post")
	}

	return err
}

// AddSubscriber adds user to the post list of subscribers
func (s *PostStorage) AddSubscriber(post *models.Post, user *models.User) error {
	return s.internalAddSubscriber(post, user, true)
}

func (s *PostStorage) internalAddSubscriber(post *models.Post, user *models.User, force bool) error {
	conflict := " DO NOTHING"
	if force {
		conflict = "(user_id, post_id) DO UPDATE SET status = $5, updated_at = $4"
	}

	_, err := s.trx.Execute(fmt.Sprintf(`
	INSERT INTO post_subscribers (tenant_id, user_id, post_id, created_at, updated_at, status)
	VALUES ($1, $2, $3, $4, $4, $5)  ON CONFLICT %s`, conflict),
		s.tenant.ID, user.ID, post.ID, time.Now(), models.SubscriberActive,
	)
	if err != nil {
		return errors.Wrap(err, "failed insert post subscriber")
	}
	return nil
}

// RemoveSubscriber removes user from post list of subscribers
func (s *PostStorage) RemoveSubscriber(post *models.Post, user *models.User) error {
	_, err := s.trx.Execute(`
		INSERT INTO post_subscribers (tenant_id, user_id, post_id, created_at, updated_at, status)
		VALUES ($1, $2, $3, $4, $4, 0) ON CONFLICT (user_id, post_id)
		DO UPDATE SET status = 0, updated_at = $4`,
		s.tenant.ID, user.ID, post.ID, time.Now(),
	)
	if err != nil {
		return errors.Wrap(err, "failed remove post subscriber")
	}
	return nil
}

// GetActiveSubscribers based on input and settings
func (s *PostStorage) GetActiveSubscribers(number int, channel models.NotificationChannel, event models.NotificationEvent) ([]*models.User, error) {
	var (
		users []*dbUser
		err   error
	)

	if len(event.RequiresSubscriptionUserRoles) == 0 {
		err = s.trx.Select(&users, `
			SELECT DISTINCT u.id, u.name, u.email, u.tenant_id, u.role, u.status
			FROM users u
			LEFT JOIN user_settings set
			ON set.user_id = u.id
			AND set.tenant_id = u.tenant_id
			AND set.key = $1
			WHERE u.tenant_id = $2
			AND u.status = $5
			AND (
				(set.value IS NULL AND u.role = ANY($3))
				OR CAST(set.value AS integer) & $4 > 0
			)
			ORDER by u.id`,
			event.UserSettingsKeyName,
			s.tenant.ID,
			pq.Array(event.DefaultEnabledUserRoles),
			channel,
			models.UserActive,
		)
	} else {
		err = s.trx.Select(&users, `
			SELECT DISTINCT u.id, u.name, u.email, u.tenant_id, u.role, u.status
			FROM users u
			LEFT JOIN post_subscribers sub
			ON sub.user_id = u.id
			AND sub.post_id = (SELECT id FROM posts p WHERE p.tenant_id = $4 and p.number = $1 LIMIT 1)
			AND sub.tenant_id = u.tenant_id
			LEFT JOIN user_settings set
			ON set.user_id = u.id
			AND set.key = $3
			AND set.tenant_id = u.tenant_id
			WHERE u.tenant_id = $4
			AND u.status = $8
			AND ( sub.status = $2 OR (sub.status IS NULL AND NOT u.role = ANY($7)) )
			AND (
				(set.value IS NULL AND u.role = ANY($5))
				OR CAST(set.value AS integer) & $6 > 0
			)
			ORDER by u.id`,
			number,
			models.SubscriberActive,
			event.UserSettingsKeyName,
			s.tenant.ID,
			pq.Array(event.DefaultEnabledUserRoles),
			channel,
			pq.Array(event.RequiresSubscriptionUserRoles),
			models.UserActive,
		)
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get post number '%d' subscribers", number)
	}

	var result = make([]*models.User, len(users))
	for i, user := range users {
		result[i] = user.toModel()
	}
	return result, nil
}

// SetResponse changes current post response
func (s *PostStorage) SetResponse(post *models.Post, text string, status models.PostStatus) error {
	if status == models.PostDuplicate {
		return errors.New("Use MarkAsDuplicate to change an post status to Duplicate")
	}

	respondedAt := time.Now()
	if post.Status == status && post.Response != nil {
		respondedAt = post.Response.RespondedAt
	}

	_, err := s.trx.Execute(`
	UPDATE posts 
	SET response = $3, original_id = NULL, response_date = $4, response_user_id = $5, status = $6 
	WHERE id = $1 and tenant_id = $2
	`, post.ID, s.tenant.ID, text, respondedAt, s.user.ID, status)
	if err != nil {
		return errors.Wrap(err, "failed to update post's response")
	}

	post.Status = status
	post.Response = &models.PostResponse{
		Text:        text,
		RespondedAt: respondedAt,
		User:        s.user,
	}
	return nil
}

// MarkAsDuplicate set post as a duplicate of another post
func (s *PostStorage) MarkAsDuplicate(post *models.Post, original *models.Post) error {
	respondedAt := time.Now()
	if post.Status == models.PostDuplicate && post.Response != nil {
		respondedAt = post.Response.RespondedAt
	}

	var users []*dbUser
	err := s.trx.Select(&users, "SELECT user_id AS id FROM post_votes WHERE post_id = $1 AND tenant_id = $2", post.ID, s.tenant.ID)
	if err != nil {
		return errors.Wrap(err, "failed to get votes of post with id '%d'", post.ID)
	}

	for _, u := range users {
		if err := s.AddVote(original, u.toModel()); err != nil {
			return err
		}
	}

	_, err = s.trx.Execute(`
	UPDATE posts 
	SET response = '', original_id = $3, response_date = $4, response_user_id = $5, status = $6 
	WHERE id = $1 and tenant_id = $2
	`, post.ID, s.tenant.ID, original.ID, respondedAt, s.user.ID, models.PostDuplicate)
	if err != nil {
		return errors.Wrap(err, "failed to update post's response")
	}

	post.Status = models.PostDuplicate
	post.Response = &models.PostResponse{
		RespondedAt: respondedAt,
		User:        s.user,
		Original: &models.OriginalPost{
			Number: original.Number,
			Title:  original.Title,
			Slug:   original.Slug,
			Status: original.Status,
		},
	}
	return nil
}

// IsReferenced returns true if another post is referencing given post
func (s *PostStorage) IsReferenced(post *models.Post) (bool, error) {
	exists, err := s.trx.Exists(`
		SELECT 1 FROM posts p 
		INNER JOIN posts o
		ON o.tenant_id = p.tenant_id
		AND o.id = p.original_id
		WHERE p.tenant_id = $1
		AND o.id = $2`, s.tenant.ID, post.ID)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if post is referenced")
	}
	return exists, nil
}

// VotedBy returns a list of Post ID voted by given user
func (s *PostStorage) VotedBy() ([]int, error) {
	posts, err := s.trx.QueryIntArray("SELECT post_id FROM post_votes WHERE user_id = $1 AND tenant_id = $2", s.user.ID, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user's voted posts")
	}
	return posts, nil
}

// ListVotes returns a list of all votes on given post
func (s *PostStorage) ListVotes(post *models.Post, limit int) ([]*models.Vote, error) {
	sqlLimit := "ALL"
	if limit > 0 {
		sqlLimit = strconv.Itoa(limit)
	}

	votes := []*models.Vote{}
	err := s.trx.Select(&votes, `
		SELECT 
			pv.created_at, 
			u.id AS user_id,
			u.name AS user_name,
			u.email AS user_email
		FROM post_votes pv
		INNER JOIN users u
		ON u.id = pv.user_id
		AND u.tenant_id = pv.tenant_id 
		WHERE pv.post_id = $1  
		AND pv.tenant_id = $2
		ORDER BY pv.created_at
		LIMIT `+sqlLimit, post.ID, s.tenant.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get votes of post")
	}
	return votes, nil
}
