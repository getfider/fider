package validate_test

import (
	"testing"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/rand"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/storage/inmemory"
)

func TestInvalidEmail(t *testing.T) {
	RegisterT(t)

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
		rand.String(200) + "@gmail.com",
	} {
		result := validate.Email(email)
		Expect(result.Ok).IsFalse()
		Expect(len(result.Messages) > 0).IsTrue()
		Expect(result.Error).IsNil()
	}
}

func TestValidEmail(t *testing.T) {
	RegisterT(t)

	for _, email := range []string{
		"hello@company.com",
		"hello+alias@company.com",
		"abc@gmail.com",
	} {
		result := validate.Email(email)
		Expect(result.Ok).IsTrue()
		Expect(result.Messages).HasLen(0)
		Expect(result.Error).IsNil()
	}
}

func TestInvalidURL(t *testing.T) {
	RegisterT(t)

	for _, rawurl := range []string{
		"http//google.com",
		"google.com",
		"google",
		rand.String(301),
		"my@company",
	} {
		result := validate.URL(rawurl)
		Expect(result.Ok).IsFalse()
		Expect(len(result.Messages) > 0).IsTrue()
		Expect(result.Error).IsNil()
	}
}

func TestValidURL(t *testing.T) {
	RegisterT(t)

	for _, rawurl := range []string{
		"http://example.org",
		"https://example.org/oauth",
		"https://example.org/oauth?test=abc",
	} {
		result := validate.URL(rawurl)
		Expect(result.Ok).IsTrue()
		Expect(result.Messages).HasLen(0)
		Expect(result.Error).IsNil()
	}
}

func TestInvalidCNAME(t *testing.T) {
	RegisterT(t)

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
		result := validate.CNAME(inmemory.NewTenantStorage(), cname)
		Expect(result.Ok).IsFalse()
		Expect(len(result.Messages) > 0).IsTrue()
		Expect(result.Error).IsNil()
	}
}

func TestValidHostname(t *testing.T) {
	RegisterT(t)
	for _, cname := range []string{
		"google.com",
		"feedback.fider.io",
		"my.super.domain.com",
		"jon-snow.got.com",
		"got.com",
		"hi.m",
	} {
		result := validate.CNAME(inmemory.NewTenantStorage(), cname)
		Expect(result.Ok).IsTrue()
		Expect(result.Messages).HasLen(0)
		Expect(result.Error).IsNil()
	}
}

func TestValidCNAME_Availability(t *testing.T) {
	RegisterT(t)
	tenants := inmemory.NewTenantStorage()
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
		Expect(result.Ok).IsFalse()
		Expect(len(result.Messages) > 0).IsTrue()
		Expect(result.Error).IsNil()
	}
	for _, cname := range []string{
		"fider.footbook.com",
		"yourcompany.com",
		"anything.com",
	} {
		result := validate.CNAME(tenants, cname)
		Expect(result.Ok).IsTrue()
		Expect(result.Messages).HasLen(0)
		Expect(result.Error).IsNil()
	}
}
