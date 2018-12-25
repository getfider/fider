package inmemory

import (
	"time"

	"github.com/gosimple/slug"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models"
)

// PostStorage contains read and write operations for posts
type PostStorage struct {
	lastID          int
	lastCommentID   int
	posts           []*models.Post
	postsVotedBy    map[int][]int
	postSubscribers map[int][]*models.User
	postComments    map[int][]*models.Comment
	tenant          *models.Tenant
	user            *models.User
}

// NewPostStorage creates a new PostStorage
func NewPostStorage() *PostStorage {
	return &PostStorage{
		posts:           make([]*models.Post, 0),
		postsVotedBy:    make(map[int][]int, 0),
		postSubscribers: make(map[int][]*models.User, 0),
		postComments:    make(map[int][]*models.Comment, 0),
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

// GetByID returns post by given id
func (s *PostStorage) GetByID(postID int) (*models.Post, error) {
	for _, post := range s.posts {
		if post.ID == postID {
			return post, nil
		}
	}
	return nil, app.ErrNotFound
}

// Update given post
func (s *PostStorage) Update(post *models.Post, title, description string) (*models.Post, error) {
	post.Title = title
	post.Description = description
	return post, nil
}

// GetByNumber returns post by tenant and number
func (s *PostStorage) GetByNumber(number int) (*models.Post, error) {
	for _, post := range s.posts {
		if post.Number == number && post.Status != models.PostDeleted {
			return post, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetBySlug returns post by tenant and slug
func (s *PostStorage) GetBySlug(slug string) (*models.Post, error) {
	for _, post := range s.posts {
		if post.Slug == slug && post.Status != models.PostDeleted {
			return post, nil
		}
	}
	return nil, app.ErrNotFound
}

// GetAll returns all tenant posts
func (s *PostStorage) GetAll() ([]*models.Post, error) {
	return s.posts, nil
}

// CountPerStatus returns total number of posts per status
func (s *PostStorage) CountPerStatus() (map[models.PostStatus]int, error) {
	return make(map[models.PostStatus]int, 0), nil
}

// Search existing posts based on input
func (s *PostStorage) Search(query, view, limit string, tags []string) ([]*models.Post, error) {
	return s.posts, nil
}

// GetCommentsByPost returns all comments from given post
func (s *PostStorage) GetCommentsByPost(post *models.Post) ([]*models.Comment, error) {
	return s.postComments[post.ID], nil
}

// Add a new post in the database
func (s *PostStorage) Add(title, description string) (*models.Post, error) {
	s.lastID = s.lastID + 1
	post := &models.Post{
		ID:          s.lastID,
		Number:      s.lastID,
		Title:       title,
		Slug:        slug.Make(title),
		Description: description,
		User:        s.user,
	}
	s.posts = append(s.posts, post)
	s.postsVotedBy[s.user.ID] = append(s.postsVotedBy[s.user.ID], post.ID)
	s.AddSubscriber(post, s.user)
	return post, nil
}

// AddComment places a new comment on an post
func (s *PostStorage) AddComment(post *models.Post, content string) (int, error) {
	s.lastCommentID++
	s.postComments[post.ID] = append(s.postComments[post.ID], &models.Comment{
		ID:        s.lastCommentID,
		Content:   content,
		CreatedAt: time.Now(),
		User:      s.user,
	})

	return s.lastCommentID, nil
}

// DeleteComment by its id
func (s *PostStorage) DeleteComment(id int) error {
	for _, p := range s.posts {
		if comments, ok := s.postComments[p.ID]; ok {
			for idx, comment := range comments {
				if comment.ID == id {
					s.postComments[p.ID] = append(comments[:idx], comments[idx+1:]...)
				}
			}
		}
	}

	return nil
}

// GetCommentByID returns a comment by given ID
func (s *PostStorage) GetCommentByID(id int) (*models.Comment, error) {
	for _, comments := range s.postComments {
		for _, comment := range comments {
			if comment.ID == id && comment.User.Tenant == s.tenant {
				return comment, nil
			}
		}
	}
	return nil, app.ErrNotFound
}

// UpdateComment with given ID and content
func (s *PostStorage) UpdateComment(id int, content string) error {
	now := time.Now()
	comment, err := s.GetCommentByID(id)
	if err != nil {
		return err
	}
	comment.Content = content
	comment.EditedAt = &now
	comment.EditedBy = s.user
	return nil
}

// AddVote adds user to post list of votes
func (s *PostStorage) AddVote(post *models.Post, user *models.User) error {
	s.postsVotedBy[user.ID] = append(s.postsVotedBy[user.ID], post.ID)
	post.VotesCount = post.VotesCount + 1
	return nil
}

// RemoveVote removes user from post list of votes
func (s *PostStorage) RemoveVote(post *models.Post, user *models.User) error {
	for i, id := range s.postsVotedBy[user.ID] {
		if id == post.ID {
			s.postsVotedBy[user.ID] = append(s.postsVotedBy[user.ID][:i], s.postsVotedBy[user.ID][i+1:]...)
			break
		}
	}
	post.VotesCount = post.VotesCount - 1
	return nil
}

// SetResponse changes current post response
func (s *PostStorage) SetResponse(post *models.Post, text string, status models.PostStatus) error {
	for _, storedPost := range s.posts {
		if storedPost.Number == post.Number {
			storedPost.Status = status
			storedPost.Response = &models.PostResponse{
				Text:        text,
				User:        s.user,
				RespondedAt: time.Now(),
			}
		}
	}
	return nil
}

// MarkAsDuplicate set post as a duplicate of another post
func (s *PostStorage) MarkAsDuplicate(post *models.Post, original *models.Post) error {
	post.Status = models.PostDuplicate
	post.Response = &models.PostResponse{
		Original: &models.OriginalPost{
			Number: original.Number,
			Title:  original.Title,
			Slug:   original.Slug,
			Status: original.Status,
		},
		Text:        "",
		User:        s.user,
		RespondedAt: time.Now(),
	}
	return nil
}

// IsReferenced returns true if another post is referencing given post
func (s *PostStorage) IsReferenced(post *models.Post) (bool, error) {
	for _, i := range s.posts {
		if i.Status == models.PostDuplicate && i.Response.Original.Number == post.Number {
			return true, nil
		}
	}
	return false, nil
}

// VotedBy returns a list of Post ID voted by given user
func (s *PostStorage) VotedBy() ([]int, error) {
	return s.postsVotedBy[s.user.ID], nil
}

// ListVotes returns a list of all votes on given post
func (s *PostStorage) ListVotes(post *models.Post, limit int) ([]*models.Vote, error) {
	return make([]*models.Vote, 0), nil
}

// AddSubscriber adds user to the post list of subscribers
func (s *PostStorage) AddSubscriber(post *models.Post, user *models.User) error {
	s.postSubscribers[post.ID] = append(s.postSubscribers[post.ID], user)
	return nil
}

// RemoveSubscriber removes user from post list of subscribers
func (s *PostStorage) RemoveSubscriber(post *models.Post, user *models.User) error {
	for i, u := range s.postSubscribers[post.ID] {
		if u.ID == user.ID {
			s.postSubscribers[post.ID] = append(s.postSubscribers[post.ID][:i], s.postSubscribers[post.ID][i+1:]...)
			break
		}
	}
	return nil
}

// GetActiveSubscribers based on input and settings
func (s *PostStorage) GetActiveSubscribers(number int, channel models.NotificationChannel, event models.NotificationEvent) ([]*models.User, error) {
	var post *models.Post
	for _, p := range s.posts {
		if p.Number == number {
			post = p
		}
	}

	if post == nil {
		return nil, app.ErrNotFound
	}

	subscribers, ok := s.postSubscribers[post.ID]
	if ok {
		users := make([]*models.User, len(subscribers))
		for i, user := range subscribers {
			users[i] = user
		}
		return users, nil
	}
	return make([]*models.User, 0), nil
}
