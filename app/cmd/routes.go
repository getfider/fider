package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/web"
)

func routes(r *web.Engine) *web.Engine {
	r.Worker().Use(middlewares.WorkerSetup())

	r.Use(middlewares.Secure())
	r.Use(middlewares.Compress())

	files := r.Group()
	{
		files.Get("/robots.txt", handlers.RobotsTXT())
		files.Get("/privacy", handlers.LegalPage("Privacy Policy", "privacy.md"))
		files.Get("/terms", handlers.LegalPage("Terms of Service", "terms.md"))
	}

	assets := r.Group()
	{
		assets.Use(middlewares.CORS())
		assets.Use(middlewares.ClientCache(365 * 24 * time.Hour))
		assets.Static("/favicon.ico", "favicon.ico")
		assets.Static("/assets/*filepath", "dist")
	}

	r.Use(middlewares.WebSetup())
	r.Use(middlewares.Tenant())

	noTenant := r.Group()
	{
		noTenant.Get("/-/health", handlers.Health())

		noTenant.Post("/api/tenants", handlers.CreateTenant())
		noTenant.Get("/api/tenants/:subdomain/availability", handlers.CheckAvailability())
		noTenant.Get("/signup", handlers.SignUp())

		noTenant.Get("/oauth/:provider", handlers.SignInByOAuth())
		noTenant.Get("/oauth/:provider/callback", handlers.OAuthCallback())
	}

	r.Use(middlewares.RequireTenant())

	tenantAssets := r.Group()
	{
		tenantAssets.Use(middlewares.ClientCache(72 * time.Hour))
		tenantAssets.Get("/avatars/:size/:id/:name", handlers.Avatar())
		tenantAssets.Get("/images/:size/:id", handlers.ViewUploadedImage())
		tenantAssets.Get("/custom/:md5.css", func(c web.Context) error {
			return c.Blob(http.StatusOK, "text/css", []byte(c.Tenant().CustomCSS))
		})
	}

	open := r.Group()
	{
		open.Get("/-/ui", handlers.Page("UI Toolkit", "A preview of Fider UI elements"))
		open.Get("/signup/verify", handlers.VerifySignUpKey())
		open.Get("/oauth/:provider/token", handlers.OAuthToken())
		open.Get("/oauth/:provider/echo", handlers.OAuthEcho())
		open.Use(middlewares.OnlyActiveTenants())
		open.Get("/signin", handlers.SignInPage())
		open.Get("/not-invited", handlers.NotInvitedPage())
		open.Get("/signin/verify", handlers.VerifySignInKey(models.EmailVerificationKindSignIn))
		open.Get("/invite/verify", handlers.VerifySignInKey(models.EmailVerificationKindUserInvitation))
		open.Post("/api/signin/complete", handlers.CompleteSignInProfile())
		open.Post("/api/signin", handlers.SignInByEmail())
	}

	r.Use(middlewares.User())

	page := r.Group()
	{
		page.Use(middlewares.OnlyActiveTenants())
		page.Use(middlewares.CheckTenantPrivacy())

		public := page.Group()
		{
			public.Get("/", handlers.Index())
			public.Get("/api/v1/posts", apiv1.SearchPosts())
			public.Get("/posts/:number", handlers.PostDetails())
			public.Get("/posts/:number/*all", handlers.PostDetails())
			public.Get("/signout", handlers.SignOut())

			/* This is a temporary redirect and should be removed in the future */
			public.Get("/ideas/:number", func(c web.Context) error {
				return c.Redirect(strings.Replace(c.Request.URL.Path, "/ideas/", "/posts/", 1))
			})
			public.Get("/ideas/:number/*all", func(c web.Context) error {
				return c.Redirect(strings.Replace(c.Request.URL.Path, "/ideas/", "/posts/", 1))
			})
			/* This is a temporary redirect and should be removed in the future */
		}

		private := page.Group()
		{
			private.Use(middlewares.IsAuthenticated())
			private.Get("/settings", handlers.UserSettings())
			private.Get("/notifications", handlers.Notifications())
			private.Get("/notifications/:id", handlers.ReadNotification())
			private.Get("/change-email/verify", handlers.VerifyChangeEmailKey())

			private.Post("/api/v1/posts", apiv1.CreatePost())
			private.Post("/api/posts/:number", handlers.UpdatePost())
			private.Post("/api/posts/:number/comments", handlers.PostComment())
			private.Post("/api/posts/:number/comments/:id", handlers.UpdateComment())
			private.Post("/api/posts/:number/status", handlers.SetResponse())
			private.Post("/api/posts/:number/support", handlers.AddSupporter())
			private.Post("/api/posts/:number/unsupport", handlers.RemoveSupporter())
			private.Post("/api/posts/:number/subscribe", handlers.Subscribe())
			private.Post("/api/posts/:number/unsubscribe", handlers.Unsubscribe())
			private.Post("/api/posts/:number/tags/:slug", handlers.AssignTag())
			private.Delete("/api/posts/:number/tags/:slug", handlers.UnassignTag())
			private.Delete("/api/user", handlers.DeleteUser())
			private.Post("/api/user/settings", handlers.UpdateUserSettings())
			private.Post("/api/user/change-email", handlers.ChangeUserEmail())
			private.Post("/api/notifications/read-all", handlers.ReadAllNotifications())
			private.Get("/api/notifications/unread/total", handlers.TotalUnreadNotifications())

			private.Use(middlewares.IsAuthorized(models.RoleCollaborator, models.RoleAdministrator))

			private.Get("/admin", handlers.GeneralSettingsPage())
			private.Get("/admin/advanced", handlers.AdvancedSettingsPage())
			private.Get("/admin/privacy", handlers.Page("Privacy · Site Settings", ""))
			private.Get("/admin/invitations", handlers.Page("Invitations · Site Settings", ""))
			private.Get("/admin/members", handlers.ManageMembers())
			private.Get("/admin/tags", handlers.ManageTags())
			private.Get("/admin/authentication", handlers.ManageAuthentication())
			private.Get("/api/admin/oauth/:provider", handlers.GetOAuthConfig())
			private.Post("/api/admin/invitations/send", handlers.SendInvites())
			private.Post("/api/admin/invitations/sample", handlers.SendSampleInvite())

			private.Use(middlewares.IsAuthorized(models.RoleAdministrator))

			private.Get("/admin/export", handlers.Page("Export · Site Settings", ""))
			private.Get("/admin/export/posts.csv", handlers.ExportPostsToCSV())
			private.Delete("/api/posts/:number", handlers.DeletePost())
			private.Post("/api/admin/settings/general", handlers.UpdateSettings())
			private.Post("/api/admin/settings/advanced", handlers.UpdateAdvancedSettings())
			private.Post("/api/admin/settings/privacy", handlers.UpdatePrivacy())
			private.Delete("/api/admin/tags/:slug", handlers.DeleteTag())
			private.Post("/api/admin/tags/:slug", handlers.CreateEditTag())
			private.Post("/api/admin/tags", handlers.CreateEditTag())
			private.Post("/api/admin/oauth", handlers.SaveOAuthConfig())
			private.Post("/api/admin/users/:user_id/role", handlers.ChangeUserRole())
		}
	}

	return r
}
