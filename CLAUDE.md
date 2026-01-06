# CLAUDE.md

This file provides guidance to Claude Code when working with the Fider codebase.

## Repository Structure

**Backend (Go):**

- `app/handlers/` - HTTP request handlers
- `app/models/` - Data models (entity, cmd, query, action, dto)
- `app/services/` - Business logic and external service integrations
- `app/pkg/bus/` - Service registry and dispatch system
- `app/cmd/routes.go` - All HTTP routes defined here
- `migrations/` - Database migrations (numbered SQL files)

**Frontend (React/TypeScript):**

- `public/pages/` - Page components (lazy-loaded)
- `public/components/` - Reusable UI components
- `public/services/` - Client-side services and API calls
- `public/hooks/` - Custom React hooks
- `public/assets/styles/` - SCSS styles with utility classes

**Configuration:**

- `.env` - Local environment config (copy from `.example.env`)
- `Makefile` - Build and development commands

## Local Development

All commands are defined in the Makefile.

**Development:**

- `make watch` - Hot reload for both server and UI (use this for active development)
- `make migrate` - Run database migrations

**Testing:**

- `make test` - Run all tests (both Go and Jest)
- `make lint` - Lint both server and UI code

Essential: Run "make lint" and "make test" when you've completed any changes, to check the formatting and tests.

**Building:**

- `make build` - Build production-ready binaries and assets

## Code Patterns & Examples

### Backend: Adding a New API Endpoint

**1. Define the route** in `app/cmd/routes.go`:

```go
// Public API
engine.Get("/api/v1/posts", handlers.SearchPosts())

// Authenticated API (requires login)
engine.Get("/api/v1/user/settings", middlewares.IsAuthenticated(), handlers.UserSettings())

// Admin only
engine.Post("/api/v1/admin/members", middlewares.IsAuthenticated(), middlewares.IsAuthorized(enum.RoleAdministrator), handlers.CreateMember())
```

**2. Create the handler** in `app/handlers/`:

```go
package handlers

import (
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/web"
)

// SearchPosts returns posts matching criteria
func SearchPosts() web.HandlerFunc {
	return func(c *web.Context) error {
		q := &query.SearchPosts{
			Query: c.QueryParam("query"),
			Limit: 10,
		}

		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}

		return c.Ok(q.Result)
	}
}
```

**3. Define the query** in `app/models/query/`:

```go
package query

import "github.com/getfider/fider/app/models/entity"

type SearchPosts struct {
	Query  string
	Limit  int
	Result []*entity.Post
}
```

**4. Implement the service** in `app/services/`:

```go
package postgres

import (
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/dbx"
)

func searchPosts(ctx context.Context, q *query.SearchPosts) error {
	return using(ctx, func(trx *dbx.Trx, tenant *entity.Tenant, user *entity.User) error {
		q.Result = []*entity.Post{}
		return trx.Select(&q.Result,
			"SELECT * FROM posts WHERE title ILIKE $1 LIMIT $2",
			"%"+q.Query+"%", q.Limit,
		)
	})
}

func init() {
	bus.Register(Service{})
}

func (s Service) Name() string {
	return "Postgres"
}

func (s Service) Category() string {
	return "sqlstore"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	bus.AddHandler(searchPosts)
}
```

### Frontend: Calling the API from React

```typescript
// In a React component
import { http } from "@fider/services"

export const PostSearch: React.FC = () => {
  const [posts, setPosts] = React.useState([])
  const [query, setQuery] = React.useState("")

  const handleSearch = async () => {
    const result = await http.get<Post[]>(`/api/v1/posts?query=${query}`)
    if (result.ok) {
      setPosts(result.data)
    }
  }

  return (
    <div className="c-post-search">
      <input className="c-post-search__input" value={query} onChange={(e) => setQuery(e.target.value)} />
      <Button onClick={handleSearch}>Search</Button>
    </div>
  )
}
```

### Model Naming Conventions

Use these strict prefixes in `app/models/`:

```go
// entity.* - Database tables
type entity.User struct {
    ID    int
    Name  string
    Email string
}

// action.* - User input for POST/PUT/PATCH (maps 1:1 to commands)
type action.CreatePost struct {
    Title       string
    Description string
}

// cmd.* - Commands to execute (potentially return values)
type cmd.SendEmail struct {
    To      string
    Subject string
    Body    string
}

// query.* - Queries to fetch data
type query.GetPostByID struct {
    PostID int
    Result *entity.Post
}

// dto.* - Data transfer between packages
type dto.PostInfo struct {
    Title     string
    AuthorName string
}
```

### CSS Conventions

Most pages have a corresponding scss file. For example Home.Page.tsx imports Home.Page.scss which contains styles for various components in that page. These follow a BEM style:

```scss
// Page ID: p-<page_name>
#p-home {
  // Component: c-<component>
  .c-post-list {
    // Element: c-<component>__<element>
    &__item {
      padding: 1rem;
    }

    &__title {
      font-size: 1.2rem;
    }

    // Modifier: c-<component>--<state>
    &--loading {
      opacity: 0.5;
    }
  }
}
```

However, more recently, Utility classes were added to the codebase, to make it possible to apply classes in a utility style, similar to tailwind. If possible favour the utility classes over adding new styles for components and pages:

```scss
// Utility classes (no prefix) - defined in public/assets/styles/utility/
.mt-2 {
  margin-top: 0.5rem;
}
.flex {
  display: flex;
}
```

Important: Do not just apply tailwind classes to new code, tailwind is not in use in this project. Check what is available in the utility scss files.

## Common Tasks

### Creating a Database Migration

"up migrations" only.

```bash
# Create new migration file in migrations/ directory
# Format: YYYYMMDDHHMMSS_description.sql
# Example: migrations/202601041200_add_user_preferences.sql

CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    theme VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

```

Run with: `make migrate`

### Adding a New Page Component

```typescript
// public/pages/MyNewPage/MyNewPage.page.tsx
import React from "react"
import { Page } from "@fider/components"

interface MyNewPageProps {
  // Props here
}

const MyNewPage: React.FC<MyNewPageProps> = (props) => {
  return (
    <Page id="p-my-new-page">
      <h1>My New Page</h1>
    </Page>
  )
}

export default MyNewPage
```

```typescript
// public/pages/MyNewPage/index.ts
export { default as MyNewPage } from "./MyNewPage.page"
```

Register route in `app/handlers/` and lazy-load in `public/index.tsx`.

### Using the Bus System

The bus handles service registration and dispatch:

```go
// Register a service
bus.Register(MyService{})

// Dispatch a query
q := &query.GetUserByID{UserID: 123}
if err := bus.Dispatch(ctx, q); err != nil {
    return err
}
user := q.Result

// Dispatch a command
c := &cmd.SendEmail{To: "user@example.com", Subject: "Welcome"}
if err := bus.Dispatch(ctx, c); err != nil {
    return err
}
```

## Architecture Patterns

**Backend (Go):**

- CQRS: Separate commands (write) and queries (read)
- Bus system: Dependency injection and service dispatch
- Middleware chain: Auth, tenant resolution, CORS applied via routes
- SQL-based data access: Direct SQL queries, no ORM

**Frontend (React):**

- SSR Support: Server-side rendering with React 18 hydration
- i18n: LinguiJS for translations
- Code splitting: Lazy-loaded page components
- Type safety: Full TypeScript coverage

## Troubleshooting

**Database connection errors:**

- Ensure Docker is running: `docker compose ps`
- Check `.env` has correct `DATABASE_URL`
- Run migrations: `make migrate`

**Build failures:**

- Clear build cache: `make clean && make build`
- Check Go version: `go version` (need 1.22+)
- Check Node version: `node --version` (need 21/22)

**Tests failing:**

- Ensure test database is migrated (happens automatically with `make test`)
- Check `.test.env` configuration
- Run specific test: `go test ./app/handlers -v -run TestName`

**Port conflicts:**

- Default ports: 3000 (app), 5432 (postgres), 8025 (mailhog)
- Change in `.env` if needed

## Development Tips

- Use `make watch` for active development - it handles hot reload for both frontend and backend
- MailHog captures all emails at http://localhost:8025
- Frontend changes: Webpack rebuilds automatically
- Backend changes: Air restarts server automatically
- All routes are centralized in `app/cmd/routes.go` - start there when tracing request flow
- Bus handlers are registered in service `Init()` methods
- Utility CSS classes are in `public/assets/styles/utility/` - use these before adding custom styles
