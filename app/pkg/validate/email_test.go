package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/validate"
	. "github.com/onsi/gomega"
)

func TestInvalidEmail(t *testing.T) {
	RegisterTestingT(t)

	for _, email := range []string{
		"hello",
		"",
		"my@company",
		"my @company.com",
		"my@.company.com",
		"my+company.com",
		".my@company.com",
		"my@company@other.com",
		"@gmail.com",
		"abc12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890@gmail.com",
	} {
		result := validate.Email(email)
		Expect(result.Ok).To(BeFalse())
		Expect(len(result.Messages) > 0).To(BeTrue())
		Expect(result.Error).To(BeNil())
	}
}

func TestValidEmail(t *testing.T) {
	RegisterTestingT(t)

	for _, email := range []string{
		"hello@company.com",
		"hello+alias@company.com",
		"abc@gmail.com",
	} {
		result := validate.Email(email)
		Expect(result.Ok).To(BeTrue())
		Expect(len(result.Messages)).To(Equal(0))
		Expect(result.Error).To(BeNil())
	}
}
