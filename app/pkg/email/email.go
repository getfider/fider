package email

//Sender is used to send e-mails
type Sender interface {
	Send(from, to, subject, message string) error
}
