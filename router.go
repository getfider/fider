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
func GetMainEngine(settings *models.AppSettings, db *dbx.Database) *web.Engine {
	r := web.New(settings)

	assets := r.Group("/assets")
	{
		assets.Use(middlewares.OneYearCache())
		assets.Static("/", "dist")
	}

	public := r.Group("")
	{
		public.Use(middlewares.Setup(db))
		public.Use(middlewares.Tenant())
		public.Use(middlewares.AddServices())
		public.Use(middlewares.JwtGetter())
		public.Use(middlewares.JwtSetter())

		public.Get("/", handlers.Index())
		public.Get("/ideas/:number", handlers.IdeaDetails())
		public.Get("/ideas/:number/*", handlers.IdeaDetails())
		public.Get("/logout", handlers.Logout())
		public.Get("/api/status", handlers.Status(settings))
	}

	private := r.Group("")
	{
		private.Use(middlewares.Setup(db))
		private.Use(middlewares.Tenant())
		private.Use(middlewares.AddServices())
		private.Use(middlewares.JwtGetter())
		private.Use(middlewares.JwtSetter())
		private.Use(middlewares.IsAuthenticated())

		private.Post("/api/ideas", handlers.PostIdea())
		private.Post("/api/ideas/:number/comments", handlers.PostComment())
		private.Post("/api/ideas/:number/status", handlers.SetResponse())
		private.Post("/api/ideas/:number/support", handlers.AddSupporter())
		private.Post("/api/ideas/:number/unsupport", handlers.RemoveSupporter())
	}

	auth := r.Group("/oauth")
	{
		auth.Use(middlewares.Setup(db))
		auth.Use(middlewares.Tenant())
		auth.Use(middlewares.AddServices())

		auth.Get("/facebook", handlers.Login(oauth.FacebookProvider))
		auth.Get("/facebook/callback", handlers.OAuthCallback(oauth.FacebookProvider))
		auth.Get("/google", handlers.Login(oauth.GoogleProvider))
		auth.Get("/google/callback", handlers.OAuthCallback(oauth.GoogleProvider))
	}

	admin := r.Group("/admin")
	{
		admin.Use(middlewares.Setup(db))
		admin.Use(middlewares.Tenant())
		admin.Use(middlewares.AddServices())
		admin.Use(middlewares.JwtGetter())
		admin.Use(middlewares.JwtSetter())
		admin.Use(middlewares.IsAuthenticated())
		admin.Use(middlewares.IsAuthorized(models.RoleMember, models.RoleAdministrator))

		admin.Get("", func(ctx web.Context) error {
			return ctx.Page(echo.Map{})
		})
	}

	return r
}
