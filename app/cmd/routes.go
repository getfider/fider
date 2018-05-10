package cmd

import (
	"time"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
)

func routes(r *web.Engine) *web.Engine {
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
		noTenant.Get("/-/health", handlers.Health())

		noTenant.Post("/api/tenants", handlers.CreateTenant())
		noTenant.Get("/api/tenants/:subdomain/availability", handlers.CheckAvailability())
		noTenant.Get("/signup", handlers.SignUp())

		noTenant.Get("/oauth/facebook", handlers.SignInByOAuth(oauth.FacebookProvider))
		noTenant.Get("/oauth/facebook/callback", handlers.OAuthCallback(oauth.FacebookProvider))
		noTenant.Get("/oauth/google", handlers.SignInByOAuth(oauth.GoogleProvider))
		noTenant.Get("/oauth/google/callback", handlers.OAuthCallback(oauth.GoogleProvider))
		noTenant.Get("/oauth/github", handlers.SignInByOAuth(oauth.GitHubProvider))
		noTenant.Get("/oauth/github/callback", handlers.OAuthCallback(oauth.GitHubProvider))
	}

	r.Use(middlewares.Tenant())

	avatar := r.Group()
	{
		avatar.Use(middlewares.ClientCache(72 * time.Hour))
		avatar.Get("/avatars/:size/:id/:name", handlers.Avatar())
		avatar.Get("/logo/:id", handlers.Logo())
	}

	open := r.Group()
	{
		open.Get("/signup/verify", handlers.VerifySignUpKey())
		open.Use(middlewares.OnlyActiveTenants())
		open.Get("/signin", handlers.SignInPage())
		open.Get("/not-invited", handlers.NotInvitedPage())
		open.Get("/signin/verify", handlers.VerifySignInKey(models.EmailVerificationKindSignIn))
		open.Get("/invite/verify", handlers.VerifySignInKey(models.EmailVerificationKindUserInvitation))
		open.Post("/api/signin/complete", handlers.CompleteSignInProfile())
		open.Post("/api/signin", handlers.SignInByEmail())
	}

	r.Use(middlewares.JwtGetter())
	r.Use(middlewares.JwtSetter())

	page := r.Group()
	{
		page.Use(middlewares.OnlyActiveTenants())
		page.Use(middlewares.CheckTenantPrivacy())

		public := page.Group()
		{
			public.Get("/", handlers.Index())
			public.Get("/api/ideas/search", handlers.SearchIdeas())
			public.Get("/ideas/:number", handlers.IdeaDetails())
			public.Get("/ideas/:number/*all", handlers.IdeaDetails())
			public.Get("/signout", handlers.SignOut())
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

			private.Use(middlewares.IsAuthorized(models.RoleCollaborator, models.RoleAdministrator))

			private.Get("/admin", handlers.GeneralSettingsPage())
			private.Get("/admin/privacy", handlers.Page("Privacy · Site Settings", ""))
			private.Get("/admin/invitations", handlers.Page("Invitations · Site Settings", ""))
			private.Get("/admin/members", handlers.ManageMembers())
			private.Get("/admin/tags", handlers.ManageTags())
			private.Post("/api/admin/invitations/send", handlers.SendInvites())
			private.Post("/api/admin/invitations/sample", handlers.SendSampleInvite())

			private.Use(middlewares.IsAuthorized(models.RoleAdministrator))

			private.Get("/admin/export", handlers.Page("Export · Site Settings", ""))
			private.Get("/admin/export/ideas.csv", handlers.ExportIdeasToCSV())
			private.Delete("/api/ideas/:number", handlers.DeleteIdea())
			private.Post("/api/admin/settings/general", handlers.UpdateSettings())
			private.Post("/api/admin/settings/privacy", handlers.UpdatePrivacy())
			private.Delete("/api/admin/tags/:slug", handlers.DeleteTag())
			private.Post("/api/admin/tags/:slug", handlers.CreateEditTag())
			private.Post("/api/admin/tags", handlers.CreateEditTag())
			private.Post("/api/admin/users/:user_id/role", handlers.ChangeUserRole())
		}
	}

	return r
}
