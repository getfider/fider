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
		result := validate.Subdomain(nil, subdomain)
		Expect(result.Ok).IsFalse()
		Expect(len(result.Messages) > 0).IsTrue()
		Expect(result.Err).IsNil()
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
		result := validate.Subdomain(tenants, subdomain)
		Expect(result.Ok).IsFalse()
		Expect(len(result.Messages) > 0).IsTrue()
		Expect(result.Err).IsNil()
	}

	for _, subdomain := range []string{
		"my-company",
		"123-company",
	} {
		result := validate.Subdomain(tenants, subdomain)
		Expect(result.Ok).IsTrue()
		Expect(result.Messages).HasLen(0)
		Expect(result.Err).IsNil()
	}
}
