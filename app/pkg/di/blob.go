package di

import (
	"strings"

	"github.com/getfider/fider/app/pkg/blob"
	"github.com/getfider/fider/app/pkg/blob/fs"
	"github.com/getfider/fider/app/pkg/blob/s3"
	"github.com/getfider/fider/app/pkg/blob/sql"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/email/mailgun"
	"github.com/getfider/fider/app/pkg/email/noop"
	"github.com/getfider/fider/app/pkg/email/smtp"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web/http"
)

// NewEmailer creates a new emailer instance based on current configuration
func NewEmailer(logger log.Logger) email.Sender {
	if env.IsTest() {
		return noop.NewSender()
	}

	if env.Config.Email.Mailgun.APIKey != "" {
		return mailgun.NewSender(
			logger,
			http.NewClient(),
			env.Config.Email.Mailgun.Domain,
			env.Config.Email.Mailgun.APIKey,
		)
	}

	return smtp.NewSender(
		logger,
		env.Config.Email.SMTP.Host,
		env.Config.Email.SMTP.Port,
		env.Config.Email.SMTP.Username,
		env.Config.Email.SMTP.Password,
	)
}

// NewBlobStorage creates a new blob storage instance based on current configuration
func NewBlobStorage(db *dbx.Database) blob.Storage {
	storageType := strings.ToLower(env.Config.BlobStorage.Type)
	switch storageType {
	case "sql":
		return sql.NewStorage(db)
	case "s3":
		return s3.NewStorage(env.Config.BlobStorage.S3.BucketName)
	case "fs":
		return fs.NewStorage(env.Config.BlobStorage.FS.Path)
	}
	panic("Invalid blob storage type: " + storageType)
}
