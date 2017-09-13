package email

import "github.com/getfider/fider/app/pkg/env"

// NoReply is the default 'from' address
var NoReply = env.MustGet("NOREPLY_EMAIL")

// Sender is used to send e-mails
type Sender interface {
	Send(from, to, subject, message string) error
}
