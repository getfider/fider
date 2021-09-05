package mailgun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/getfider/fider/app"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/services/email"
)

func sendMail(ctx context.Context, c *cmd.SendMail) {
	if len(c.To) == 0 {
		return
	}

	if c.Props == nil {
		c.Props = dto.Props{}
	}

	isBatch := len(c.To) > 1

	var message *email.Message
	if isBatch {
		// Replace recipient specific Go templates variables with Mailgun template variables
		if c.To[0].Props != nil {
			for k := range c.To[0].Props {
				c.Props[k] = fmt.Sprintf("%%recipient.%s%%", k)
			}
		}
		message = email.RenderMessage(ctx, c.TemplateName, c.Props)
	} else {
		message = email.RenderMessage(ctx, c.TemplateName, c.Props.Merge(c.To[0].Props))
	}

	form := url.Values{}
	form.Add("from", dto.NewRecipient(c.From, email.NoReply, dto.Props{}).String())
	form.Add("h:Reply-To", email.NoReply)
	form.Add("subject", message.Subject)
	form.Add("html", message.Body)
	form.Add("o:tag", fmt.Sprintf("template:%s", c.TemplateName))

	tenant, ok := ctx.Value(app.TenantCtxKey).(*entity.Tenant)
	if ok && !env.IsSingleHostMode() {
		form.Add("o:tag", fmt.Sprintf("tenant:%s", tenant.Subdomain))
	}

	// Set Mailgun's var based on each recipient's variables
	recipientVariables := make(map[string]dto.Props)
	for _, r := range c.To {
		if r.Address != "" {
			if email.CanSendTo(r.Address) {
				form.Add("to", r.String())
				recipientVariables[r.Address] = r.Props
			} else {
				log.Warnf(ctx, "Skipping email to '@{Name} <@{Address}>'.", dto.Props{
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
		log.Debugf(ctx, "Sending email to @{CountRecipients} recipients with template @{TemplateName}.", dto.Props{
			"CountRecipients": len(recipientVariables),
			"TemplateName":    c.TemplateName,
		})
	} else {
		log.Debugf(ctx, "Sending email to @{Address} with template @{TemplateName}.", dto.Props{
			"Address":      c.To[0].Address,
			"TemplateName": c.TemplateName,
		})
	}

	req := &cmd.HTTPRequest{
		Method: "POST",
		URL:    getEndpoint(ctx, env.Config.Email.Mailgun.Domain, "/messages"),
		Body:   strings.NewReader(form.Encode()),
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		BasicAuth: &dto.BasicAuth{
			User:     "api",
			Password: env.Config.Email.Mailgun.APIKey,
		},
	}
	err := bus.Dispatch(ctx, req)
	if err != nil {
		panic(errors.Wrap(err, "failed to send email with template %s", c.TemplateName))
	}
	log.Debugf(ctx, "Email sent with response code @{StatusCode}.", dto.Props{
		"StatusCode": req.ResponseStatusCode,
	})
}
