package apiv1

import (
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// UploadImage uploads an image without associating it with a post or comment
func UploadImage() web.HandlerFunc {
	return func(c *web.Context) error {
		input := new(actions.UploadImage)

		if result := c.BindTo(input); !result.Ok {
			return c.HandleValidation(result)
		}

		// Upload the image to the "attachments" folder
		uploadCmd := &cmd.UploadImage{
			Image:  input.Image,
			Folder: "attachments",
		}

		if err := bus.Dispatch(c, uploadCmd); err != nil {
			return c.Failure(err)
		}

		// Return the bkey so it can be used in the markdown
		return c.Ok(web.Map{
			"bkey": input.Image.BlobKey,
		})
	}
}
