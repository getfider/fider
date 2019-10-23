package validate_test

import (
	"context"
	"testing"

	"github.com/getfider/fider/app/models/query"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/validate"
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
		messages, err := validate.Subdomain(context.Background(), subdomain)
		Expect(len(messages) > 0).IsTrue()
		Expect(err).IsNil()
	}
}

func TestValidSubdomains_Availability(t *testing.T) {
	RegisterT(t)

	bus.AddHandler(func(ctx context.Context, q *query.IsSubdomainAvailable) error {
		q.Result = q.Subdomain != "footbook" && q.Subdomain != "yourcompany" && q.Subdomain != "newyork"
		return nil
	})

	for _, subdomain := range []string{
		"footbook",
		"yourcompany",
		"newyork",
		"NewYork",
	} {
		messages, err := validate.Subdomain(context.Background(), subdomain)
		Expect(len(messages) > 0).IsTrue()
		Expect(err).IsNil()
	}

	for _, subdomain := range []string{
		"my-company",
		"123-company",
	} {
		messages, err := validate.Subdomain(context.Background(), subdomain)
		Expect(messages).HasLen(0)
		Expect(err).IsNil()
	}
}
