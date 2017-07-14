package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/validate"
	. "github.com/onsi/gomega"
)

func TestInvalidIdeaTitles(t *testing.T) {
	RegisterTestingT(t)

	for _, title := range []string{
		"me",
		"",
		"  ",
		"signup",
		"my company",
		"my@company",
		"my.company",
		"my+company",
		"1234567890123456789012345678901234567890ABC",
	} {
		ok, messages, err := validate.Idea(title, "")
		Expect(ok).To(BeFalse())
		Expect(len(messages) > 0).To(BeTrue())
		Expect(err).To(BeNil())
	}
}

func TestValidIdeaTitles(t *testing.T) {
	RegisterTestingT(t)

	for _, title := range []string{
		"this is my new idea",
		"this idea is very descriptive",
	} {
		ok, messages, err := validate.Idea(title, "")
		Expect(ok).To(BeTrue())
		Expect(len(messages)).To(Equal(0))
		Expect(err).To(BeNil())
	}
}
