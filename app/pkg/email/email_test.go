package email_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/email"

	. "github.com/onsi/gomega"
)

func TestRenderMessage(t *testing.T) {
	RegisterTestingT(t)

	message := email.RenderMessage("echo_test", email.Params{
		"name": "Fider",
	})
	Expect(message.Subject).To(Equal("Message to: Fider"))
	Expect(message.Body).To(Equal("Hello World Fider!"))
}

func TestEmailWhitelist_Valid(t *testing.T) {
	RegisterTestingT(t)

	email.SetWhitelist("(^.+@fider.io$)|(^darthvader\\.fider(\\+.*)?@gmail\\.com$)")

	for _, address := range []string{
		"me@fider.io",
		"john@fider.io",
		"me+123@fider.io",
		"darthvader.fider@gmail.com",
		"darthvader.fider+1234@gmail.com",
		"darthvader.fider+434@gmail.com",
	} {
		Expect(email.CanSendTo(address)).To(BeTrue())
	}

	for _, address := range []string{
		"",
		"me@fidero.io",
		"@fider.io",
		"me+123@fider.iod",
		"darthvader@gmail.com",
		"darthvader.fider@gmaila.com",
		"darthvader.fidera@gmail.com",
	} {
		Expect(email.CanSendTo(address)).To(BeFalse())
	}
}
