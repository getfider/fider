package main

import (
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/email"
	"github.com/getfider/fider/app/pkg/email/mailgun"
	"github.com/getfider/fider/app/pkg/email/smtp"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
)

func initEmailer() email.Sender {
	if env.IsDefined("MAILGUN_API") {
		return mailgun.NewSender(env.MustGet("MAILGUN_DOMAIN"), env.MustGet("MAILGUN_API"))
	}
	return smtp.NewSender(
		env.MustGet("SMTP_HOST"),
		env.MustGet("SMTP_PORT"),
		env.MustGet("SMTP_USERNAME"),
		env.MustGet("SMTP_PASSWORD"),
	)
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(settings *models.AppSettings) *web.Engine {
	r := web.New(settings)

	db, err := dbx.NewWithLogger(r.Logger())
	if err != nil {
		panic(err)
	}
	db.Migrate()
	emailer := initEmailer()

	r.Use(middlewares.Compress())

	assets := r.Group()
	{
		assets.Use(middlewares.OneYearCache())
		assets.Get("/avatars/:size/:name", handlers.LetterAvatar())
		assets.Static("/favicon.ico", "favicon.ico")
		assets.Static("/assets/*filepath", "dist")
	}

	noTenant := r.Group()
	{
		noTenant.Use(middlewares.Setup(db, emailer))
		noTenant.Use(middlewares.AddServices())

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

	verify := r.Group()
	{
		verify.Use(middlewares.Setup(db, emailer))
		verify.Use(middlewares.Tenant())
		verify.Use(middlewares.AddServices())
		verify.Get("/signup/verify", handlers.VerifySignUpKey())
	}

	page := r.Group()
	{
		page.Use(middlewares.Setup(db, emailer))
		page.Use(middlewares.Tenant())
		page.Use(middlewares.AddServices())
		page.Use(middlewares.JwtGetter())
		page.Use(middlewares.JwtSetter())
		page.Use(middlewares.OnlyActiveTenants())

		public := page.Group()
		{
			public.Get("/", handlers.Index())
			public.Get("/ideas/:number", handlers.IdeaDetails())
			public.Get("/ideas/:number/*all", handlers.IdeaDetails())
			public.Get("/signout", handlers.SignOut())
			public.Get("/signin/verify", handlers.VerifySignInKey())
			public.Get("/api/status", handlers.Status(settings))
			public.Post("/api/signin/complete", handlers.CompleteSignInProfile())
			public.Post("/api/signin", handlers.SignInByEmail())
		}

		private := page.Group()
		{
			private.Use(middlewares.IsAuthenticated())
			private.Get("/settings", handlers.Page())

			private.Post("/api/ideas", handlers.PostIdea())
			private.Post("/api/ideas/:number/comments", handlers.PostComment())
			private.Post("/api/ideas/:number/status", handlers.SetResponse())
			private.Post("/api/ideas/:number/support", handlers.AddSupporter())
			private.Post("/api/ideas/:number/unsupport", handlers.RemoveSupporter())
			private.Post("/api/user/settings", handlers.UpdateUserSettings())

			private.Use(middlewares.IsAuthorized(models.RoleAdministrator))

			private.Post("/api/admin/settings", handlers.UpdateSettings())
			private.Post("/api/admin/users/:user_id/role", handlers.ChangeUserRole())
		}

		admin := page.Group()
		{
			admin.Use(middlewares.IsAuthenticated())
			admin.Use(middlewares.IsAuthorized(models.RoleCollaborator, models.RoleAdministrator))

			admin.Get("/admin", handlers.Page())
			admin.Get("/admin/members", handlers.ManageMembers())
		}
	}

	return r
}
