package jobs

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"

	stripe "github.com/stripe/stripe-go/v83"
	stripesub "github.com/stripe/stripe-go/v83/subscription"
)

// DeleteScheduledTenantsJobHandler permanently deletes tenants whose grace period has
// elapsed. It processes at most one tenant per run so each deletion commits independently
// and a failure on one tenant cannot corrupt another (deletes also trigger non-transactional
// side effects — Stripe cancellation and blob storage wipes).
type DeleteScheduledTenantsJobHandler struct {
}

func (j DeleteScheduledTenantsJobHandler) Schedule() string {
	return "0 */5 * * * *" // every 5 minutes
}

func (j DeleteScheduledTenantsJobHandler) Run(ctx Context) error {
	pending := &query.GetTenantsPendingDeletion{}
	if err := bus.Dispatch(ctx, pending); err != nil {
		return errors.Wrap(err, "failed to fetch tenants pending deletion")
	}

	if len(pending.Result) == 0 {
		return nil
	}

	tenant := pending.Result[0]
	tenantCtx := context.WithValue(ctx.Context, app.TenantCtxKey, tenant)

	log.Warnf(ctx, "Deleting tenant @{TenantID} (@{Subdomain}) — grace period elapsed", dto.Props{
		"TenantID":  tenant.ID,
		"Subdomain": tenant.Subdomain,
	})

	// Capture the owner before deletion so we can send the completion email afterwards.
	owner := &query.GetTenantOwner{TenantID: tenant.ID}
	if err := bus.Dispatch(ctx, owner); err != nil {
		log.Warnf(ctx, "Could not resolve owner for tenant @{TenantID}; no completion email will be sent", dto.Props{
			"TenantID": tenant.ID,
		})
		owner.Result = nil
	}

	// Stripe cancellation is a non-transactional side effect: do it first and bail (retry next
	// run) on a real error, but treat an already-missing subscription as success.
	if err := cancelStripeSubscription(tenantCtx); err != nil {
		return errors.Wrap(err, "failed to cancel stripe subscription for tenant '%d'", tenant.ID)
	}

	// Wipe blob storage (S3/fs files live under a tenant prefix; sql-backed blobs also fall to
	// the DB delete below, but listing+deleting here is harmless and correct for all backends).
	if err := wipeTenantBlobs(tenantCtx); err != nil {
		return errors.Wrap(err, "failed to wipe blobs for tenant '%d'", tenant.ID)
	}

	if err := bus.Dispatch(ctx, &cmd.DeleteTenant{TenantID: tenant.ID}); err != nil {
		return errors.Wrap(err, "failed to delete tenant '%d'", tenant.ID)
	}

	sendCompletionEmail(ctx, tenant, owner.Result)
	return nil
}

func cancelStripeSubscription(tenantCtx context.Context) error {
	billing := &query.GetStripeBillingState{}
	if err := bus.Dispatch(tenantCtx, billing); err != nil {
		return errors.Wrap(err, "failed to read stripe billing state")
	}

	subID := billing.Result.SubscriptionID
	if subID == "" {
		return nil
	}

	stripe.Key = env.Config.Stripe.SecretKey
	if _, err := stripesub.Cancel(subID, nil); err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok && stripeErr.Code == stripe.ErrorCodeResourceMissing {
			// Subscription already cancelled or never existed — nothing to do.
			return nil
		}
		return err
	}
	return nil
}

func wipeTenantBlobs(tenantCtx context.Context) error {
	list := &query.ListBlobs{}
	if err := bus.Dispatch(tenantCtx, list); err != nil {
		return errors.Wrap(err, "failed to list blobs")
	}

	for _, key := range list.Result {
		if err := bus.Dispatch(tenantCtx, &cmd.DeleteBlob{Key: key}); err != nil {
			return errors.Wrap(err, "failed to delete blob '%s'", key)
		}
	}
	return nil
}

func sendCompletionEmail(ctx Context, tenant *entity.Tenant, owner *entity.User) {
	if owner == nil || owner.Email == "" {
		return
	}

	to := dto.NewRecipient(owner.Name, owner.Email, dto.Props{})
	bus.Publish(ctx, &cmd.SendMail{
		From:         dto.Recipient{Name: "Fider"},
		To:           []dto.Recipient{to},
		TemplateName: "delete_account_completed",
		Props: dto.Props{
			"tenantName": tenant.Name,
			"subdomain":  tenant.Subdomain,
			"logo":       web.LogoURL(ctx),
		},
	})
}
