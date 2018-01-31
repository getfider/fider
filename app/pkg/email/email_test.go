package email_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/web"

	. "github.com/onsi/gomega"
)

func TestRenderMessage(t *testing.T) {
	RegisterTestingT(t)

	message := email.RenderMessage("echo", web.Map{
		"name": "Fider",
	})
	Expect(message).To(Equal("Hello World Fider!"))

	message = email.RenderMessage("echo", web.Map{
		"name": "Golang",
	})
	Expect(message).To(Equal("Hello World Golang!"))
}
