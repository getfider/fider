package mailgun

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web/http"
)

// Known base URLs
// Should Mailgun add other regions we'll just need to add their URLs here
// Use upper case keys - incoming env var values are normalized before being used
var baseURLs = map[string]string{
	"US" : "https://api.mailgun.net/v3/%s/messages",
	"EU" : "https://api.eu.mailgun.net/v3/%s/messages",
}

//Sender is used to send emails
type Sender struct {
	logger log.Logger
	client http.Client
	domain string
	apiKey string
}

//NewSender creates a new mailgun email sender
func NewSender(logger log.Logger, client http.Client, domain, apiKey string) *Sender {
	return &Sender{logger, client, domain, apiKey}
}

//Send an email
func (s *Sender) Send(ctx email.Context, templateName string, params email.Params, from string, to email.Recipient) error {
	return s.BatchSend(ctx, templateName, params, from, []email.Recipient{to})
}

// Try getting the base URL of the Mailgun API using Environment vars and the Sender's domain
// Fall back to the URL for the US region if that fails to maintain compatibility with older installs
func (s *Sender) GetBaseURL() string {
	var regionCode = env.Config.Email.Mailgun.Region
	regionCode = strings.ToUpper(regionCode)

	// Default to the US domain if no region code was provided (ENV not set)
	// or if the provided code isn't valid
	if len(regionCode) < 1 {
		regionCode = "US"
	} else if len(baseURLs[regionCode]) < 1 {
		s.logger.Warnf(
			"Unknown Mailgun region code '@{Code}' configured - falling back to 'US'", 
			log.Props{
				"Code": env.Config.Email.Mailgun.Region,
			},
		)

		regionCode = "US"
	}

	return fmt.Sprintf(baseURLs[regionCode], s.domain)
}

// BatchSend an email to multiple recipients
func (s *Sender) BatchSend(ctx email.Context, templateName string, params email.Params, from string, to []email.Recipient) error {
	if len(to) == 0 {
		return nil
	}

	isBatch := len(to) > 1

	var message *email.Message
	if isBatch {
		// Replace recipient specific Go templates variables with Mailgun template variables
		for k := range to[0].Params {
			params[k] = fmt.Sprintf("%%recipient.%s%%", k)
		}
		message = email.RenderMessage(ctx, templateName, params)
	} else {
		message = email.RenderMessage(ctx, templateName, params.Merge(to[0].Params))
	}

	form := url.Values{}
	form.Add("from", email.NewRecipient(from, email.NoReply, email.Params{}).String())
	form.Add("h:Reply-To", email.NoReply)
	form.Add("subject", message.Subject)
	form.Add("html", message.Body)
	form.Add("o:tag", fmt.Sprintf("template:%s", templateName))
	if ctx.Tenant() != nil && !env.IsSingleHostMode() {
		form.Add("o:tag", fmt.Sprintf("tenant:%s", ctx.Tenant().Subdomain))
	}

	// Set Mailgun's var based on each recipient's variables
	recipientVariables := make(map[string]email.Params, 0)
	for _, r := range to {
		if r.Address != "" {
			if email.CanSendTo(r.Address) {
				form.Add("to", r.String())
				recipientVariables[r.Address] = r.Params
			} else {
				s.logger.Warnf("Skipping email to '@{Name} <@{Address}>'.", log.Props{
					"Name":    r.Name,
					"Address": r.Address,
				})
			}
		}
	}

	// If we skipped all recipients, just return
	if len(recipientVariables) == 0 {
		return nil
	}

	if isBatch {
		json, err := json.Marshal(recipientVariables)
		if err != nil {
			return errors.Wrap(err, "failed to marshal recipient variables")
		}

		form.Add("recipient-variables", string(json))
	}

	if isBatch {
		s.logger.Debugf("Sending email to @{CountRecipients} recipients with template @{TemplateName}.", log.Props{
			"CountRecipients": len(recipientVariables),
			"TemplateName":    templateName,
		})
	} else {
		s.logger.Debugf("Sending email to @{Address} with template @{TemplateName}.", log.Props{
			"Address":      to[0].Address,
			"TemplateName": templateName,
		})
	}

	request, err := http.NewRequest("POST", s.GetBaseURL(), strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "failed to create POST request")
	}

	request.SetBasicAuth("api", s.apiKey)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to send email with template %s", templateName)
	}

	defer resp.Body.Close()
	s.logger.Debugf("Email sent with response code @{StatusCode}.", log.Props{
		"StatusCode": resp.StatusCode,
	})
	return nil
}
