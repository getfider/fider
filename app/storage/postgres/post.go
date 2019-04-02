package postgres

import (
	"context"
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

func (i *dbPost) toModel(ctx context.Context) *models.Post {
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
		User:          i.User.toModel(ctx),
		Tags:          i.Tags,
	}

	if i.Response.Valid {
		post.Response = &models.PostResponse{
			Text:        i.Response.String,
			RespondedAt: i.RespondedAt.Time,
			User:        i.ResponseUser.toModel(ctx),
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
	ID          int          `db:"id"`
	Content     string       `db:"content"`
	CreatedAt   time.Time    `db:"created_at"`
	User        *dbUser      `db:"user"`
	Attachments []string     `db:"attachment_bkeys"`
	EditedAt    dbx.NullTime `db:"edited_at"`
	EditedBy    *dbUser      `db:"edited_by"`
}

func (c *dbComment) toModel(ctx context.Context) *models.Comment {
	comment := &models.Comment{
		ID:          c.ID,
		Content:     c.Content,
		CreatedAt:   c.CreatedAt,
		User:        c.User.toModel(ctx),
		Attachments: c.Attachments,
	}
	if c.EditedAt.Valid {
		comment.EditedBy = c.EditedBy.toModel(ctx)
		comment.EditedAt = &c.EditedAt.Time
	}
	return comment
}

// PostStorage contains read and write operations for posts
type PostStorage struct {
	trx    *dbx.Trx
	tenant *models.Tenant
	user   *models.User
	ctx    context.Context
}

// NewPostStorage creates a new PostStorage
func NewPostStorage(trx *dbx.Trx, ctx context.Context) *PostStorage {
	return &PostStorage{
		trx: trx,
		ctx: ctx,
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
													agg_tags AS ( 
														SELECT 
																post_id, 
																ARRAY_REMOVE(ARRAY_AGG(tags.slug), NULL) as tags
														FROM post_tags
														INNER JOIN tags
														ON tags.ID = post_tags.TAG_ID
														AND tags.tenant_id = post_tags.tenant_id
														WHERE post_tags.tenant_id = $1
														%s
														GROUP BY post_id 
													), 
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
																u.avatar_type AS user_avatar_type,
																u.avatar_bkey AS user_avatar_bkey,
																p.response,
																p.response_date,
																r.id AS response_user_id, 
																r.name AS response_user_name, 
																r.email AS response_user_email, 
																r.role AS response_user_role,
																r.status AS response_user_status,
																r.avatar_type AS response_user_avatar_type,
																r.avatar_bkey AS response_user_avatar_bkey,
																d.number AS original_number,
																d.title AS original_title,
																d.slug AS original_slug,
																d.status AS original_status,
																COALESCE(agg_t.tags, ARRAY[]::text[]) AS tags,
																COALESCE(%s, false) AS has_voted
													FROM posts p
													INNER JOIN users u
													ON u.id = p.user_id
													AND u.tenant_id = $1
													LEFT JOIN users r
													ON r.id = p.response_user_id
													AND r.tenant_id = $1
													LEFT JOIN posts d
													ON d.id = p.original_id
													AND d.tenant_id = $1
													LEFT JOIN agg_comments agg_c
													ON agg_c.post_id = p.id
													LEFT JOIN agg_votes agg_s
													ON agg_s.post_id = p.id
													LEFT JOIN agg_tags agg_t 
													ON agg_t.post_id = p.id
													WHERE p.status != ` + strconv.Itoa(int(models.PostDeleted)) + ` AND %s`
)

func (s *PostStorage) getPostQuery(filter string) string {
	tagCondition := `AND tags.is_public = true`
	if s.user != nil && s.user.IsCollaborator() {
		tagCondition = ``
	}
	hasVotedSubQuery := "null"
	if s.user != nil {
		hasVotedSubQuery = fmt.Sprintf("(SELECT true FROM post_votes WHERE post_id = p.id AND user_id = %d)", s.user.ID)
	}
	return fmt.Sprintf(sqlSelectPostsWhere, tagCondition, hasVotedSubQuery, filter)
}

func (s *PostStorage) getSingle(query string, args ...interface{}) (*models.Post, error) {
	post := dbPost{}

	if err := s.trx.Get(&post, query, args...); err != nil {
		return nil, err
	}

	return post.toModel(s.ctx), nil
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
		result[i] = post.toModel(s.ctx)
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
