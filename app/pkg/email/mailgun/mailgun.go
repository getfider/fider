package mailgun

import (
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
func (s *Sender) Send(from, to, templateName string, params map[string]interface{}) (*email.Message, error) {
	s.logger.Debugf("Sending e-mail to %s with template %s and params %s.", to, templateName, params)

	message := email.RenderMessage(templateName, params)
	message.From = fmt.Sprintf("%s <%s>", from, email.NoReply)
	message.To = to

	form := url.Values{}
	form.Add("from", message.From)
	form.Add("to", message.To)
	form.Add("subject", message.Subject)
	form.Add("html", message.Body)

	url := fmt.Sprintf(baseURL, s.domain)
	request, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth("api", s.apiKey)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(request)
	if err == nil {
		s.logger.Debugf("E-mail sent with response code %d.", resp.StatusCode)
		return message, nil
	}
	s.logger.Errorf("Failed to send e-mail")
	return nil, err
}
