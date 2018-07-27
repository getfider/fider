package validate

import (
	"fmt"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/img"
)

//ImageUpload validates given image upload
func ImageUpload(upload *models.ImageUpload, minWidth, minHeight, maxKilobytes int) ([]string, error) {
	messages := []string{}

	if upload != nil && upload.Upload != nil && len(upload.Upload.Content) > 0 {
		logo, err := img.Parse(upload.Upload.Content)
		if err != nil {
			if err == img.ErrNotSupported {
				messages = append(messages, "This file format not supported.")
			} else {
				return nil, err
			}
		} else {
			if logo.Width < minWidth || logo.Height < minHeight {
				messages = append(messages, fmt.Sprintf("The image must have minimum dimensions of %dx%d pixels.", minWidth, minHeight))
			}

			if logo.Width != logo.Height {
				messages = append(messages, "The image must have an aspect ratio of 1:1.")
			}

			if logo.Size > (maxKilobytes * 1024) {
				messages = append(messages, fmt.Sprintf("The image size must be smaller than %dKB.", maxKilobytes))
			}
		}
	}

	return messages, nil
}
