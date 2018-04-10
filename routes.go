package main

import (
	"time"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
)

// GetMainEngine returns main HTTP engine
func GetMainEngine(settings *models.SystemSettings) *web.Engine {
	r := web.New(settings)

	r.Worker().Use(middlewares.WorkerSetup(r.Worker().Logger()))

	r.Use(middlewares.Secure())
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
		page.Use(middlewares.CheckTenantPrivacy())

		public := page.Group()
		{
			public.Get("/", handlers.Index())
			public.Get("/api/ideas/search", handlers.SearchIdeas())
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
			private.Get("/settings", handlers.UserSettings())
			private.Get("/notifications", handlers.Notifications())
			private.Get("/notifications/:id", handlers.ReadNotification())
			private.Get("/change-email/verify", handlers.VerifyChangeEmailKey())

			private.Post("/api/ideas", handlers.PostIdea())
			private.Post("/api/ideas/:number", handlers.UpdateIdea())
			private.Post("/api/ideas/:number/comments", handlers.PostComment())
			private.Post("/api/ideas/:number/comments/:id", handlers.UpdateComment())
			private.Post("/api/ideas/:number/status", handlers.SetResponse())
			private.Post("/api/ideas/:number/support", handlers.AddSupporter())
			private.Post("/api/ideas/:number/unsupport", handlers.RemoveSupporter())
			private.Post("/api/ideas/:number/subscribe", handlers.Subscribe())
			private.Post("/api/ideas/:number/unsubscribe", handlers.Unsubscribe())
			private.Post("/api/ideas/:number/tags/:slug", handlers.AssignTag())
			private.Delete("/api/ideas/:number/tags/:slug", handlers.UnassignTag())
			private.Post("/api/user/settings", handlers.UpdateUserSettings())
			private.Post("/api/user/change-email", handlers.ChangeUserEmail())
			private.Post("/api/notifications/read-all", handlers.ReadAllNotifications())
			private.Get("/api/notifications/unread/total", handlers.TotalUnreadNotifications())

			private.Use(middlewares.IsAuthorized(models.RoleAdministrator))

			private.Delete("/api/ideas/:number", handlers.DeleteIdea())
			private.Post("/api/admin/settings/general", handlers.UpdateSettings())
			private.Post("/api/admin/settings/privacy", handlers.UpdatePrivacy())
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
			admin.Get("/admin/privacy", handlers.Page())
			admin.Get("/admin/members", handlers.ManageMembers())
			admin.Get("/admin/tags", handlers.ManageTags())
		}
	}

	return r
}
