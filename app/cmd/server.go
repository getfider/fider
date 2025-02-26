package cmd

import (
	"context"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/Spicy-Bush/fider-tarkov-community/app/jobs"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/query"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/env"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/errors"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/log"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/web"
	"github.com/robfig/cron"

	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/billing/paddle"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/blob/fs"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/blob/s3"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/blob/sql"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/email/awsses"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/email/mailgun"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/email/smtp"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/httpclient"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/log/console"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/log/file"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/log/sql"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/oauth"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/sqlstore/postgres"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/userlist"
	_ "github.com/Spicy-Bush/fider-tarkov-community/app/services/webhook"
)

// RunServer starts the Fider Server
// Returns an exitcode, 0 for OK and 1 for ERROR
func RunServer() int {
	svcs := bus.Init()
	ctx := log.WithProperty(context.Background(), log.PropertyKeyTag, "BOOTSTRAP")
	for _, s := range svcs {
		log.Debugf(ctx, "Service '@{ServiceCategory}.@{ServiceName}' has been initialized.", dto.Props{
			"ServiceCategory": s.Category(),
			"ServiceName":     s.Name(),
		})
	}

	copyEtcFiles(ctx)
	startJobs(ctx)

	e := routes(web.New())
	go e.Start(":" + env.Config.Port)
	return listenSignals(e)
}

// Starts all scheduled jobs
func startJobs(ctx context.Context) {
	c := cron.New()
	_ = c.AddJob(jobs.NewJob(ctx, "PurgeExpiredNotificationsJob", jobs.PurgeExpiredNotificationsJobHandler{}))
	_ = c.AddJob(jobs.NewJob(ctx, "EmailSupressionJob", jobs.EmailSupressionJobHandler{}))

	if env.IsBillingEnabled() {
		_ = c.AddJob(jobs.NewJob(ctx, "LockExpiredTenantsJob", jobs.LockExpiredTenantsJobHandler{}))
	}

	c.Start()
}

// on startup, copy all etc/ files from configured blob storage into local etc/ folder
// this can be used to avoid having to mount volumes on ephemeral environments
func copyEtcFiles(ctx context.Context) {
	q := &query.ListBlobs{Prefix: "etc/"}
	if err := bus.Dispatch(ctx, q); err != nil {
		panic(errors.Wrap(err, "failed to list etc/ blobs"))
	}

	if len(q.Result) == 0 {
		log.Debug(ctx, "No etc/ files to copy")
		return
	}

	for _, blobKey := range q.Result {
		getBlob := &query.GetBlobByKey{Key: blobKey}
		if err := bus.Dispatch(ctx, getBlob); err != nil {
			panic(errors.Wrap(err, "failed to get blob by key '%s'", blobKey))
		}
		if err := os.MkdirAll(path.Dir(blobKey), 0774); err != nil {
			panic(errors.Wrap(err, "failed to create dir"))
		}
		if err := os.WriteFile(blobKey, getBlob.Result.Content, 0774); err != nil {
			panic(errors.Wrap(err, "failed to write blob '%s'", blobKey))
		}

		log.Debugf(ctx, "Copied '@{BlobKey}' to etc/ folder.", dto.Props{
			"BlobKey": blobKey,
		})
	}
}

func listenSignals(e *web.Engine) int {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, append([]os.Signal{syscall.SIGTERM, syscall.SIGINT}, extraSignals...)...)
	for {
		s := <-signals
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			err := e.Stop()
			if err != nil {
				return 1
			}
			return 0
		default:
			ret := handleExtraSignal(s, e)
			if ret >= 0 {
				return ret
			}
		}
	}
}
