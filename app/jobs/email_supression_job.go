package jobs

import (
	"time"

	"github.com/Spicy-Bush/fider-tarkov-community/app/models/cmd"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/dto"
	"github.com/Spicy-Bush/fider-tarkov-community/app/models/query"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/bus"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/errors"
	"github.com/Spicy-Bush/fider-tarkov-community/app/pkg/log"
)

type EmailSupressionJobHandler struct {
}

func (e EmailSupressionJobHandler) Schedule() string {
	return "0 5 * * * *" // every hour at minute 5
}

func (e EmailSupressionJobHandler) Run(ctx Context) error {
	startTime := ctx.LastSuccessfulRun
	if startTime == nil {
		twoDaysAgo := time.Now().AddDate(0, 0, -2)
		startTime = &twoDaysAgo
	}

	q := &query.FetchRecentSupressions{
		StartTime: *startTime,
	}

	if err := bus.Dispatch(ctx, q); err != nil {
		return errors.Wrap(err, "failed to fetch recent supressions")
	}

	c := &cmd.SupressEmail{
		EmailAddresses: q.EmailAddresses,
	}
	if err := bus.Dispatch(ctx, c); err != nil {
		return errors.Wrap(err, "failed to supress emails")
	}

	log.Debugf(ctx, "@{Count} account(s) marked with supressed email", dto.Props{
		"Count": c.NumOfSupressedEmailAddresses,
	})

	return nil
}
