package dbEntities

import (
	"github.com/getfider/fider/app/pkg/dbx"
)

type StripeBillingState struct {
	StripeCustomerID     dbx.NullString `db:"stripe_customer_id"`
	StripeSubscriptionID dbx.NullString `db:"stripe_subscription_id"`
	LicenseKey           dbx.NullString `db:"license_key"`
}
