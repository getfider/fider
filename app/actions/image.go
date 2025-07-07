package actions

import (
	"context"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/validate"
)

type DeleteImage struct {
	BlobKey string `route:"bkey"`
}

func (input *DeleteImage) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil
}

func (input *DeleteImage) Validate(ctx context.Context, user *entity.User) *validate.Result {
	return validate.Success()
}

// UploadImage is used to upload an image without associating it with a post or comment
type UploadImage struct {
	Image *dto.ImageUpload `json:"image"`
}

// Initialize the model
func (input *UploadImage) Initialize() {}

// IsAuthorized returns true if current user is authorized to perform this action
func (input *UploadImage) IsAuthorized(ctx context.Context, user *entity.User) bool {
	return user != nil
}

// Validate if current model is valid
func (input *UploadImage) Validate(ctx context.Context, user *entity.User) *validate.Result {
	result := validate.Success()

	messages, err := validate.ImageUpload(ctx, input.Image, validate.ImageUploadOpts{
		IsRequired:   true,
		MinHeight:    50,
		MinWidth:     50,
		ExactRatio:   false,
		MaxKilobytes: 5120,
	})

	if err != nil {
		return validate.Error(err)
	}
	result.AddFieldFailure("image", messages...)
	return result
}
