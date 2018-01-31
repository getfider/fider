package mailgun

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/pkg/email"
)

var baseURL = "https://api.mailgun.net/v3/%s/messages"

//Sender is used to send e-mails
type Sender struct {
	domain string
	apiKey string
}

//NewSender creates a new mailgun e-mail sender
func NewSender(domain, apiKey string) *Sender {
	return &Sender{domain, apiKey}
}

//Send an e-mail
func (s *Sender) Send(from, to, templateName string, params map[string]interface{}) error {

	message := email.RenderMessage(templateName, params)

	form := url.Values{}
	form.Add("from", fmt.Sprintf("%s <%s>", from, email.NoReply))
	form.Add("to", to)
	form.Add("subject", message.Subject)
	form.Add("html", message.Body)

	url := fmt.Sprintf(baseURL, s.domain)
	request, err := http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	request.SetBasicAuth("api", s.apiKey)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = http.DefaultClient.Do(request)
	return err
}
