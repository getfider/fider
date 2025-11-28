package postgres

import (
	"context"
	"time"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

func isCNAMEAvailable(ctx context.Context, q *query.IsCNAMEAvailable) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		tenantID := 0
		if tenant != nil {
			tenantID = tenant.ID
		}

		exists, err := trx.Exists("SELECT id FROM tenants WHERE cname = $1 AND id <> $2", q.CNAME, tenantID)
		if err != nil {
			q.Result = false
			return errors.Wrap(err, "failed to check if tenant exists with CNAME '%s'", q.CNAME)
		}
		q.Result = !exists
		return nil
	})
}

func isSubdomainAvailable(ctx context.Context, q *query.IsSubdomainAvailable) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		exists, err := trx.Exists("SELECT id FROM tenants WHERE subdomain = $1", q.Subdomain)
		if err != nil {
			q.Result = false
			return errors.Wrap(err, "failed to check if tenant exists with subdomain '%s'", q.Subdomain)
		}
		q.Result = !exists
		return nil
	})
}

func updateTenantPrivacySettings(ctx context.Context, c *cmd.UpdateTenantPrivacySettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute("UPDATE tenants SET is_private = $1 WHERE id = $2", c.IsPrivate, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant privacy setting")
		}
		_, err = trx.Execute("UPDATE tenants SET is_feed_enabled = $1 WHERE id = $2", c.IsFeedEnabled, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant feed setting")
		}
		_, err = trx.Execute("UPDATE tenants SET is_moderation_enabled = $1 WHERE id = $2", c.IsModerationEnabled, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant moderation setting")
		}
		return nil
	})
}

func updateTenantEmailAuthAllowedSettings(ctx context.Context, c *cmd.UpdateTenantEmailAuthAllowedSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		_, err := trx.Execute("UPDATE tenants SET is_email_auth_allowed = $1 WHERE id = $2", c.IsEmailAuthAllowed, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant allowing email auth settings")
		}
		return nil
	})
}

func updateTenantSettings(ctx context.Context, c *cmd.UpdateTenantSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if c.Logo.Remove {
			c.Logo.BlobKey = ""
		}

		query := "UPDATE tenants SET name = $1, invitation = $2, welcome_message = $3, cname = $4, logo_bkey = $5, locale = $6 WHERE id = $7"
		_, err := trx.Execute(query, c.Title, c.Invitation, c.WelcomeMessage, c.CNAME, c.Logo.BlobKey, c.Locale, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant settings")
		}

		tenant.Name = c.Title
		tenant.Invitation = c.Invitation
		tenant.CNAME = c.CNAME
		tenant.WelcomeMessage = c.WelcomeMessage

		return nil
	})
}

func updateTenantAdvancedSettings(ctx context.Context, c *cmd.UpdateTenantAdvancedSettings) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		AllowedSchemes := c.AllowedSchemes
		if !env.Config.AllowAllowedSchemes {
			AllowedSchemes = ""
		}

		query := "UPDATE tenants SET custom_css = $1, allowed_schemes = $2 WHERE id = $3"
		_, err := trx.Execute(query, c.CustomCSS, AllowedSchemes, tenant.ID)
		if err != nil {
			return errors.Wrap(err, "failed update tenant advanced settings")
		}

		tenant.CustomCSS = c.CustomCSS
		tenant.AllowedSchemes = AllowedSchemes
		return nil
	})
}

func activateTenant(ctx context.Context, c *cmd.ActivateTenant) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		query := "UPDATE tenants SET status = $1 WHERE id = $2"
		_, err := trx.Execute(query, enum.TenantActive, c.TenantID)
		if err != nil {
			return errors.Wrap(err, "failed to activate tenant with id '%d'", c.TenantID)
		}
		return nil
	})
}

func getVerificationByKey(ctx context.Context, q *query.GetVerificationByKey) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		verification := dbEntities.EmailVerification{}

		query := "SELECT id, email, name, key, created_at, verified_at, expires_at, kind, user_id FROM email_verifications WHERE key = $1 AND kind = $2 LIMIT 1"
		err := trx.Get(&verification, query, q.Key, q.Kind)
		if err != nil {
			return errors.Wrap(err, "failed to get email verification by its key")
		}

		q.Result = verification.ToModel()
		return nil
	})
}

func getVerificationByEmailAndCode(ctx context.Context, q *query.GetVerificationByEmailAndCode) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		verification := dbEntities.EmailVerification{}

		query := "SELECT id, email, name, key, created_at, verified_at, expires_at, kind, user_id FROM email_verifications WHERE tenant_id = $1 AND email = $2 AND key = $3 AND kind = $4 LIMIT 1"
		err := trx.Get(&verification, query, tenant.ID, q.Email, q.Code, q.Kind)
		if err != nil {
			return errors.Wrap(err, "failed to get email verification by email and code")
		}

		q.Result = verification.ToModel()
		return nil
	})
}

func saveVerificationKey(ctx context.Context, c *cmd.SaveVerificationKey) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var userID any
		if c.Request.GetUser() != nil {
			userID = c.Request.GetUser().ID
		}

		query := "INSERT INTO email_verifications (tenant_id, email, created_at, expires_at, key, name, kind, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
		_, err := trx.Execute(query, tenant.ID, c.Request.GetEmail(), time.Now(), time.Now().Add(c.Duration), c.Key, c.Request.GetName(), c.Request.GetKind(), userID)
		if err != nil {
			return errors.Wrap(err, "failed to save verification key for kind '%d'", c.Request.GetKind())
		}
		return nil
	})
}

func setKeyAsVerified(ctx context.Context, c *cmd.SetKeyAsVerified) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		query := "UPDATE email_verifications SET verified_at = $1 WHERE tenant_id = $2 AND key = $3 AND verified_at IS NULL"
		_, err := trx.Execute(query, time.Now(), tenant.ID, c.Key)
		if err != nil {
			return errors.Wrap(err, "failed to update verified date of email verification request")
		}
		return nil
	})
}

func createTenant(ctx context.Context, c *cmd.CreateTenant) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		now := time.Now()

		var id int
		err := trx.Get(&id,
			`INSERT INTO tenants (name, subdomain, created_at, cname, invitation, welcome_message, status, is_private, custom_css, logo_bkey, locale, is_email_auth_allowed, is_feed_enabled, prevent_indexing, is_moderation_enabled)
			 VALUES ($1, $2, $3, '', '', '', $4, false, '', '', $5, true, true, true, false)
			 RETURNING id`, c.Name, c.Subdomain, now, c.Status, env.Config.Locale)
		if err != nil {
			return err
		}

		byDomain := &query.GetTenantByDomain{Domain: c.Subdomain}
		err = bus.Dispatch(ctx, byDomain)
		c.Result = byDomain.Result
		return err
	})
}

func getFirstTenant(ctx context.Context, q *query.GetFirstTenant) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		tenant := dbEntities.Tenant{}

		err := trx.Get(&tenant, `
			SELECT id, name, subdomain, cname, invitation, locale, welcome_message, status, is_private, logo_bkey, custom_css, allowed_schemes, is_email_auth_allowed, is_feed_enabled, is_moderation_enabled, prevent_indexing, is_pro
			FROM tenants
			ORDER BY id LIMIT 1
		`)
		if err != nil {
			return errors.Wrap(err, "failed to get first tenant")
		}

		q.Result = tenant.ToModel()
		return nil
	})
}

func getTenantByDomain(ctx context.Context, q *query.GetTenantByDomain) error {
	return using(ctx, func(trx *dbx.Trx, _ *entity.Tenant, _ *entity.User) error {
		tenant := dbEntities.Tenant{}

		err := trx.Get(&tenant, `
			SELECT id, name, subdomain, cname, invitation, locale, welcome_message, status, is_private, logo_bkey, custom_css, allowed_schemes, is_email_auth_allowed, is_feed_enabled, is_moderation_enabled, prevent_indexing, is_pro
			FROM tenants t
			WHERE subdomain = $1 OR subdomain = $2 OR cname = $3
			ORDER BY cname DESC
		`, env.Subdomain(q.Domain), q.Domain, q.Domain)
		if err != nil {
			return errors.Wrap(err, "failed to get tenant with domain '%s'", q.Domain)
		}

		q.Result = tenant.ToModel()
		return nil
	})
}

func getPendingSignUpVerification(ctx context.Context, q *query.GetPendingSignUpVerification) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		verification := dbEntities.EmailVerification{}

		query := `SELECT id, email, name, key, created_at, verified_at, expires_at, kind, user_id 
		          FROM email_verifications 
		          WHERE tenant_id = $1 AND kind = $2 AND verified_at IS NULL
		          ORDER BY created_at DESC 
		          LIMIT 1`
		err := trx.Get(&verification, query, tenant.ID, enum.EmailVerificationKindSignUp)
		if err != nil {
			return errors.Wrap(err, "failed to get pending signup verification for tenant '%d'", tenant.ID)
		}

		q.Result = verification.ToModel()
		return nil
	})
}

func invalidatePreviousSignUpKeys(ctx context.Context, c *cmd.InvalidatePreviousSignUpKeys) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		query := `UPDATE email_verifications 
		          SET expires_at = $1 
		          WHERE tenant_id = $2 AND kind = $3 AND verified_at IS NULL AND expires_at > $1`
		_, err := trx.Execute(query, time.Now(), tenant.ID, enum.EmailVerificationKindSignUp)
		if err != nil {
			return errors.Wrap(err, "failed to invalidate previous signup keys for tenant '%d'", tenant.ID)
		}
		return nil
	})
}
