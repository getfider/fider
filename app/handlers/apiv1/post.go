package apiv1

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/metrics"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/tasks"
)

// SearchPosts return existing posts based on search criteria
func SearchPosts() web.HandlerFunc {
	return func(c *web.Context) error {
		searchPosts := &query.SearchPosts{
			Query: c.QueryParam("query"),
			View:  c.QueryParam("view"),
			Limit: c.QueryParam("limit"),
			Tags:  c.QueryParamAsArray("tags"),
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

		err := bus.Dispatch(c,
			&cmd.UploadImages{
				Images: action.Attachments,
				Folder: "attachments",
			},
			&cmd.UpdatePost{
				Post:        action.Post,
				Title:       action.Title,
				Description: action.Description,
			},
			&cmd.SetAttachments{
				Post:        action.Post,
				Attachments: action.Attachments,
			},
		)
		if err != nil {
			return c.Failure(err)
		}

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

		if action.Text != "" {
			// Only send notification if user wrote a comment.
			c.Enqueue(tasks.NotifyAboutDeletedPost(action.Post))
		}

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

		return c.Ok(commentByID.Result)
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

		if err := bus.Dispatch(c, &cmd.SetAttachments{
			Post:        getPost.Result,
			Comment:     addNewComment.Result,
			Attachments: action.Attachments,
		}); err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewComment(getPost.Result, action.Content))

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
		err = bus.Dispatch(c, listVotes)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(listVotes.Result)
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

	cmd := getCommand(getPost.Result, c.User())
	err = bus.Dispatch(c, cmd)
	if err != nil {
		return c.Failure(err)
	}

	return c.Ok(web.Map{})
}
