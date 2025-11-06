package postgres

import (
	"context"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/services/sqlstore/dbEntities"
)

func getCustomOAuthConfigByProvider(ctx context.Context, q *query.GetCustomOAuthConfigByProvider) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if tenant == nil {
			return app.ErrNotFound
		}

		config := &dbEntities.OAuthConfig{}
		err := trx.Get(config, `
		SELECT id, provider, display_name, status, is_trusted, logo_bkey,
					 client_id, client_secret, authorize_url,
					 profile_url, token_url, scope, json_user_id_path,
					 json_user_name_path, json_user_email_path
		FROM oauth_providers
		WHERE tenant_id = $1 AND provider = $2
		`, tenant.ID, q.Provider)
		if err != nil {
			return err
		}

		q.Result = config.ToModel()
		return nil
	})
}

func listCustomOAuthConfig(ctx context.Context, q *query.ListCustomOAuthConfig) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		if trx == nil || tenant == nil {
			return nil
		}

		configs := []*dbEntities.OAuthConfig{}
		if tenant != nil {
			err := trx.Select(&configs, `
			SELECT id, provider, display_name, status, is_trusted, logo_bkey,
						 client_id, client_secret, authorize_url,
						 profile_url, token_url, scope, json_user_id_path,
						 json_user_name_path, json_user_email_path
			FROM oauth_providers
			WHERE tenant_id = $1
			ORDER BY id`, tenant.ID)
			if err != nil {
				return err
			}
		}

		q.Result = make([]*entity.OAuthConfig, len(configs))
		for i, config := range configs {
			q.Result[i] = config.ToModel()
		}
		return nil
	})
}

func saveCustomOAuthConfig(ctx context.Context, c *cmd.SaveCustomOAuthConfig) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		var err error

		if c.Logo.Remove {
			c.Logo.BlobKey = ""
		}

		if c.ID == 0 {
			query := `INSERT INTO oauth_providers (
				tenant_id, provider, display_name, status, is_trusted,
				client_id, client_secret, authorize_url,
				profile_url, token_url, scope, json_user_id_path,
				json_user_name_path, json_user_email_path, logo_bkey
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING id`

			err = trx.Get(&c.ID, query, tenant.ID, c.Provider,
				c.DisplayName, c.Status, c.IsTrusted, c.ClientID, c.ClientSecret,
				c.AuthorizeURL, c.ProfileURL, c.TokenURL,
				c.Scope, c.JSONUserIDPath, c.JSONUserNamePath,
				c.JSONUserEmailPath, c.Logo.BlobKey)
		} else {
			query := `
				UPDATE oauth_providers 
				SET display_name = $3, status = $4, client_id = $5, client_secret = $6, 
						authorize_url = $7, profile_url = $8, token_url = $9, scope = $10, 
						json_user_id_path = $11, json_user_name_path = $12, json_user_email_path = $13,
						logo_bkey = $14, is_trusted = $15
			WHERE tenant_id = $1 AND id = $2`

			_, err = trx.Execute(query, tenant.ID, c.ID,
				c.DisplayName, c.Status, c.ClientID, c.ClientSecret,
				c.AuthorizeURL, c.ProfileURL, c.TokenURL,
				c.Scope, c.JSONUserIDPath, c.JSONUserNamePath,
				c.JSONUserEmailPath, c.Logo.BlobKey, c.IsTrusted)
		}

		if err != nil {
			return errors.Wrap(err, "failed to save OAuth Provider")
		}

		return nil
	})
}
