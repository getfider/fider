package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestInvalidSubdomains(t *testing.T) {
	RegisterTestingT(t)

	for _, subdomain := range []string{
		"me",
		"i",
		"signup",
		"",
		"my company",
		"my@company",
		"my.company",
		"my+company",
		"1234567890123456789012345678901234567890ABC",
	} {
		result := validate.Subdomain(nil, subdomain)
		Expect(result.Ok).To(BeFalse())
		Expect(len(result.Messages) > 0).To(BeTrue())
		Expect(result.Error).To(BeNil())
	}
}

func TestValidSubdomains_Availability(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	tenants.Add("Footbook", "footbook")
	tenants.Add("Your Company", "yourcompany")
	tenants.Add("New York", "newyork")

	for _, subdomain := range []string{
		"footbook",
		"yourcompany",
		"newyork",
		"NewYork",
	} {
		result := validate.Subdomain(tenants, subdomain)
		Expect(result.Ok).To(BeFalse())
		Expect(len(result.Messages) > 0).To(BeTrue())
		Expect(result.Error).To(BeNil())
	}

	for _, subdomain := range []string{
		"my-company",
		"123-company",
	} {
		result := validate.Subdomain(tenants, subdomain)
		Expect(result.Ok).To(BeTrue())
		Expect(len(result.Messages)).To(Equal(0))
		Expect(result.Error).To(BeNil())
	}
}
