package apiv1

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/metrics"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// SearchPosts return existing posts based on search criteria
func SearchPosts() web.HandlerFunc {
	return func(c *web.Context) error {
		viewQueryParams := c.QueryParam("view")
		if viewQueryParams == "" {
			viewQueryParams = "all" // Set default value to "all" if not provided
		}
		searchPosts := &query.SearchPosts{
			Query: c.QueryParam("query"),
			View:  viewQueryParams,
			Limit: c.QueryParam("limit"),
			Tags:  c.QueryParamAsArray("tags"),
		}
		if myVotesOnly, err := c.QueryParamAsBool("myvotes"); err == nil {
			searchPosts.MyVotesOnly = myVotesOnly
		}
		if noTagsOnly, err := c.QueryParamAsBool("notags"); err == nil {
			searchPosts.NoTagsOnly = noTagsOnly
		}
		if myPostsOnly, err := c.QueryParamAsBool("myposts"); err == nil {
			searchPosts.MyPostsOnly = myPostsOnly
		}
		searchPosts.SetStatusesFromStrings(c.QueryParamAsArray("statuses"))

		if err := bus.Dispatch(c, searchPosts); err != nil {
			return c.Failure(err)
		}

		return c.Ok(searchPosts.Result)
	}
}

// FindSimilarPosts return posts similar to query
func FindSimilarPosts() web.HandlerFunc {
	return func(c *web.Context) error {
		searchPosts := &query.FindSimilarPosts{
			Query: c.QueryParam("query"),
		}
		if err := bus.Dispatch(c, searchPosts); err != nil {
			return c.Failure(err)
		}
		return c.Ok(searchPosts.Result)
	}
}

// CreatePost creates a new post on current tenant
func CreatePost() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.CreateNewPost)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		if err := bus.Dispatch(c, &cmd.UploadImages{Images: action.Attachments, Folder: "attachments"}); err != nil {
			return c.Failure(err)
		}

		newPost := &cmd.AddNewPost{
			Title:       action.Title,
			Description: action.Description,
		}
		err := bus.Dispatch(c, newPost)
		if err != nil {
			return c.Failure(err)
		}

		setAttachments := &cmd.SetAttachments{Post: newPost.Result, Attachments: action.Attachments}
		addVote := &cmd.AddVote{Post: newPost.Result, User: c.User()}
		if err = bus.Dispatch(c, setAttachments, addVote); err != nil {
			return c.Failure(err)
		}

		if env.Config.PostCreationWithTagsEnabled {
			for _, tag := range action.Tags {
				assignTag := &cmd.AssignTag{Tag: tag, Post: newPost.Result}
				if err := bus.Dispatch(c, assignTag); err != nil {
					return c.Failure(err)
				}
			}
		}

		c.Enqueue(tasks.NotifyAboutNewPost(newPost.Result))

		metrics.TotalPosts.Inc()
		return c.Ok(web.Map{
			"id":     newPost.Result.ID,
			"number": newPost.Result.Number,
			"title":  newPost.Result.Title,
			"slug":   newPost.Result.Slug,
		})
	}
}

// GetPost retrieves the existing post by number
func GetPost() web.HandlerFunc {
	return func(c *web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.NotFound()
		}

		getPost := &query.GetPostByNumber{Number: number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		return c.Ok(getPost.Result)
	}
}

// UpdatePost updates an existing post of current tenant
func UpdatePost() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.UpdatePost)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		updatePost := &cmd.UpdatePost{
			Post:        action.Post,
			Title:       action.Title,
			Description: action.Description,
		}

		err := bus.Dispatch(c,
			&cmd.UploadImages{
				Images: action.Attachments,
				Folder: "attachments",
			},
			updatePost,
			&cmd.SetAttachments{
				Post:        action.Post,
				Attachments: action.Attachments,
			},
		)

		if err != nil {
			return c.Failure(err)
		}

		// Notify about mentions in the updated post
		c.Enqueue(tasks.NotifyAboutUpdatedPost(updatePost.Result))

		return c.Ok(web.Map{})
	}
}

// SetResponse changes current post staff response
func SetResponse() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.SetResponse)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		getPost := &query.GetPostByNumber{Number: action.Number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		prevStatus := getPost.Result.Status

		var command bus.Msg
		if action.Status == enum.PostDuplicate {
			command = &cmd.MarkPostAsDuplicate{Post: getPost.Result, Original: action.Original}
		} else {
			command = &cmd.SetPostResponse{
				Post:   getPost.Result,
				Text:   action.Text,
				Status: action.Status,
			}
		}

		if err := bus.Dispatch(c, command); err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutStatusChange(getPost.Result, prevStatus))

		return c.Ok(web.Map{})
	}
}

// DeletePost deletes an existing post of current tenant
func DeletePost() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.DeletePost)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.SetPostResponse{
			Post:   action.Post,
			Text:   action.Text,
			Status: enum.PostDeleted,
		})
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutDeletedPost(action.Post, action.Text != ""))

		return c.Ok(web.Map{})
	}
}

// ListComments returns a list of all comments of a post
func ListComments() web.HandlerFunc {
	return func(c *web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.NotFound()
		}

		getPost := &query.GetPostByNumber{Number: number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		getComments := &query.GetCommentsByPost{Post: getPost.Result}
		if err := bus.Dispatch(c, getComments); err != nil {
			return c.Failure(err)
		}

		for _, comment := range getComments.Result {
			commentString := entity.CommentString(comment.Content)
			comment.Content = commentString.SanitizeMentions()
		}

		return c.Ok(getComments.Result)
	}
}

// GetComment returns a single comment by its ID
func GetComment() web.HandlerFunc {
	return func(c *web.Context) error {
		id, err := c.ParamAsInt("id")
		if err != nil {
			return c.NotFound()
		}

		commentByID := &query.GetCommentByID{CommentID: id}
		if err := bus.Dispatch(c, commentByID); err != nil {
			return c.Failure(err)
		}

		commentString := entity.CommentString(commentByID.Result.Content)
		commentByID.Result.Content = commentString.SanitizeMentions()

		return c.Ok(commentByID.Result)
	}
}

// ToggleReaction adds or removes a reaction on a comment
func ToggleReaction() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.ToggleCommentReaction)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		getComment := &query.GetCommentByID{CommentID: action.Comment}
		if err := bus.Dispatch(c, getComment); err != nil {
			return c.Failure(err)
		}

		toggleReaction := &cmd.ToggleCommentReaction{
			Comment: getComment.Result,
			Emoji:   action.Reaction,
			User:    c.User(),
		}
		if err := bus.Dispatch(c, toggleReaction); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{
			"added": toggleReaction.Result,
		})
	}
}

// PostComment creates a new comment on given post
func PostComment() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.AddNewComment)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		getPost := &query.GetPostByNumber{Number: action.Number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		if err := bus.Dispatch(c, &cmd.UploadImages{Images: action.Attachments, Folder: "attachments"}); err != nil {
			return c.Failure(err)
		}

		addNewComment := &cmd.AddNewComment{
			Post:    getPost.Result,
			Content: action.Content,
		}
		if err := bus.Dispatch(c, addNewComment); err != nil {
			return c.Failure(err)
		}

		// For processing, restore the original content
		addNewComment.Result.Content = action.Content

		if err := bus.Dispatch(c, &cmd.SetAttachments{
			Post:        getPost.Result,
			Comment:     addNewComment.Result,
			Attachments: action.Attachments,
		}); err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewComment(addNewComment.Result, getPost.Result))

		metrics.TotalComments.Inc()
		return c.Ok(web.Map{
			"id": addNewComment.Result.ID,
		})
	}
}

// UpdateComment changes an existing comment with new content
func UpdateComment() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.EditComment)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		getPost := &query.GetPostByID{PostID: action.Post.ID}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		comment := &entity.Comment{
			ID:      action.ID,
			Content: action.Content,
		}

		err := bus.Dispatch(c,
			&cmd.UploadImages{
				Images: action.Attachments,
				Folder: "attachments",
			},
			&cmd.UpdateComment{
				CommentID: action.ID,
				Content:   action.Content,
			},
			&cmd.SetAttachments{
				Post:        action.Post,
				Comment:     action.Comment,
				Attachments: action.Attachments,
			},
		)
		if err != nil {
			return c.Failure(err)
		}

		// Update the content

		c.Enqueue(tasks.NotifyAboutUpdatedComment(getPost.Result, comment))

		return c.Ok(web.Map{})
	}
}

// DeleteComment deletes an existing comment by its ID
func DeleteComment() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.DeleteComment)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.DeleteComment{
			CommentID: action.CommentID,
		})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// AddVote adds current user to given post list of votes
func AddVote() web.HandlerFunc {
	return func(c *web.Context) error {
		err := addOrRemove(c, func(post *entity.Post, user *entity.User) bus.Msg {
			return &cmd.AddVote{Post: post, User: user}
		})

		if err == nil {
			metrics.TotalVotes.Inc()
		}

		return err
	}
}

// RemoveVote removes current user from given post list of votes
func RemoveVote() web.HandlerFunc {
	return func(c *web.Context) error {
		return addOrRemove(c, func(post *entity.Post, user *entity.User) bus.Msg {
			return &cmd.RemoveVote{Post: post, User: user}
		})
	}
}

func ToggleVote() web.HandlerFunc {
	return func(c *web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.NotFound()
		}

		getPost := &query.GetPostByNumber{Number: number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		if getPost.Result == nil {
			return c.NotFound()
		}

		listVotes := &query.ListPostVotes{PostID: getPost.Result.ID}
		if err := bus.Dispatch(c, listVotes); err != nil {
			return c.Failure(err)
		}

		hasVoted := false
		for _, vote := range listVotes.Result {
			if vote.User.ID == c.User().ID {
				hasVoted = true
				break
			}
		}

		if hasVoted {
			err := bus.Dispatch(c, &cmd.RemoveVote{Post: getPost.Result, User: c.User()})
			if err != nil {
				return c.Failure(err)
			}
			return c.Ok(web.Map{"voted": false})
		}

		err = bus.Dispatch(c, &cmd.AddVote{Post: getPost.Result, User: c.User()})
		if err != nil {
			return c.Failure(err)
		}
		metrics.TotalVotes.Inc()
		return c.Ok(web.Map{"voted": true})
	}
}

// Subscribe adds current user to list of subscribers of given post
func Subscribe() web.HandlerFunc {
	return func(c *web.Context) error {
		return addOrRemove(c, func(post *entity.Post, user *entity.User) bus.Msg {
			return &cmd.AddSubscriber{Post: post, User: user}
		})
	}
}

// Unsubscribe removes current user from list of subscribers of given post
func Unsubscribe() web.HandlerFunc {
	return func(c *web.Context) error {
		return addOrRemove(c, func(post *entity.Post, user *entity.User) bus.Msg {
			return &cmd.RemoveSubscriber{Post: post, User: user}
		})
	}
}

// ListVotes returns a list of all votes on given post
func ListVotes() web.HandlerFunc {
	return func(c *web.Context) error {
		number, err := c.ParamAsInt("number")
		if err != nil {
			return c.NotFound()
		}

		getPost := &query.GetPostByNumber{Number: number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		listVotes := &query.ListPostVotes{PostID: getPost.Result.ID, IncludeEmail: true}
		if err := bus.Dispatch(c, listVotes); err != nil {
			return c.Failure(err)
		}

		return c.Ok(listVotes.Result)
	}
}

// AddVoteOnBehalf allows administrators to add votes on behalf of users
func AddVoteOnBehalf() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.AddVoteOnBehalf)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		// Get the post
		getPost := &query.GetPostByNumber{Number: action.Number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		// Check if post can be voted on
		if !getPost.Result.CanBeVoted() {
			return c.BadRequest(web.Map{
				"errors": []web.Map{
					{"message": "This post cannot receive votes."},
				},
			})
		}

		// Find or create user
		var user *entity.User
		getByEmail := &query.GetUserByEmail{Email: action.Email}
		err := bus.Dispatch(c, getByEmail)

		if err != nil && errors.Cause(err) == app.ErrNotFound {
			// Create new user
			user = &entity.User{
				Tenant: c.Tenant(),
				Name:   action.Name,
				Email:  action.Email,
				Role:   enum.RoleVisitor,
			}
			if err := bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
				return c.Failure(err)
			}
		} else if err != nil {
			return c.Failure(err)
		} else {
			user = getByEmail.Result
		}

		// Check if user already voted
		listVotes := &query.ListPostVotes{PostID: getPost.Result.ID}
		if err := bus.Dispatch(c, listVotes); err != nil {
			return c.Failure(err)
		}

		for _, vote := range listVotes.Result {
			if vote.User.ID == user.ID {
				return c.BadRequest(web.Map{
					"errors": []web.Map{
						{"field": "email", "message": "This user has already voted for this post."},
					},
				})
			}
		}

		// Add vote
		if err := bus.Dispatch(c, &cmd.AddVote{Post: getPost.Result, User: user}); err != nil {
			return c.Failure(err)
		}

		metrics.TotalVotes.Inc()
		return c.Ok(web.Map{
			"user": web.Map{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		})
	}
}

// RemoveVoteOnBehalf allows administrators to remove votes on behalf of users
func RemoveVoteOnBehalf() web.HandlerFunc {
	return func(c *web.Context) error {
		action := new(actions.RemoveVoteOnBehalf)
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		// Get the post
		getPost := &query.GetPostByNumber{Number: action.Number}
		if err := bus.Dispatch(c, getPost); err != nil {
			return c.Failure(err)
		}

		// Get the user
		getUser := &query.GetUserByID{UserID: action.UserID}
		if err := bus.Dispatch(c, getUser); err != nil {
			return c.Failure(err)
		}

		// Remove vote
		if err := bus.Dispatch(c, &cmd.RemoveVote{Post: getPost.Result, User: getUser.Result}); err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

func addOrRemove(c *web.Context, getCommand func(post *entity.Post, user *entity.User) bus.Msg) error {
	number, err := c.ParamAsInt("number")
	if err != nil {
		return c.NotFound()
	}

	getPost := &query.GetPostByNumber{Number: number}
	if err := bus.Dispatch(c, getPost); err != nil {
		return c.Failure(err)
	}

	command := getCommand(getPost.Result, c.User())
	if err := bus.Dispatch(c, command); err != nil {
		return c.Failure(err)
	}

	return c.Ok(web.Map{})
}
