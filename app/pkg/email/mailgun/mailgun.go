package mailgun

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/log"
)

var baseURL = "https://api.mailgun.net/v3/%s/messages"

//Sender is used to send e-mails
type Sender struct {
	logger log.Logger
	domain string
	apiKey string
}

//NewSender creates a new mailgun e-mail sender
func NewSender(logger log.Logger, domain, apiKey string) *Sender {
	return &Sender{logger, domain, apiKey}
}

//Send an e-mail
func (s *Sender) Send(templateName string, params email.Params, from string, to email.Recipient) error {
	return s.BatchSend(templateName, params, from, []email.Recipient{to})
}

// BatchSend an e-mail to multiple recipients
func (s *Sender) BatchSend(templateName string, params email.Params, from string, to []email.Recipient) error {
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

	// Set Mailgun's var based on each recipient's variables
	recipientVariables := make(map[string]email.Params, 0)
	for _, r := range to {
		if email.CanSendTo(r.Address) {
			form.Add("to", fmt.Sprintf("%s <%s>", r.Name, r.Address))
			recipientVariables[r.Address] = r.Params
		} else {
			s.logger.Warnf("Skipping e-mail to '%s <%s>' due to whitelist.", r.Name, r.Address)
		}
	}

	// If we skipped all recipients due to whitelist, just return
	if len(recipientVariables) == 0 {
		return nil
	}

	if isBatch {
		json, err := json.Marshal(recipientVariables)
		if err != nil {
			return err
		}

		form.Add("recipient-variables", string(json))
	}

	if isBatch {
		s.logger.Debugf("Sending e-mail to %d recipients with template %s.", len(recipientVariables), templateName)
	} else {
		s.logger.Debugf("Sending e-mail to %s with template %s.", to[0].Address, templateName)
	}

	url := fmt.Sprintf(baseURL, s.domain)
	request, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	request.SetBasicAuth("api", s.apiKey)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		s.logger.Errorf("Failed to send e-mail")
		return err
	}
	s.logger.Debugf("E-mail sent with response code %d.", resp.StatusCode)
	return nil
}
