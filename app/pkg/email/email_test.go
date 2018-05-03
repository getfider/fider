package email_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/email"

	. "github.com/getfider/fider/app/pkg/assert"
)

func TestRenderMessage(t *testing.T) {
	RegisterT(t)

	message := email.RenderMessage("echo_test", email.Params{
		"name": "Fider",
	})
	Expect(message.Subject).Equals("Message to: Fider")
	Expect(message.Body).Equals("Hello World Fider!")
}

func TestEmailWhitelist_Valid(t *testing.T) {
	RegisterT(t)

	email.SetWhitelist("(^.+@fider.io$)|(^darthvader\\.fider(\\+.*)?@gmail\\.com$)")

	for _, address := range []string{
		"me@fider.io",
		"john@fider.io",
		"me+123@fider.io",
		"darthvader.fider@gmail.com",
		"darthvader.fider+1234@gmail.com",
		"darthvader.fider+434@gmail.com",
	} {
		Expect(email.CanSendTo(address)).IsTrue()
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
		Expect(email.CanSendTo(address)).IsFalse()
	}
}

func TestParamsMerge(t *testing.T) {
	RegisterT(t)

	p1 := email.Params{
		"name": "Jon",
		"age":  26,
	}
	p2 := p1.Merge(email.Params{
		"age":   30,
		"email": "john.snow@got.com",
	})
	Expect(p2).Equals(email.Params{
		"name":  "Jon",
		"age":   30,
		"email": "john.snow@got.com",
	})
}
