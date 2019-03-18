package mailgun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/services/email"
	"github.com/getfider/fider/app/services/httpclient"
)

var baseURL = "https://api.mailgun.net/v3/%s/messages"

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Category() string {
	return "email"
}

func (s Service) Enabled() bool {
	return env.Config.Email.Mailgun.APIKey != ""
}

func (s Service) Init() {
	bus.AddEventListener(sendEmail)
}

func sendEmail(ctx context.Context, cmd *email.SendMessageCommand) {
	if len(cmd.To) == 0 {
		return
	}

	if cmd.Params == nil {
		cmd.Params = email.Params{}
	}

	isBatch := len(cmd.To) > 1

	var message *email.Message
	if isBatch {
		// Replace recipient specific Go templates variables with Mailgun template variables
		if cmd.To[0].Params != nil {
			for k := range cmd.To[0].Params {
				cmd.Params[k] = fmt.Sprintf("%%recipient.%s%%", k)
			}
		}
		message = email.RenderMessage(cmd.TemplateName, cmd.Params)
	} else {
		message = email.RenderMessage(cmd.TemplateName, cmd.Params.Merge(cmd.To[0].Params))
	}

	form := url.Values{}
	form.Add("from", email.NewRecipient(cmd.From, email.NoReply, email.Params{}).String())
	form.Add("h:Reply-To", email.NoReply)
	form.Add("subject", message.Subject)
	form.Add("html", message.Body)
	form.Add("o:tag", fmt.Sprintf("template:%s", cmd.TemplateName))

	tenant, ok := ctx.Value(app.TenantCtxKey).(*models.Tenant)
	if ok && !env.IsSingleHostMode() {
		form.Add("o:tag", fmt.Sprintf("tenant:%s", tenant.Subdomain))
	}

	// Set Mailgun's var based on each recipient's variables
	recipientVariables := make(map[string]email.Params)
	for _, r := range cmd.To {
		if r.Address != "" {
			if email.CanSendTo(r.Address) {
				form.Add("to", r.String())
				recipientVariables[r.Address] = r.Params
			} else {
				log.Warnf(ctx, "Skipping email to '@{Name} <@{Address}>'.", log.Props{
					"Name":    r.Name,
					"Address": r.Address,
				})
			}
		}
	}

	// If we skipped all recipients, just return
	if len(recipientVariables) == 0 {
		return
	}

	if isBatch {
		json, err := json.Marshal(recipientVariables)
		if err != nil {
			panic(errors.Wrap(err, "failed to marshal recipient variables"))
		}

		form.Add("recipient-variables", string(json))
	}

	if isBatch {
		log.Debugf(ctx, "Sending email to @{CountRecipients} recipients with template @{TemplateName}.", log.Props{
			"CountRecipients": len(recipientVariables),
			"TemplateName":    cmd.TemplateName,
		})
	} else {
		log.Debugf(ctx, "Sending email to @{Address} with template @{TemplateName}.", log.Props{
			"Address":      cmd.To[0].Address,
			"TemplateName": cmd.TemplateName,
		})
	}

	url := fmt.Sprintf(baseURL, env.Config.Email.Mailgun.Domain)

	req := &httpclient.Request{
		Method: "POST",
		URL:    url,
		Body:   strings.NewReader(form.Encode()),
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		BasicAuth: &httpclient.BasicAuth{
			User:     "api",
			Password: env.Config.Email.Mailgun.APIKey,
		},
	}
	err := bus.Dispatch(ctx, req)
	if err != nil {
		panic(errors.Wrap(err, "failed to send email with template %s", cmd.TemplateName))
	}
	log.Debugf(ctx, "Email sent with response code @{StatusCode}.", log.Props{
		"StatusCode": req.ResponseStatusCode,
	})
}
