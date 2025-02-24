package resend

import (
	"context"
	"github.com/resend/resend-go/v2"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/services/email"
)

var resendClient *resend.Client

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Resend"
}

func (s Service) Category() string {
	return "email"
}

func (s Service) Enabled() bool {
	return env.Config.Email.Type == "resend"
}

func (s Service) Init() {
	resendEnvConfig := env.Config.Email.Resend
	resendClient = resend.NewClient(resendEnvConfig.APIKey)
	log.Debug(context.Background(), "Resend client initialized.")
	bus.AddListener(sendMail)
	bus.AddHandler(fetchRecentSupressions)
}

func sendMail(ctx context.Context, c *cmd.SendMail) {
	if c.Props == nil {
		c.Props = dto.Props{}
	}

	if c.From.Address == "" {
		c.From.Address = email.NoReply
	}

	for _, to := range c.To {
		if to.Address == "" {
			return
		}

		if !email.CanSendTo(to.Address) {
			log.Warnf(ctx, "Skipping email to '@{Name} <@{Address}>'.", dto.Props{
				"Name":    to.Name,
				"Address": to.Address,
			})
			return
		}

		log.Debugf(ctx, "Sending email to @{Address} with template @{TemplateName} and params @{Props}.", dto.Props{
			"Address":      to.Address,
			"TemplateName": c.TemplateName,
			"Props":        to.Props,
		})

		message := email.RenderMessage(ctx, c.TemplateName, c.From.Address, c.Props.Merge(to.Props))
		tags := []resend.Tag{
			{Name: "template", Value: c.TemplateName},
		}

		tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
		if ok && !env.IsSingleHostMode() {
			tags = append(tags, resend.Tag{Name: "tenant", Value: tenant.Subdomain})
		}

		resendInput := resend.SendEmailRequest{
			From:    c.From.String(),
			To:      []string{to.String()},
			Subject: message.Subject,
			Html:    message.Body,
			Tags:    tags,
		}

		result, err := resendClient.Emails.SendWithContext(ctx, &resendInput)
		if err != nil {
			log.Error(ctx, err)
			panic(errors.Wrap(err, "failed to send email with template %s", c.TemplateName))
		}

		log.Debugf(ctx, "Email sent with ID @{MessageId}.", dto.Props{
			"MessageId": result.Id,
		})
	}
}

func fetchRecentSupressions(ctx context.Context, q *query.FetchRecentSupressions) error {

	return nil
}
