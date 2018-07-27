package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/storage/inmemory"
)

func TestInvalidSubdomains(t *testing.T) {
	RegisterT(t)

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
		messages, err := validate.Subdomain(nil, subdomain)
		Expect(len(messages) > 0).IsTrue()
		Expect(err).IsNil()
	}
}

func TestValidSubdomains_Availability(t *testing.T) {
	RegisterT(t)

	tenants := inmemory.NewTenantStorage()
	tenants.Add("Footbook", "footbook", models.TenantActive)
	tenants.Add("Your Company", "yourcompany", models.TenantActive)
	tenants.Add("New York", "newyork", models.TenantActive)

	for _, subdomain := range []string{
		"footbook",
		"yourcompany",
		"newyork",
		"NewYork",
	} {
		messages, err := validate.Subdomain(tenants, subdomain)
		Expect(len(messages) > 0).IsTrue()
		Expect(err).IsNil()
	}

	for _, subdomain := range []string{
		"my-company",
		"123-company",
	} {
		messages, err := validate.Subdomain(tenants, subdomain)
		Expect(messages).HasLen(0)
		Expect(err).IsNil()
	}
}
