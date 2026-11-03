package validate

import (
	"context"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/i18n"
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

// MultiImageUpload validates multiple image uploads
func MultiImageUpload(ctx context.Context, currentAttachments []string, uploads []*dto.ImageUpload, opts MultiImageUploadOpts) ([]string, error) {
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

		messages, err := ImageUpload(ctx, upload, ImageUploadOpts{
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
		return []string{i18n.T(ctx, "validation.custom.maxattachments", i18n.Params{"number": opts.MaxUploads})}, nil
	}

	return []string{}, nil
}

// ImageUpload validates given image upload
func ImageUpload(ctx context.Context, upload *dto.ImageUpload, opts ImageUploadOpts) ([]string, error) {
	messages := []string{}

	if opts.IsRequired {
		if upload == nil || (upload.BlobKey == "" && upload.Upload == nil) || upload.Remove {
			messages = append(messages, i18n.T(ctx, "validation.required",
				i18n.Params{"name": i18n.T(ctx, "property.image")},
			))
		}
	}

	if upload != nil && upload.Upload != nil && len(upload.Upload.Content) > 0 {
		// Only reads the image header, so decompression bombs are rejected before anything
		// decodes the pixel data. Budgeted for one operation: the resize below.
		header, err := CheckDecodeBudget(upload.Upload.Content, 1)
		switch err {
		case imagic.ErrNotSupported:
			messages = append(messages, i18n.T(ctx, "validation.custom.unsupportedfileformat"))
		case ErrImageTooLarge:
			messages = append(messages, i18n.T(ctx, "validation.custom.maximagepixels"))
			return messages, nil
		case nil:
			if header.Width < opts.MinWidth || header.Height < opts.MinHeight {
				messages = append(messages, i18n.T(ctx, "validation.custom.minimagedimensions",
					i18n.Params{"width": opts.MinWidth, "height": opts.MinHeight},
				))
			}

			if opts.ExactRatio && header.Width != header.Height {
				messages = append(messages, i18n.T(ctx, "validation.custom.imagesquareratio"))
			}

			if len(upload.Upload.Content) > (opts.MaxKilobytes * 1024) {
				messages = append(messages, i18n.T(ctx, "validation.custom.maximagesize",
					i18n.Params{"kilobytes": opts.MaxKilobytes},
				))
			}

			// The upload is rejected anyway, so don't spend time/memory decoding it.
			if len(messages) > 0 {
				return messages, nil
			}

			if header.Height > MaxDimensionSize && header.Width > MaxDimensionSize {
				// The decode budget was already checked above
				newImageBytes, err := decodeWithinBudget(ctx, upload.Upload.Content, imagic.Resize(MaxDimensionSize))
				if err == ErrDecodeBusy {
					messages = append(messages, i18n.T(ctx, "validation.custom.imageprocessingbusy"))
					return messages, nil
				}
				if err != nil {
					return nil, err
				}
				upload.Upload.Content = newImageBytes
			}
		default:
			return nil, err
		}
	}

	return messages, nil
}
