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

	r.Get("/-/health", handlers.Health())
	r.Get("/robots.txt", handlers.RobotsTXT())
	r.Get("/privacy", handlers.LegalPage("Privacy Policy", "privacy.md"))
	r.Get("/terms", handlers.LegalPage("Terms of Service", "terms.md"))

	assets := r.Group()
	{
		assets.Use(middlewares.CORS())
		assets.Use(middlewares.ClientCache(365 * 24 * time.Hour))
		assets.Static("/favicon.ico", "favicon.ico")
		assets.Static("/assets/*filepath", "dist")
	}

	r.Use(middlewares.WebSetup())
	r.Use(middlewares.Tenant())

	r.Post("/_api/tenants", handlers.CreateTenant())
	r.Get("/_api/tenants/:subdomain/availability", handlers.CheckAvailability())
	r.Get("/signup", handlers.SignUp())
	r.Get("/oauth/:provider", handlers.SignInByOAuth())
	r.Get("/oauth/:provider/callback", handlers.OAuthCallback())

	//From this step, a Tenant is required (regardless of status)
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

	r.Get("/-/ui", handlers.Page("UI Toolkit", "A preview of Fider UI elements"))
	r.Get("/signup/verify", handlers.VerifySignUpKey())
	r.Get("/oauth/:provider/token", handlers.OAuthToken())
	r.Get("/oauth/:provider/echo", handlers.OAuthEcho())

	//From this step, a only active Tenants are allowed
	r.Use(middlewares.OnlyActiveTenants())

	r.Get("/signin", handlers.SignInPage())
	r.Get("/not-invited", handlers.NotInvitedPage())
	r.Get("/signin/verify", handlers.VerifySignInKey(models.EmailVerificationKindSignIn))
	r.Get("/invite/verify", handlers.VerifySignInKey(models.EmailVerificationKindUserInvitation))
	r.Post("/_api/signin/complete", handlers.CompleteSignInProfile())
	r.Post("/_api/signin", handlers.SignInByEmail())

	//From this step, the User might be available
	r.Use(middlewares.User())

	//From this step, block if it's private tenant with unauthenticated user
	r.Use(middlewares.CheckTenantPrivacy())

	r.Get("/", handlers.Index())
	r.Get("/api/v1/posts", apiv1.SearchPosts())
	r.Get("/posts/:number", handlers.PostDetails())
	r.Get("/posts/:number/*all", handlers.PostDetails())
	r.Get("/signout", handlers.SignOut())

	/*
	** This is a temporary redirect and should be removed in the future
	** START
	 */
	r.Get("/ideas/:number", func(c web.Context) error {
		return c.Redirect(strings.Replace(c.Request.URL.Path, "/ideas/", "/posts/", 1))
	})
	r.Get("/ideas/:number/*all", func(c web.Context) error {
		return c.Redirect(strings.Replace(c.Request.URL.Path, "/ideas/", "/posts/", 1))
	})
	/*
	** END
	 */

	//From this step, a User is required
	r.Use(middlewares.IsAuthenticated())

	r.Get("/settings", handlers.UserSettings())
	r.Get("/notifications", handlers.Notifications())
	r.Get("/notifications/:id", handlers.ReadNotification())
	r.Get("/change-email/verify", handlers.VerifyChangeEmailKey())

	r.Post("/api/v1/posts", apiv1.CreatePost())
	r.Put("/api/v1/posts/:number", apiv1.UpdatePost())
	r.Post("/api/v1/posts/:number/comments", apiv1.PostComment())
	r.Put("/api/v1/posts/:number/comments/:id", apiv1.UpdateComment())
	r.Put("/api/v1/posts/:number/status", apiv1.SetResponse())
	r.Post("/api/v1/posts/:number/support", apiv1.AddSupporter())
	r.Delete("/api/v1/posts/:number/support", apiv1.RemoveSupporter())
	r.Post("/api/v1/posts/:number/subscription", apiv1.Subscribe())
	r.Delete("/api/v1/posts/:number/subscription", apiv1.Unsubscribe())
	r.Post("/api/v1/posts/:number/tags/:slug", apiv1.AssignTag())
	r.Delete("/api/v1/posts/:number/tags/:slug", apiv1.UnassignTag())
	r.Delete("/_api/user", handlers.DeleteUser())
	r.Post("/_api/user/settings", handlers.UpdateUserSettings())
	r.Post("/_api/user/change-email", handlers.ChangeUserEmail())
	r.Post("/_api/notifications/read-all", handlers.ReadAllNotifications())
	r.Get("/_api/notifications/unread/total", handlers.TotalUnreadNotifications())

	//From this step, only Collaborators and Administrators are allowed
	r.Use(middlewares.IsAuthorized(models.RoleCollaborator, models.RoleAdministrator))

	r.Get("/admin", handlers.GeneralSettingsPage())
	r.Get("/admin/advanced", handlers.AdvancedSettingsPage())
	r.Get("/admin/privacy", handlers.Page("Privacy · Site Settings", ""))
	r.Get("/admin/invitations", handlers.Page("Invitations · Site Settings", ""))
	r.Get("/admin/members", handlers.ManageMembers())
	r.Get("/admin/tags", handlers.ManageTags())
	r.Get("/admin/authentication", handlers.ManageAuthentication())
	r.Get("/_api/admin/oauth/:provider", handlers.GetOAuthConfig())
	r.Post("/api/v1/invitations/send", apiv1.SendInvites())
	r.Post("/api/v1/invitations/sample", apiv1.SendSampleInvite())

	//From this step, only Administrators are allowed
	r.Use(middlewares.IsAuthorized(models.RoleAdministrator))

	r.Get("/admin/export", handlers.Page("Export · Site Settings", ""))
	r.Get("/admin/export/posts.csv", handlers.ExportPostsToCSV())
	r.Delete("/api/v1/posts/:number", apiv1.DeletePost())
	r.Post("/_api/admin/settings/general", handlers.UpdateSettings())
	r.Post("/_api/admin/settings/advanced", handlers.UpdateAdvancedSettings())
	r.Post("/_api/admin/settings/privacy", handlers.UpdatePrivacy())
	r.Post("/api/v1/tags", apiv1.CreateEditTag())
	r.Put("/api/v1/tags/:slug", apiv1.CreateEditTag())
	r.Delete("/api/v1/tags/:slug", apiv1.DeleteTag())
	r.Post("/_api/admin/oauth", handlers.SaveOAuthConfig())
	r.Post("/api/v1/roles/:role/users", handlers.ChangeUserRole())

	return r
}
