package email_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/web"

	. "github.com/onsi/gomega"
)

func TestRenderMessage(t *testing.T) {
	RegisterTestingT(t)

	message := email.RenderMessage("echo_test", web.Map{
		"name": "Fider",
	})
	Expect(message.Subject).To(Equal("Message to: Fider"))
	Expect(message.Body).To(Equal("Hello World Fider!"))
}
