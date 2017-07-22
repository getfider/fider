package im_test

import (
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/im"
	"github.com/getfider/fider/app/storage/inmemory"
	. "github.com/onsi/gomega"
)

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: ""}
	result := input.Validate(nil)
	ExpectFailed(result)
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: jonSnowToken, Name: ""}
	result := input.Validate(nil)
	ExpectFailed(result)
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: jonSnowToken, Name: "My Company"}
	result := input.Validate(&app.Services{
		Tenants: &inmemory.TenantStorage{},
	})
	ExpectFailed(result)
}
