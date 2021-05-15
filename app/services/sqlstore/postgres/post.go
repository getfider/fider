package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/gosimple/slug"
	"github.com/lib/pq"

	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
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

func (i *dbPost) toModel(ctx context.Context) *entity.Post {
	post := &entity.Post{
		ID:            i.ID,
		Number:        i.Number,
		Title:         i.Title,
		Slug:          i.Slug,
		Description:   i.Description,
		CreatedAt:     i.CreatedAt,
		HasVoted:      i.HasVoted,
		VotesCount:    i.VotesCount,
		CommentsCount: i.CommentsCount,
		Status:        enum.PostStatus(i.Status),
		User:          i.User.toModel(ctx),
		Tags:          i.Tags,
	}

	if i.Response.Valid {
		post.Response = &entity.PostResponse{
			Text:        i.Response.String,
			RespondedAt: i.RespondedAt.Time,
			User:        i.ResponseUser.toModel(ctx),
		}
		if post.Status == enum.PostDuplicate && i.OriginalNumber.Valid {
			post.Response.Original = &entity.OriginalPost{
				Number: int(i.OriginalNumber.Int64),
				Slug:   i.OriginalSlug.String,
				Title:  i.OriginalTitle.String,
				Status: enum.PostStatus(i.OriginalStatus.Int64),
			}
		}
	}
	return post
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
													WHERE p.status != ` + strconv.Itoa(int(enum.PostDeleted)) + ` AND %s`
)

func postIsReferenced(ctx context.Context, q *query.PostIsReferenced) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = false

		exists, err := trx.Exists(`
			SELECT 1 FROM posts p 
			INNER JOIN posts o
			ON o.tenant_id = p.tenant_id
			AND o.id = p.original_id
			WHERE p.tenant_id = $1
			AND o.id = $2`, tenant.ID, q.PostID)
		if err != nil {
			return errors.Wrap(err, "failed to check if post is referenced")
		}

		q.Result = exists
		return nil
	})
}

func setPostResponse(ctx context.Context, c *cmd.SetPostResponse) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if c.Status == enum.PostDuplicate {
			return errors.New("Use MarkAsDuplicate to change an post status to Duplicate")
		}

		respondedAt := time.Now()
		if c.Post.Status == c.Status && c.Post.Response != nil {
			respondedAt = c.Post.Response.RespondedAt
		}

		_, err := trx.Execute(`
		UPDATE posts 
		SET response = $3, original_id = NULL, response_date = $4, response_user_id = $5, status = $6 
		WHERE id = $1 and tenant_id = $2
		`, c.Post.ID, tenant.ID, c.Text, respondedAt, user.ID, c.Status)
		if err != nil {
			return errors.Wrap(err, "failed to update post's response")
		}

		c.Post.Status = c.Status
		c.Post.Response = &entity.PostResponse{
			Text:        c.Text,
			RespondedAt: respondedAt,
			User:        user,
		}
		return nil
	})
}

func markPostAsDuplicate(ctx context.Context, c *cmd.MarkPostAsDuplicate) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		respondedAt := time.Now()
		if c.Post.Status == enum.PostDuplicate && c.Post.Response != nil {
			respondedAt = c.Post.Response.RespondedAt
		}

		var users []*dbUser
		err := trx.Select(&users, "SELECT user_id AS id FROM post_votes WHERE post_id = $1 AND tenant_id = $2", c.Post.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to get votes of post with id '%d'", c.Post.ID)
		}

		for _, u := range users {
			err := bus.Dispatch(ctx, &cmd.AddVote{Post: c.Original, User: u.toModel(ctx)})
			if err != nil {
				return err
			}
		}

		_, err = trx.Execute(`
		UPDATE posts 
		SET response = '', original_id = $3, response_date = $4, response_user_id = $5, status = $6 
		WHERE id = $1 and tenant_id = $2
		`, c.Post.ID, tenant.ID, c.Original.ID, respondedAt, user.ID, enum.PostDuplicate)
		if err != nil {
			return errors.Wrap(err, "failed to update post's response")
		}

		c.Post.Status = enum.PostDuplicate
		c.Post.Response = &entity.PostResponse{
			RespondedAt: respondedAt,
			User:        user,
			Original: &entity.OriginalPost{
				Number: c.Original.Number,
				Title:  c.Original.Title,
				Slug:   c.Original.Slug,
				Status: c.Original.Status,
			},
		}
		return nil
	})
}

func countPostPerStatus(ctx context.Context, q *query.CountPostPerStatus) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {

		type dbStatusCount struct {
			Status enum.PostStatus `db:"status"`
			Count  int             `db:"count"`
		}

		q.Result = make(map[enum.PostStatus]int)
		stats := []*dbStatusCount{}
		err := trx.Select(&stats, "SELECT status, COUNT(*) AS count FROM posts WHERE tenant_id = $1 GROUP BY status", tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed to count posts per status")
		}

		for _, v := range stats {
			q.Result[v.Status] = v.Count
		}
		return nil
	})
}

func addNewPost(ctx context.Context, c *cmd.AddNewPost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var id int
		err := trx.Get(&id,
			`INSERT INTO posts (title, slug, number, description, tenant_id, user_id, created_at, status) 
			 VALUES ($1, $2, (SELECT COALESCE(MAX(number), 0) + 1 FROM posts p WHERE p.tenant_id = $4), $3, $4, $5, $6, 0) 
			 RETURNING id`, c.Title, slug.Make(c.Title), c.Description, tenant.ID, user.ID, time.Now())
		if err != nil {
			return errors.Wrap(err, "failed add new post")
		}

		q := &query.GetPostByID{PostID: id}
		if err := getPostByID(ctx, q); err != nil {
			return err
		}
		c.Result = q.Result

		if err := internalAddSubscriber(trx, q.Result, tenant, user, false); err != nil {
			return err
		}

		return nil
	})
}

func updatePost(ctx context.Context, c *cmd.UpdatePost) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute(`UPDATE posts SET title = $1, slug = $2, description = $3 
													 WHERE id = $4 AND tenant_id = $5`, c.Title, slug.Make(c.Title), c.Description, c.Post.ID, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update post")
		}

		q := &query.GetPostByID{PostID: c.Post.ID}
		if err := getPostByID(ctx, q); err != nil {
			return err
		}
		c.Result = q.Result
		return nil
	})
}

func getPostByID(ctx context.Context, q *query.GetPostByID) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		post, err := querySinglePost(ctx, trx, buildPostQuery(user, "p.tenant_id = $1 AND p.id = $2"), tenant.ID, q.PostID)
		if err != nil {
			return errors.Wrap(err, "failed to get post with id '%d'", q.PostID)
		}
		q.Result = post
		return nil
	})
}

func getPostBySlug(ctx context.Context, q *query.GetPostBySlug) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		post, err := querySinglePost(ctx, trx, buildPostQuery(user, "p.tenant_id = $1 AND p.slug = $2"), tenant.ID, q.Slug)
		if err != nil {
			return errors.Wrap(err, "failed to get post with slug '%s'", q.Slug)
		}
		q.Result = post
		return nil
	})
}

func getPostByNumber(ctx context.Context, q *query.GetPostByNumber) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		post, err := querySinglePost(ctx, trx, buildPostQuery(user, "p.tenant_id = $1 AND p.number = $2"), tenant.ID, q.Number)
		if err != nil {
			return errors.Wrap(err, "failed to get post with number '%d'", q.Number)
		}
		q.Result = post
		return nil
	})
}

func searchPosts(ctx context.Context, q *query.SearchPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		innerQuery := buildPostQuery(user, "p.tenant_id = $1 AND p.status = ANY($2)")

		if q.Tags == nil {
			q.Tags = []string{}
		}

		if q.Limit != "all" {
			if _, err := strconv.Atoi(q.Limit); err != nil {
				q.Limit = "30"
			}
		}

		var (
			posts []*dbPost
			err   error
		)
		if q.Query != "" {
			scoreField := "ts_rank(setweight(to_tsvector(title), 'A') || setweight(to_tsvector(description), 'B'), to_tsquery('english', $3)) + similarity(title, $4) + similarity(description, $4)"
			sql := fmt.Sprintf(`
				SELECT * FROM (%s) AS q 
				WHERE %s > 0.1
				ORDER BY %s DESC
				LIMIT %s
			`, innerQuery, scoreField, scoreField, q.Limit)
			err = trx.Select(&posts, sql, tenant.ID, pq.Array([]enum.PostStatus{
				enum.PostOpen,
				enum.PostStarted,
				enum.PostPlanned,
				enum.PostCompleted,
				enum.PostDeclined,
			}), ToTSQuery(q.Query), q.Query)
		} else {
			condition, statuses, sort := getViewData(q.View)
			sql := fmt.Sprintf(`
				SELECT * FROM (%s) AS q 
				WHERE tags @> $3 %s
				ORDER BY %s DESC
				LIMIT %s
			`, innerQuery, condition, sort, q.Limit)
			err = trx.Select(&posts, sql, tenant.ID, pq.Array(statuses), pq.Array(q.Tags))
		}

		if err != nil {
			return errors.Wrap(err, "failed to search posts")
		}

		q.Result = make([]*entity.Post, len(posts))
		for i, post := range posts {
			q.Result[i] = post.toModel(ctx)
		}
		return nil
	})
}

func getAllPosts(ctx context.Context, q *query.GetAllPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		searchQuery := &query.SearchPosts{View: "all", Limit: "all"}
		if err := searchPosts(ctx, searchQuery); err != nil {
			return errors.Wrap(err, "failed to get all posts")
		}
		q.Result = searchQuery.Result
		return nil
	})
}

func querySinglePost(ctx context.Context, trx *dbx.Trx, query string, args ...interface{}) (*entity.Post, error) {
	post := dbPost{}

	if err := trx.Get(&post, query, args...); err != nil {
		return nil, err
	}

	return post.toModel(ctx), nil
}

func buildPostQuery(user *entity.User, filter string) string {
	tagCondition := `AND tags.is_public = true`
	if user != nil && user.IsCollaborator() {
		tagCondition = ``
	}
	hasVotedSubQuery := "null"
	if user != nil {
		hasVotedSubQuery = fmt.Sprintf("(SELECT true FROM post_votes WHERE post_id = p.id AND user_id = %d)", user.ID)
	}
	return fmt.Sprintf(sqlSelectPostsWhere, tagCondition, hasVotedSubQuery, filter)
}
