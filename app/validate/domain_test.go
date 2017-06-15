package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/validate"
	. "github.com/onsi/gomega"
)

func TestInvalidSubdomains(t *testing.T) {
	RegisterTestingT(t)

	for _, subdomain := range []string{
		"me",
		"i",
		"signup",
		"",
		"1234567890123456789012345678901234567890ABC",
	} {
		ok, messages := validate.Subdomain(subdomain)
		Expect(ok).To(BeFalse())
		Expect(len(messages) > 0).To(BeTrue())
	}
}

func TestValidSubdomains(t *testing.T) {
	RegisterTestingT(t)

	for _, subdomain := range []string{
		"footbook",
		"yourcompany",
		"newyork",
	} {
		ok, messages := validate.Subdomain(subdomain)
		Expect(ok).To(BeTrue())
		Expect(len(messages)).To(Equal(0))
	}
}
