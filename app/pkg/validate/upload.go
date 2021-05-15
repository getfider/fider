package validate

import (
	"fmt"

	"github.com/getfider/fider/app/models/dto"
	"github.com/goenning/imagic"
)

// MaxDimensionSize is the max width/height of an image. If image is bigger than this, it'll be resized.
const MaxDimensionSize = 1500

// MultiImageUploadOpts arguments to validate mulitple image upload process
type MultiImageUploadOpts struct {
	MaxUploads   int
	IsRequired   bool
	MinWidth     int
	MinHeight    int
	ExactRatio   bool
	MaxKilobytes int
}

// ImageUploadOpts arguments to validate given upload
type ImageUploadOpts struct {
	IsRequired   bool
	MinWidth     int
	MinHeight    int
	ExactRatio   bool
	MaxKilobytes int
}

//MultiImageUpload validates multiple image uploads
func MultiImageUpload(currentAttachments []string, uploads []*dto.ImageUpload, opts MultiImageUploadOpts) ([]string, error) {
	if currentAttachments == nil {
		currentAttachments = []string{}
	}

	totalCount := len(currentAttachments)

	for _, upload := range uploads {
		if upload.Remove {
			for _, attachment := range currentAttachments {
				if attachment == upload.BlobKey {
					totalCount--
				}
			}
		} else if upload.Upload != nil {
			totalCount++
		}

		messages, err := ImageUpload(upload, ImageUploadOpts{
			IsRequired:   opts.IsRequired,
			MinWidth:     opts.MinWidth,
			MinHeight:    opts.MinHeight,
			ExactRatio:   opts.ExactRatio,
			MaxKilobytes: opts.MaxKilobytes,
		})
		if err != nil {
			return nil, err
		}
		if len(messages) > 0 {
			return messages, nil
		}
	}

	if totalCount > opts.MaxUploads {
		return []string{fmt.Sprintf("A maximum of %d attachments are allowed per post.", opts.MaxUploads)}, nil
	}

	return []string{}, nil
}

//ImageUpload validates given image upload
func ImageUpload(upload *dto.ImageUpload, opts ImageUploadOpts) ([]string, error) {
	messages := []string{}

	if opts.IsRequired {
		if upload == nil || (upload.BlobKey == "" && upload.Upload == nil) || upload.Remove {
			messages = append(messages, "An image is required.")
		}
	}

	if upload != nil && upload.Upload != nil && len(upload.Upload.Content) > 0 {
		logo, err := imagic.Parse(upload.Upload.Content)
		if err != nil {
			if err == imagic.ErrNotSupported {
				messages = append(messages, "This file format not supported.")
			} else {
				return nil, err
			}
		} else {

			if logo.Width < opts.MinWidth || logo.Height < opts.MinHeight {
				messages = append(messages, fmt.Sprintf("The image must have minimum dimensions of %dx%d pixels.", opts.MinWidth, opts.MinHeight))
			}

			if opts.ExactRatio && logo.Width != logo.Height {
				messages = append(messages, "The image must have an aspect ratio of 1:1.")
			}

			if logo.Size > (opts.MaxKilobytes * 1024) {
				messages = append(messages, fmt.Sprintf("The image size must be smaller than %dKB.", opts.MaxKilobytes))
			}

			if logo.Height > MaxDimensionSize && logo.Width > MaxDimensionSize {
				newImageBytes, err := imagic.Apply(upload.Upload.Content, imagic.Resize(MaxDimensionSize))
				if err != nil {
					return nil, err
				}
				upload.Upload.Content = newImageBytes
			}
		}
	}

	return messages, nil
}
