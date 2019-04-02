package apiv1

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
	"github.com/getfider/fider/app/tasks"
)

// SearchPosts return existing posts based on search criteria
func SearchPosts() web.HandlerFunc {
	return func(c *web.Context) error {
		posts, err := c.Services().Posts.Search(
			c.QueryParam("query"),
			c.QueryParam("view"),
			c.QueryParam("limit"),
			c.QueryParamAsArray("tags"),
		)
		if err != nil {
			return c.Failure(err)
		}
		return c.Ok(posts)
	}
}

// CreatePost creates a new post on current tenant
func CreatePost() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.CreateNewPost)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		posts := c.Services().Posts

		err := webutil.ProcessMultiImageUpload(c, input.Model.Attachments, "attachments")
		if err != nil {
			return c.Failure(err)
		}

		post, err := posts.Add(input.Model.Title, input.Model.Description)
		if err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.SetAttachments{Post: post, Attachments: input.Model.Attachments})
		if err != nil {
			return c.Failure(err)
		}

		if err := bus.Dispatch(c, &cmd.AddVote{Post: post, User: c.User()}); err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewPost(post))

		return c.Ok(web.Map{
			"id":     post.ID,
			"number": post.Number,
			"title":  post.Title,
			"slug":   post.Slug,
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

		post, err := c.Services().Posts.GetByNumber(number)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(post)
	}
}

// UpdatePost updates an existing post of current tenant
func UpdatePost() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.UpdatePost)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := webutil.ProcessMultiImageUpload(c, input.Model.Attachments, "attachments")
		if err != nil {
			return c.Failure(err)
		}

		_, err = c.Services().Posts.Update(input.Post, input.Model.Title, input.Model.Description)
		if err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.SetAttachments{Post: input.Post, Attachments: input.Model.Attachments})
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// SetResponse changes current post staff response
func SetResponse() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.SetResponse)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		post, err := c.Services().Posts.GetByNumber(input.Model.Number)
		if err != nil {
			return c.Failure(err)
		}

		prevStatus := post.Status
		if input.Model.Status == models.PostDuplicate {
			err = bus.Dispatch(c, &cmd.MarkPostAsDuplicate{Post: post, Original: input.Original})
		} else {
			err = bus.Dispatch(c, &cmd.SetPostResponse{
				Post:   post,
				Text:   input.Model.Text,
				Status: input.Model.Status,
			})
		}

		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutStatusChange(post, prevStatus))

		return c.Ok(web.Map{})
	}
}

// DeletePost deletes an existing post of current tenant
func DeletePost() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.DeletePost)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.SetPostResponse{
			Post:   input.Post,
			Text:   input.Model.Text,
			Status: models.PostDeleted,
		})
		if err != nil {
			return c.Failure(err)
		}

		if input.Model.Text != "" {
			// Only send notification if user wrote a comment.
			c.Enqueue(tasks.NotifyAboutDeletedPost(input.Post))
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

		post, err := c.Services().Posts.GetByNumber(number)
		if err != nil {
			return c.Failure(err)
		}

		getComments := &query.GetCommentsByPost{Post: post}
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
		if bus.Dispatch(c, commentByID); err != nil {
			return c.Failure(err)
		}

		return c.Ok(commentByID.Result)
	}
}

// PostComment creates a new comment on given post
func PostComment() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.AddNewComment)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := webutil.ProcessMultiImageUpload(c, input.Model.Attachments, "attachments")
		if err != nil {
			return c.Failure(err)
		}

		post, err := c.Services().Posts.GetByNumber(input.Model.Number)
		if err != nil {
			return c.Failure(err)
		}

		addNewComment := &cmd.AddNewComment{
			Post:    post,
			Content: input.Model.Content,
		}
		err = bus.Dispatch(c, addNewComment)
		if err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.SetAttachments{
			Post:        post,
			Comment:     addNewComment.Result,
			Attachments: input.Model.Attachments,
		})
		if err != nil {
			return c.Failure(err)
		}

		c.Enqueue(tasks.NotifyAboutNewComment(post, input.Model))

		return c.Ok(web.Map{
			"id": addNewComment.Result.ID,
		})
	}
}

// UpdateComment changes an existing comment with new content
func UpdateComment() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.EditComment)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := webutil.ProcessMultiImageUpload(c, input.Model.Attachments, "attachments")
		if err != nil {
			return c.Failure(err)
		}

		updateComment := &cmd.UpdateComment{
			CommentID: input.Model.ID,
			Content:   input.Model.Content,
		}

		setAttachments := &cmd.SetAttachments{
			Post:        input.Post,
			Comment:     input.Comment,
			Attachments: input.Model.Attachments,
		}

		err = bus.Dispatch(c, updateComment, setAttachments)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(web.Map{})
	}
}

// DeleteComment deletes an existing comment by its ID
func DeleteComment() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.DeleteComment)
		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		err := bus.Dispatch(c, &cmd.DeleteComment{
			CommentID: input.Model.CommentID,
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
		return addOrRemove(c, func(post *models.Post, user *models.User) bus.Msg {
			return &cmd.AddVote{Post: post, User: user}
		})
	}
}

// RemoveVote removes current user from given post list of votes
func RemoveVote() web.HandlerFunc {
	return func(c *web.Context) error {
		return addOrRemove(c, func(post *models.Post, user *models.User) bus.Msg {
			return &cmd.RemoveVote{Post: post, User: user}
		})
	}
}

// Subscribe adds current user to list of subscribers of given post
func Subscribe() web.HandlerFunc {
	return func(c *web.Context) error {
		return addOrRemove(c, func(post *models.Post, user *models.User) bus.Msg {
			return &cmd.AddSubscriber{Post: post, User: user}
		})
	}
}

// Unsubscribe removes current user from list of subscribers of given post
func Unsubscribe() web.HandlerFunc {
	return func(c *web.Context) error {
		return addOrRemove(c, func(post *models.Post, user *models.User) bus.Msg {
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

		post, err := c.Services().Posts.GetByNumber(number)
		if err != nil {
			return c.Failure(err)
		}

		listVotes := &query.ListPostVotes{PostID: post.ID}
		err = bus.Dispatch(c, listVotes)
		if err != nil {
			return c.Failure(err)
		}

		return c.Ok(listVotes.Result)
	}
}

func addOrRemove(c *web.Context, getCommand func(post *models.Post, user *models.User) bus.Msg) error {
	number, err := c.ParamAsInt("number")
	if err != nil {
		return c.NotFound()
	}

	post, err := c.Services().Posts.GetByNumber(number)
	if err != nil {
		return c.Failure(err)
	}

	cmd := getCommand(post, c.User())
	err = bus.Dispatch(c, cmd)
	if err != nil {
		return c.Failure(err)
	}

	return c.Ok(web.Map{})
}
