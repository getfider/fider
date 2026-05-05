package jobs

import (
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

type EmailSuppressionJobHandler struct {
}

func (e EmailSuppressionJobHandler) Schedule() string {
	return "0 5 * * * *" // every hour at minute 5
}

func (e EmailSuppressionJobHandler) Run(ctx Context) error {
	startTime := ctx.LastSuccessfulRun
	if startTime == nil {
		twoDaysAgo := time.Now().AddDate(0, 0, -2)
		startTime = &twoDaysAgo
	}

	q := &query.FetchRecentSuppressions{
		StartTime: *startTime,
	}

	if err := bus.Dispatch(ctx, q); err != nil {
		return errors.Wrap(err, "failed to fetch recent suppressions")
	}

	c := &cmd.SuppressEmail{
		EmailAddresses: q.EmailAddresses,
	}
	if err := bus.Dispatch(ctx, c); err != nil {
		return errors.Wrap(err, "failed to suppress emails")
	}

	log.Debugf(ctx, "@{Count} account(s) marked with suppressed email", dto.Props{
		"Count": c.NumOfSuppressedEmailAddresses,
	})

	return nil
}
