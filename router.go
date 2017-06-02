package main

import (
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/storage"
	"github.com/labstack/echo"
)

// AppServices holds reference to all Fider services
type AppServices struct {
	OAuth    oauth.Service
	User     storage.User
	Tenant   storage.Tenant
	Idea     storage.Idea
	Settings *models.AppSettings
}

// GetMainEngine returns main HTTP engine
func GetMainEngine(settings *models.AppSettings) *web.Engine {
	r := web.New(settings)

	db, err := dbx.NewWithLogger(r.Logger)
	if err != nil {
		panic(err)
	}
	db.Migrate()

	assets := r.Group("/assets")
	{
		assets.Use(middlewares.OneYearCache())
		assets.Static("/favicon.ico", "favicon.ico")
		assets.Static("/", "dist")
	}

	page := r.Group("")
	{
		page.Use(middlewares.Setup(db))
		page.Use(middlewares.Tenant())
		page.Use(middlewares.AddServices())
		page.Use(middlewares.JwtGetter())
		page.Use(middlewares.JwtSetter())

		public := page.Group("")
		{
			public.Get("/", handlers.Index())
			public.Get("/ideas/:number", handlers.IdeaDetails())
			public.Get("/ideas/:number/*", handlers.IdeaDetails())
			public.Get("/logout", handlers.Logout())
			public.Get("/api/status", handlers.Status(settings))
		}

		private := page.Group("")
		{
			private.Use(middlewares.IsAuthenticated())

			private.Post("/api/ideas", handlers.PostIdea())
			private.Post("/api/ideas/:number/comments", handlers.PostComment())
			private.Post("/api/ideas/:number/status", handlers.SetResponse())
			private.Post("/api/ideas/:number/support", handlers.AddSupporter())
			private.Post("/api/ideas/:number/unsupport", handlers.RemoveSupporter())
		}

		auth := page.Group("/oauth")
		{
			auth.Get("/facebook", handlers.Login(oauth.FacebookProvider))
			auth.Get("/facebook/callback", handlers.OAuthCallback(oauth.FacebookProvider))
			auth.Get("/google", handlers.Login(oauth.GoogleProvider))
			auth.Get("/google/callback", handlers.OAuthCallback(oauth.GoogleProvider))
		}

		admin := page.Group("/admin")
		{
			admin.Use(middlewares.IsAuthenticated())
			admin.Use(middlewares.IsAuthorized(models.RoleMember, models.RoleAdministrator))

			admin.Get("", func(ctx web.Context) error {
				return ctx.Page(echo.Map{})
			})
		}
	}

	return r
}
