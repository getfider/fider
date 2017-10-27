package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/storage/inmemory"
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

func TestInvalidCNAME(t *testing.T) {
	RegisterTestingT(t)

	for _, cname := range []string{
		"hello",
		"hellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohellohello.com",
		"",
		"my",
		"name.com/abc",
		"feedback.test.fider.io",
		"test.fider.io",
		"@google.com",
	} {
		result := validate.CNAME(&inmemory.TenantStorage{}, cname)
		Expect(result.Ok).To(BeFalse())
		Expect(len(result.Messages) > 0).To(BeTrue())
		Expect(result.Error).To(BeNil())
	}
}

func TestValidHostname(t *testing.T) {
	RegisterTestingT(t)
	for _, cname := range []string{
		"google.com",
		"feedback.fider.io",
		"my.super.domain.com",
		"jon-snow.got.com",
		"got.com",
		"hi.m",
	} {
		result := validate.CNAME(&inmemory.TenantStorage{}, cname)
		Expect(result.Ok).To(BeTrue())
		Expect(len(result.Messages)).To(Equal(0))
		Expect(result.Error).To(BeNil())
	}
}

func TestValidCNAME_Availability(t *testing.T) {
	RegisterTestingT(t)
	tenants := &inmemory.TenantStorage{}
	tenant, _ := tenants.Add("Footbook", "footbook", models.TenantActive)
	tenant.CNAME = "footbook.com"
	tenant, _ = tenants.Add("Your Company", "yourcompany", models.TenantActive)
	tenant.CNAME = "fider.yourcompany.com"
	tenant, _ = tenants.Add("New York", "newyork", models.TenantActive)
	tenant.CNAME = "feedback.newyork.com"
	for _, cname := range []string{
		"footbook.com",
		"fider.yourcompany.com",
		"feedback.newyork.com",
	} {
		result := validate.CNAME(tenants, cname)
		Expect(result.Ok).To(BeFalse())
		Expect(len(result.Messages) > 0).To(BeTrue())
		Expect(result.Error).To(BeNil())
	}
	for _, cname := range []string{
		"fider.footbook.com",
		"yourcompany.com",
		"anything.com",
	} {
		result := validate.CNAME(tenants, cname)
		Expect(result.Ok).To(BeTrue())
		Expect(len(result.Messages)).To(Equal(0))
		Expect(result.Error).To(BeNil())
	}
}
