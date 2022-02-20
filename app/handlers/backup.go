package handlers

import (
	"github.com/getfider/fider/app/pkg/backup"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
)

// ExportBackupZip returns a Zip file with all content
func ExportBackupZip() web.HandlerFunc {
	return func(c *web.Context) error {

		file, err := backup.Create(c)
		if err != nil {
			log.Error(c, errors.Wrap(err, "failed to create backup"))
			return c.Failure(err)
		}

		return c.Attachment("backup.zip", "application/zip", file.Bytes())
	}
}
