package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/storage/inmemory"
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
		"my company",
		"my@company",
		"my.company",
		"my+company",
		"1234567890123456789012345678901234567890ABC",
	} {
		ok, messages, err := validate.Subdomain(nil, subdomain)
		Expect(ok).To(BeFalse())
		Expect(len(messages) > 0).To(BeTrue())
		Expect(err).To(BeNil())
	}
}

func TestValidSubdomains_Availability(t *testing.T) {
	RegisterTestingT(t)

	tenants := &inmemory.TenantStorage{}
	tenants.Add(&models.Tenant{Subdomain: "footbook"})
	tenants.Add(&models.Tenant{Subdomain: "yourcompany"})
	tenants.Add(&models.Tenant{Subdomain: "newyork"})

	for _, subdomain := range []string{
		"footbook",
		"yourcompany",
		"newyork",
		"NewYork",
	} {
		ok, messages, err := validate.Subdomain(tenants, subdomain)
		Expect(ok).To(BeFalse())
		Expect(len(messages)).To(Equal(1))
		Expect(err).To(BeNil())
	}

	for _, subdomain := range []string{
		"my-company",
		"123-company",
	} {
		ok, messages, err := validate.Subdomain(tenants, subdomain)
		Expect(ok).To(BeTrue())
		Expect(len(messages)).To(Equal(0))
		Expect(err).To(BeNil())
	}
}
