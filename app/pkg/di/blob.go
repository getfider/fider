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
	"github.com/getfider/fider/app/pkg/web"
)

// NewEmailer creates a new emailer instance based on current configuration
func NewEmailer(logger log.Logger) email.Sender {
	if env.IsTest() {
		return noop.NewSender()
	}
	if env.IsDefined("EMAIL_MAILGUN_API") {
		return mailgun.NewSender(logger, web.NewHTTPClient(), env.MustGet("EMAIL_MAILGUN_DOMAIN"), env.MustGet("EMAIL_MAILGUN_API"))
	}
	return smtp.NewSender(
		logger,
		env.MustGet("EMAIL_SMTP_HOST"),
		env.MustGet("EMAIL_SMTP_PORT"),
		env.GetEnvOrDefault("EMAIL_SMTP_USERNAME", ""),
		env.GetEnvOrDefault("EMAIL_SMTP_PASSWORD", ""),
	)
}

// NewBlobStorage creates a new blob storage instance based on current configuration
func NewBlobStorage(trx *dbx.Trx) blob.Storage {
	storageType := strings.ToLower(env.GetEnvOrDefault("BLOB_STORAGE", "sql"))
	switch storageType {
	case "sql":
		return sql.NewStorage(trx)
	case "s3":
		return s3.NewStorage(env.MustGet("BLOB_STORAGE_S3_BUCKET"))
	case "fs":
		return fs.NewStorage(env.MustGet("BLOB_STORAGE_FS_PATH"))
	}
	panic("Invalid blob storage type: " + storageType)
}
