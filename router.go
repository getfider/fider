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
	"github.com/getfider/fider/app/storage"
)

// AppServices holds reference to all Fider services
type AppServices struct {
	OAuth    oauth.Service
	User     storage.User
	Tenant   storage.Tenant
	Idea     storage.Idea
	Settings *models.AppSettings
}

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

	db, err := dbx.NewWithLogger(r.Logger)
	if err != nil {
		panic(err)
	}
	db.Migrate()
	emailer := initEmailer()

	assets := r.Group("/assets")
	{
		assets.Use(middlewares.OneYearCache())
		assets.Get("/avatars/:size/:name", handlers.LetterAvatar())
		assets.Static("/favicon.ico", "favicon.ico")
		assets.Static("/", "dist")
	}

	signup := r.Group("")
	{
		signup.Use(middlewares.Setup(db, emailer))
		signup.Use(middlewares.AddServices())

		signup.Post("/api/tenants", handlers.CreateTenant())
		signup.Get("/api/tenants/:subdomain/availability", handlers.CheckAvailability())
		signup.Get("/signup", handlers.SignUp())
	}

	auth := r.Group("/oauth")
	{
		auth.Use(middlewares.Setup(db, emailer))
		auth.Use(middlewares.AddServices())

		auth.Get("/facebook", handlers.SignIn(oauth.FacebookProvider))
		auth.Get("/facebook/callback", handlers.OAuthCallback(oauth.FacebookProvider))
		auth.Get("/google", handlers.SignIn(oauth.GoogleProvider))
		auth.Get("/google/callback", handlers.OAuthCallback(oauth.GoogleProvider))
		auth.Get("/github", handlers.SignIn(oauth.GitHubProvider))
		auth.Get("/github/callback", handlers.OAuthCallback(oauth.GitHubProvider))
	}

	page := r.Group("")
	{
		page.Use(middlewares.Setup(db, emailer))
		page.Use(middlewares.Tenant())
		page.Use(middlewares.AddServices())
		page.Use(middlewares.JwtGetter())
		page.Use(middlewares.JwtSetter())

		public := page.Group("")
		{
			public.Get("/", handlers.Index())
			public.Get("/ideas/:number", handlers.IdeaDetails())
			public.Get("/ideas/:number/*", handlers.IdeaDetails())
			public.Get("/signout", handlers.SignOut())
			public.Get("/signin/verify", handlers.VerifySignInKey())
			public.Get("/api/status", handlers.Status(settings))
			public.Post("/api/signin/complete", handlers.CompleteSignInProfile())
			public.Post("/api/signin", handlers.SignInByEmail())
		}

		private := page.Group("")
		{
			private.Use(middlewares.IsAuthenticated())
			private.Get("/settings", func(ctx web.Context) error {
				return ctx.Page(web.Map{})
			})

			private.Post("/api/ideas", handlers.PostIdea())
			private.Post("/api/ideas/:number/comments", handlers.PostComment())
			private.Post("/api/ideas/:number/status", handlers.SetResponse())
			private.Post("/api/ideas/:number/support", handlers.AddSupporter())
			private.Post("/api/ideas/:number/unsupport", handlers.RemoveSupporter())
			private.Post("/api/user/settings", handlers.UpdateUserSettings())

			private.Use(middlewares.IsAuthorized(models.RoleMember, models.RoleAdministrator))

			private.Post("/api/admin/settings", handlers.UpdateSettings())
		}

		admin := page.Group("/admin")
		{
			admin.Use(middlewares.IsAuthenticated())
			admin.Use(middlewares.IsAuthorized(models.RoleMember, models.RoleAdministrator))

			admin.Get("", func(ctx web.Context) error {
				return ctx.Page(web.Map{})
			})
		}
	}

	return r
}
