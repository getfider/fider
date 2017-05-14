package main

import (
	"net/http"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models"
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

// GetMainEngine returns main HTTP engine
func GetMainEngine(ctx *AppServices) *web.Engine {
	r := web.New(ctx.Settings)

	assets := r.Group("/assets")
	{
		assets.Use(middlewares.OneYearCache())
		assets.Static("/", "dist")
	}

	public := r.Group("")
	{
		public.Use(middlewares.Tenant(ctx.Tenant))
		public.Use(middlewares.JwtGetter(ctx.User))
		public.Use(middlewares.JwtSetter())

		public.Get("/", handlers.Handlers(ctx.Idea).List())
		public.Get("/ideas/:number", handlers.Handlers(ctx.Idea).Details())
		public.Get("/logout", handlers.Logout())
		public.Get("/api/status", handlers.Status(ctx.Settings))
	}

	private := r.Group("")
	{
		private.Use(middlewares.Tenant(ctx.Tenant))
		private.Use(middlewares.JwtGetter(ctx.User))
		private.Use(middlewares.JwtSetter())
		private.Use(middlewares.IsAuthenticated())

		private.Post("/api/ideas", handlers.Handlers(ctx.Idea).PostIdea())
		private.Post("/api/ideas/:number/comments", handlers.Handlers(ctx.Idea).PostComment())
		private.Post("/api/ideas/:number/support", handlers.Handlers(ctx.Idea).AddSupporter())
		private.Post("/api/ideas/:number/unsupport", handlers.Handlers(ctx.Idea).RemoveSupporter())
	}

	auth := r.Group("/oauth")
	{
		auth.Use(middlewares.Tenant(ctx.Tenant))
		authHandlers := handlers.OAuth(ctx.Tenant, ctx.OAuth, ctx.User)

		auth.Get("/facebook", authHandlers.Login(oauth.FacebookProvider))
		auth.Get("/facebook/callback", authHandlers.Callback(oauth.FacebookProvider))
		auth.Get("/google", authHandlers.Login(oauth.GoogleProvider))
		auth.Get("/google/callback", authHandlers.Callback(oauth.GoogleProvider))
	}

	admin := r.Group("/admin")
	{
		admin.Use(middlewares.Tenant(ctx.Tenant))
		admin.Use(middlewares.JwtGetter(ctx.User))
		admin.Use(middlewares.JwtSetter())
		admin.Use(middlewares.IsAuthenticated())
		admin.Use(middlewares.IsAuthorized(models.RoleMember, models.RoleAdministrator))

		admin.Get("", func(ctx web.Context) error {
			return ctx.HTML(http.StatusOK, "Welcome to Admin Page :)")
		})
	}

	return r
}
