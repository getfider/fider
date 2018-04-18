package mailgun

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
)

var baseURL = "https://api.mailgun.net/v3/%s/messages"

//Sender is used to send emails
type Sender struct {
	logger log.Logger
	domain string
	apiKey string
}

//NewSender creates a new mailgun email sender
func NewSender(logger log.Logger, domain, apiKey string) *Sender {
	return &Sender{logger, domain, apiKey}
}

//Send an email
func (s *Sender) Send(tenant *models.Tenant, templateName string, params email.Params, from string, to email.Recipient) error {
	return s.BatchSend(tenant, templateName, params, from, []email.Recipient{to})
}

// BatchSend an email to multiple recipients
func (s *Sender) BatchSend(tenant *models.Tenant, templateName string, params email.Params, from string, to []email.Recipient) error {
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
		message = email.RenderMessage(templateName, params)
	} else {
		message = email.RenderMessage(templateName, params.Merge(to[0].Params))
	}

	form := url.Values{}
	form.Add("from", fmt.Sprintf("%s <%s>", from, email.NoReply))
	form.Add("subject", message.Subject)
	form.Add("html", message.Body)
	form.Add("o:tag", fmt.Sprintf("template:%s", templateName))
	if tenant != nil && !env.IsSingleHostMode() {
		form.Add("o:tag", fmt.Sprintf("tenant:%s", tenant.Subdomain))
	}

	// Set Mailgun's var based on each recipient's variables
	recipientVariables := make(map[string]email.Params, 0)
	for _, r := range to {
		if r.Address != "" {
			if email.CanSendTo(r.Address) {
				form.Add("to", fmt.Sprintf("%s <%s>", r.Name, r.Address))
				recipientVariables[r.Address] = r.Params
			} else {
				s.logger.Warnf("Skipping email to '%s <%s>' due to whitelist.", r.Name, r.Address)
			}
		}
	}

	// If we skipped all recipients due to whitelist, just return
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
		s.logger.Debugf("Sending email to %d recipients with template %s.", len(recipientVariables), templateName)
	} else {
		s.logger.Debugf("Sending email to %s with template %s.", to[0].Address, templateName)
	}

	url := fmt.Sprintf(baseURL, s.domain)
	request, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return errors.Wrap(err, "failed to create POST request")
	}

	request.SetBasicAuth("api", s.apiKey)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.Wrap(err, "failed to send email with template %s", templateName)
	}

	defer resp.Body.Close()
	s.logger.Debugf("Email sent with response code %d.", resp.StatusCode)
	return nil
}
