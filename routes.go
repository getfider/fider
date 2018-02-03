package main

import (
	"time"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
)

// GetMainEngine returns main HTTP engine
func GetMainEngine(settings *models.AppSettings) *web.Engine {
	r := web.New(settings)

	r.Worker().Use(middlewares.WorkerSetup(r.Worker().Logger()))

	r.Use(middlewares.Compress())

	assets := r.Group()
	{
		assets.Use(middlewares.ClientCache(365 * 24 * time.Hour))
		assets.Static("/favicon.ico", "favicon.ico")
		assets.Static("/assets/*filepath", "dist")
	}

	r.Use(middlewares.WebSetup(r.Logger()))

	noTenant := r.Group()
	{
		noTenant.Post("/api/tenants", handlers.CreateTenant())
		noTenant.Get("/api/tenants/:subdomain/availability", handlers.CheckAvailability())
		noTenant.Get("/signup", handlers.SignUp())

		noTenant.Get("/oauth/facebook", handlers.SignIn(oauth.FacebookProvider))
		noTenant.Get("/oauth/facebook/callback", handlers.OAuthCallback(oauth.FacebookProvider))
		noTenant.Get("/oauth/google", handlers.SignIn(oauth.GoogleProvider))
		noTenant.Get("/oauth/google/callback", handlers.OAuthCallback(oauth.GoogleProvider))
		noTenant.Get("/oauth/github", handlers.SignIn(oauth.GitHubProvider))
		noTenant.Get("/oauth/github/callback", handlers.OAuthCallback(oauth.GitHubProvider))
	}

	r.Use(middlewares.Tenant())

	avatar := r.Group()
	{
		avatar.Use(middlewares.ClientCache(72 * time.Hour))
		avatar.Get("/avatars/:size/:id/:name", handlers.Avatar())
	}

	verify := r.Group()
	{
		verify.Get("/signup/verify", handlers.VerifySignUpKey())
	}

	page := r.Group()
	{
		page.Use(middlewares.JwtGetter())
		page.Use(middlewares.JwtSetter())
		page.Use(middlewares.OnlyActiveTenants())

		public := page.Group()
		{
			public.Get("/", handlers.Index())
			public.Get("/ideas/:number", handlers.IdeaDetails())
			public.Get("/ideas/:number/*all", handlers.IdeaDetails())
			public.Get("/signout", handlers.SignOut())
			public.Get("/signin/verify", handlers.VerifySignInKey(models.EmailVerificationKindSignIn))
			public.Get("/api/status", handlers.Status(settings))
			public.Post("/api/signin/complete", handlers.CompleteSignInProfile())
			public.Post("/api/signin", handlers.SignInByEmail())
		}

		private := page.Group()
		{
			private.Use(middlewares.IsAuthenticated())
			private.Get("/settings", handlers.Page())
			private.Get("/change-email/verify", handlers.VerifyChangeEmailKey())

			private.Get("/api/ideas", handlers.GetIdeas())
			private.Post("/api/ideas", handlers.PostIdea())
			private.Post("/api/ideas/:number", handlers.UpdateIdea())
			private.Post("/api/ideas/:number/comments", handlers.PostComment())
			private.Post("/api/ideas/:number/status", handlers.SetResponse())
			private.Post("/api/ideas/:number/support", handlers.AddSupporter())
			private.Post("/api/ideas/:number/unsupport", handlers.RemoveSupporter())
			private.Post("/api/ideas/:number/tags/:slug", handlers.AssignTag())
			private.Delete("/api/ideas/:number/tags/:slug", handlers.UnassignTag())
			private.Post("/api/user/settings", handlers.UpdateUserSettings())
			private.Post("/api/user/change-email", handlers.ChangeUserEmail())

			private.Use(middlewares.IsAuthorized(models.RoleAdministrator))

			private.Post("/api/admin/settings", handlers.UpdateSettings())
			private.Delete("/api/admin/tags/:slug", handlers.DeleteTag())
			private.Post("/api/admin/tags/:slug", handlers.CreateEditTag())
			private.Post("/api/admin/tags", handlers.CreateEditTag())
			private.Post("/api/admin/users/:user_id/role", handlers.ChangeUserRole())
		}

		admin := page.Group()
		{
			admin.Use(middlewares.IsAuthenticated())
			admin.Use(middlewares.IsAuthorized(models.RoleCollaborator, models.RoleAdministrator))

			admin.Get("/admin", handlers.Page())
			admin.Get("/admin/members", handlers.ManageMembers())
			admin.Get("/admin/tags", handlers.ManageTags())
		}
	}

	if env.IsDevelopment() {
		debug := r.Group()
		{
			debug.Get("/debug/stats", handlers.RuntimeStats())
		}
	}

	return r
}
