package im_test

import (
	"testing"

	"github.com/getfider/fider/app/models/im"
	. "github.com/onsi/gomega"
)

func TestCreateTenant_EmptyToken(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: ""}
	result := input.Validate(services)
	ExpectFailed(result, "token")
}

func TestCreateTenant_EmptyName(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: jonSnowToken, Name: ""}
	result := input.Validate(services)
	ExpectFailed(result, "name")
}

func TestCreateTenant_EmptySubdomain(t *testing.T) {
	RegisterTestingT(t)

	input := im.CreateTenant{Token: jonSnowToken, Name: "My Company"}
	result := input.Validate(services)
	ExpectFailed(result, "subdomain")
}
