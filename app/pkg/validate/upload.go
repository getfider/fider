package validate

import (
	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/img"
)

// ImageUploadOpts arguments to validate given upload
type ImageUploadOpts struct {
	IsRequired   bool
	MinWidth     int
	MinHeight    int
	MaxKilobytes int
}

//ImageUpload validates given image upload
func ImageUpload(upload *models.ImageUpload, opts ImageUploadOpts) ([]string, error) {
	messages := []string{}

	if opts.IsRequired {
		if upload == nil || (upload.BlobKey == "" && upload.Upload == nil) || upload.Remove {
			messages = append(messages, "An image is required.")
		}
	}

	if upload != nil && upload.Upload != nil && len(upload.Upload.Content) > 0 {
		logo, err := img.Parse(upload.Upload.Content)
		if err != nil {
			if err == img.ErrNotSupported {
				messages = append(messages, "This file format not supported.")
			} else {
				return nil, err
			}
		} else {
			if logo.Width < opts.MinWidth || logo.Height < opts.MinHeight {
				messages = append(messages, fmt.Sprintf("The image must have minimum dimensions of %dx%d pixels.", opts.MinWidth, opts.MinHeight))
			}

			if logo.Width != logo.Height {
				messages = append(messages, "The image must have an aspect ratio of 1:1.")
			}

			if logo.Size > (opts.MaxKilobytes * 1024) {
				messages = append(messages, fmt.Sprintf("The image size must be smaller than %dKB.", opts.MaxKilobytes))
			}
		}
	}

	return messages, nil
}
