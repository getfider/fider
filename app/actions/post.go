package actions

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/i18n"
	"github.com/gosimple/slug"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
)

// CreateNewPost is used to create a new post
type CreateNewPost struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	TagSlugs    []string           `json:"tags"`
	Attachments []*dto.ImageUpload `json:"attachments"`

	Tags []*entity.Tag
}

// OnPreExecute prefetches Tags for later use
func (input *CreateNewPost) OnPreExecute(ctx context.Context) error {
	if env.Config.PostCreationWithTagsEnabled {
		input.Tags = make([]*entity.Tag, 0, len(input.TagSlugs))
		for _, slug := range input.TagSlugs {
			getTag := &query.GetTagBySlug{Slug: slug}
			if err := bus.Dispatch(ctx, getTag); err != nil {
				break
			}

			input.Tags = append(input.Tags, getTag.Result)
		}
	}

	return nil
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *CreateNewPost) IsAuthorized(ctx context.Context, user *entity.User) bool {

	if user == nil {
		return false
	} else if env.Config.PostCreationWithTagsEnabled && !user.IsCollaborator() {
		for _, tag := range action.Tags {
			if !tag.IsPublic {
				return false
			}
		}
	}
	return true
}

// Validate if current model is valid
func (action *CreateNewPost) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	re := regexp.MustCompile(`\s+`)
	normalizedTitle := strings.TrimSpace(re.ReplaceAllString(action.Title, " "))

	if normalizedTitle == "" {
		result.AddFieldFailure("title", propertyIsRequired(ctx, "title"))
	} else if len(normalizedTitle) < 10 {
		result.AddFieldFailure("title", i18n.T(ctx, "validation.custom.descriptivetitle"))
	} else if len(normalizedTitle) > 100 {
		result.AddFieldFailure("title", propertyMaxStringLen(ctx, "title", 100))
	} else if env.Config.PostCreationWithTagsEnabled && len(action.TagSlugs) != len(action.Tags) {
		result.AddFieldFailure("tags", propertyIsInvalid(ctx, "tags"))
	} else {
		err := bus.Dispatch(ctx, &query.GetPostBySlug{Slug: slug.Make(action.Title)})
		if err != nil && errors.Cause(err) != app.ErrNotFound {
			return validate.Error(err)
		} else if err == nil {
			result.AddFieldFailure("title", i18n.T(ctx, "validation.custom.duplicatetitle"))
		}
	}

	messages, err := validate.MultiImageUpload(ctx, nil, action.Attachments, validate.MultiImageUploadOpts{
		MaxUploads:   3,
		MaxKilobytes: 5120,
		ExactRatio:   false,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("attachments", messages...)

	return result
}

// UpdatePost is used to edit an existing new post
type UpdatePost struct {
	Number      int                `route:"number"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Attachments []*dto.ImageUpload `json:"attachments"`

	Post *entity.Post
}

// OnPreExecute prefetches Post for later use
func (input *UpdatePost) OnPreExecute(ctx context.Context) error {
	getPost := &query.GetPostByNumber{Number: input.Number}
	if err := bus.Dispatch(ctx, getPost); err != nil {
		return err
	}

	input.Post = getPost.Result
	return nil
}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UpdatePost) IsAuthorized(ctx context.Context, user *entity.User) bool {
	if user.IsCollaborator() {
		return true
	}

	timeAgo := time.Now().UTC().Sub(input.Post.CreatedAt)
	return input.Post.User.ID == user.ID && timeAgo <= 1*time.Hour
}

// Validate if current model is valid
func (action *UpdatePost) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Title == "" {
		result.AddFieldFailure("title", propertyIsRequired(ctx, "title"))
	} else if len(action.Title) < 10 {
		result.AddFieldFailure("title", i18n.T(ctx, "validation.custom.descriptivetitle"))
	} else if len(action.Title) > 100 {
		result.AddFieldFailure("title", propertyMaxStringLen(ctx, "title", 100))
	}

	postBySlug := &query.GetPostBySlug{Slug: slug.Make(action.Title)}
	err := bus.Dispatch(ctx, postBySlug)
	if err != nil && errors.Cause(err) != app.ErrNotFound {
		return validate.Error(err)
	} else if err == nil && postBySlug.Result.ID != action.Post.ID {
		result.AddFieldFailure("title", i18n.T(ctx, "validation.custom.duplicatetitle"))
	}

	if len(action.Attachments) > 0 {
		getAttachments := &query.GetAttachments{Post: action.Post}
		err = bus.Dispatch(ctx, getAttachments)
		if err != nil {
			return validate.Error(err)
		}

		messages, err := validate.MultiImageUpload(ctx, getAttachments.Result, action.Attachments, validate.MultiImageUploadOpts{
			MaxUploads:   3,
			MaxKilobytes: 5120,
			ExactRatio:   false,
		})
		if err != nil {
			return validate.Error(err)
		}
		result.AddFieldFailure("attachments", messages...)
	}

	return result
}

type ToggleCommentReaction struct {
	Number   int    `route:"number"`
	Comment  int    `route:"id"`
	Reaction string `route:"reaction"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *ToggleCommentReaction) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil
}

// Validate if current model is valid
func (action *ToggleCommentReaction) Validate(ctx context.Context, user *entity.User) *validate.Result {

	result := validate.Success()

	allowedEmojis := []string{"👍", "👎", "😄", "🎉", "😕", "❤️", "🚀", "👀"}
	isAllowed := false
	for _, emoji := range allowedEmojis {
		if action.Reaction == emoji {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		result.AddFieldFailure("reaction", i18n.T(ctx, "validation.custom.invalidemoji"))
	}

	return result
}

// AddNewComment represents a new comment to be added
type AddNewComment struct {
	Number      int                `route:"number"`
	Content     string             `json:"content"`
	Attachments []*dto.ImageUpload `json:"attachments"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *AddNewComment) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil
}

// Validate if current model is valid
func (action *AddNewComment) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Content == "" {
		result.AddFieldFailure("content", propertyIsRequired(ctx, "comment"))
	}

	messages, err := validate.MultiImageUpload(ctx, nil, action.Attachments, validate.MultiImageUploadOpts{
		MaxUploads:   2,
		MaxKilobytes: 5120,
		ExactRatio:   false,
	})
	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("attachments", messages...)

	return result
}

// SetResponse represents the action to update an post response
type SetResponse struct {
	Number         int             `route:"number"`
	Status         enum.PostStatus `json:"status"`
	Text           string          `json:"text"`
	OriginalNumber int             `json:"originalNumber"`

	Original *entity.Post
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *SetResponse) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsCollaborator()
}

// Validate if current model is valid
func (action *SetResponse) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Status < enum.PostOpen || action.Status > enum.PostDuplicate {
		result.AddFieldFailure("status", propertyIsInvalid(ctx, "status"))
	}

	if action.Status == enum.PostDuplicate {
		if action.OriginalNumber == action.Number {
			result.AddFieldFailure("originalNumber", i18n.T(ctx, "validation.custom.selfduplicate"))
		}

		getOriginaPost := &query.GetPostByNumber{Number: action.OriginalNumber}
		err := bus.Dispatch(ctx, getOriginaPost)
		if err != nil {
			if errors.Cause(err) == app.ErrNotFound {
				result.AddFieldFailure("originalNumber", i18n.T(ctx, "validation.custom.originalpostnotfound"))
			} else {
				return validate.Error(err)
			}
		}

		if getOriginaPost.Result != nil {
			action.Original = getOriginaPost.Result
		}
	}

	return result
}

// DeletePost represents the action of an administrator deleting an existing Post
type DeletePost struct {
	Number int    `route:"number"`
	Text   string `json:"text"`

	Post *entity.Post
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *DeletePost) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil && user.IsAdministrator()
}

// Validate if current model is valid
func (action *DeletePost) Validate(ctx context.Context, user *entity.User) *validate.Result {
	getPost := &query.GetPostByNumber{Number: action.Number}
	if err := bus.Dispatch(ctx, getPost); err != nil {
		return validate.Error(err)
	}

	action.Post = getPost.Result

	isReferencedQuery := &query.PostIsReferenced{PostID: action.Post.ID}
	if err := bus.Dispatch(ctx, isReferencedQuery); err != nil {
		return validate.Error(err)
	}

	if isReferencedQuery.Result {
		return validate.Failed(i18n.T(ctx, "validation.custom.cannotdeleteduplicatepost"))
	}

	return validate.Success()
}

// EditComment represents the action to update an existing comment
type EditComment struct {
	PostNumber  int                `route:"number"`
	ID          int                `route:"id"`
	Content     string             `json:"content"`
	Attachments []*dto.ImageUpload `json:"attachments"`

	Post    *entity.Post
	Comment *entity.Comment
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *EditComment) IsAuthorized(ctx context.Context, user *entity.User) bool {
	postByNumber := &query.GetPostByNumber{Number: action.PostNumber}
	commentByID := &query.GetCommentByID{CommentID: action.ID}
	if err := bus.Dispatch(ctx, postByNumber, commentByID); err != nil {
		return false
	}

	action.Post = postByNumber.Result
	action.Comment = commentByID.Result
	return user.ID == action.Comment.User.ID || user.IsCollaborator()
}

// Validate if current model is valid
func (action *EditComment) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	if action.Content == "" {
		result.AddFieldFailure("content", propertyIsRequired(ctx, "comment"))
	}

	if len(action.Attachments) > 0 {

		nonRemovedCount := 0
		for _, v := range action.Attachments {
			if !v.Remove {
				nonRemovedCount++
			}
		}

		if nonRemovedCount > 2 {
			result.AddFieldFailure("content", i18n.T(ctx, "validation.custom.maxattachments", i18n.Params{"number": 2}))
		}

	}

	return result
}

// DeleteComment represents the action of deleting an existing comment
type DeleteComment struct {
	PostNumber int `route:"number"`
	CommentID  int `route:"id"`
}

// IsAuthorized returns true if current user is authorized to perform this action
func (action *DeleteComment) IsAuthorized(ctx context.Context, user *entity.User) bool {
	commentByID := &query.GetCommentByID{CommentID: action.CommentID}
	if err := bus.Dispatch(ctx, commentByID); err != nil {
		return false
	}

	return user.ID == commentByID.Result.User.ID || user.IsCollaborator()
}

// Validate if current model is valid
func (action *DeleteComment) Validate(ctx context.Context, user *entity.User) *validate.Result {
	return validate.Success()
}
