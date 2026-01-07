# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Building

- `make build` - Build both server and UI
- `make build-server` - Build Go server binary
- `make build-ui` - Build UI assets with webpack
- `make build-ssr` - Build SSR script and compile locales

### Running

- `make run` - Run Fider server (requires build first)
- `make watch` - Start both server and UI in watch mode (recommended for development)
- `make watch-server` - Run server in watch mode with air
- `make watch-ui` - Run UI in watch mode with webpack
- `make migrate` - Run database migrations

### Testing

- `make test` - Run both server and UI tests
- `make test-server` - Run Go server tests (includes migration)
- `make test-ui` - Run Jest tests for React components
- `make coverage-server` - Run server tests with coverage
- `make test-e2e-server` - Run E2E tests for server features
- `make test-e2e-ui` - Run E2E tests for UI features

### Linting

- `make lint` - Lint both server and UI code
- `make lint-server` - Run golangci-lint on Go code
- `make lint-ui` - Run ESLint on TypeScript/React code

### Other

- `make clean` - Remove build artifacts
- `make help` - Show all available make targets

## Project Architecture

### Backend (Go)

Fider uses a layered architecture with clean separation of concerns:

**Core Structure:**

- `main.go` - Entry point with command routing (ping, migrate, server)
- `app/cmd/` - Command implementations and server bootstrap
- `app/cmd/routes.go` - All routes are defined here
- `app/handlers/` - HTTP handlers organized by functionality
- `app/middlewares/` - HTTP middleware chain
- `app/models/` - Data models (cmd, dto, entity, enum, query)
- `app/services/` - Service implementations with dependency injection
- `app/pkg/` - Reusable packages and utilities

**Key Patterns:**

- **Bus Architecture**: Uses `app/pkg/bus` for service registration and dispatch
- **CQRS**: Commands and queries are separated in `app/models/`
- **Service Layer**: All external services (email, blob storage, oauth) are abstracted
- **Middleware Chain**: Authentication, tenant resolution, CORS, etc.

**Database:**

- PostgreSQL with custom migration system in `migrations/`
- SQL-based data access through service interfaces
- Tenant-aware data isolation

### Frontend (React/TypeScript)

Modern React application with TypeScript:

**Structure:**

- `public/index.tsx` - Application entry point with React 18
- `public/components/` - Reusable UI components
- `public/pages/` - Page-level components organized by feature
- `public/services/` - Client-side services and utilities
- `public/hooks/` - Custom React hooks

**Key Features:**

- **SSR Support**: Server-side rendering with hydration
- **Internationalization**: LinguiJS for i18n with locale switching
- **Component Library**: Extensive set of reusable components
- **State Management**: React Context for global state
- **Error Boundaries**: Comprehensive error handling

**Build System:**

- Webpack for bundling with CSS extraction
- ESBuild for SSR compilation
- SCSS for styling with utility classes
- Asset optimization and code splitting

### API Design

RESTful API with multiple access levels:

- `/api/v1/` - Public API (no auth required)
- Member API - Authenticated users
- Staff API - Collaborators and administrators
- Admin API - Administrators only

### Key Services

The application includes pluggable services for:

- **Email**: SMTP, Mailgun, AWS SES
- **Blob Storage**: Filesystem, S3, SQL
- **OAuth**: Custom providers, GitHub, Google, etc.
- **Billing**: Stripe integration (optional)
- **Webhooks**: Outbound event notifications

## Development Setup Requirements

1. **Go 1.22+** - Backend development
2. **Node.js 21/22** - Frontend build tools and TypeScript
3. **Docker** - PostgreSQL and local SMTP (MailHog)
4. **Air** - Go hot reload: `go install github.com/cosmtrek/air`
5. **Godotenv** - Environment loading: `go install github.com/joho/godotenv/cmd/godotenv`
6. **golangci-lint** - Go linting: `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2`

**Environment Setup:**

- Copy `.example.env` to `.env` for local configuration
- Run `docker compose up -d` for PostgreSQL and MailHog
- MailHog UI available at http://localhost:8025

## Testing Strategy

**Backend Testing:**

- Unit tests alongside source files (`*_test.go`)
- Integration tests with test database
- E2E tests using Cucumber for server features

**Frontend Testing:**

- Jest for unit/component tests
- Testing Library for React components
- E2E tests using Cucumber + Playwright

**Test Environment:**

- Uses `.test.env` for test-specific configuration
- Automated database migrations before test runs
- Coverage reporting available

## Code Organization Principles

### Backend Model Naming Conventions

Follow these strict naming patterns for `app/models/`:

- **`action.<something>`** - User interactions for POST/PUT/PATCH requests, map 1-to-1 with Commands (e.g., `action.CreateNewUser`)
- **`dto.<something>`** - Data transfer objects between packages/services (e.g., `dto.NewUserInfo`)
- **`entity.<something>`** - Objects mapped to database tables (e.g., `entity.User`)
- **`cmd.<something>`** - Commands that must be executed and potentially return values (e.g., `cmd.HttpRequest`, `cmd.LogDebug`, `cmd.SendMail`, `cmd.CreateNewUser`)
- **`query.<something>`** - Queries to get information from somewhere (e.g., `query.GetUserById`, `query.GetAllPosts`)

### Frontend Structure Conventions

- **Page Organization**: Each page has its own folder under `public/pages/` with:
  - `index.ts` - Module exporter
  - `[PageName].page.tsx` - Main page component
  - `[PageName].page.scss` - Page-specific styles
  - `[PageName].page.spec.tsx` - Unit tests
  - `./components/` - Page-specific components

### CSS Naming Conventions

Fider uses BEM methodology combined with utility classes:

- **`p-<page_name>`** - HTML ID for each page component (e.g., `p-home`, `p-user-settings`)
- **`c-<component_name>`** - Block class for components (e.g., `c-toggle`)
- **`c-<component_name>__<element>`** - Element classes (e.g., `c-toggle__label`)
- **`c-<component_name>--<state>`** - State modifiers (e.g., `c-toggle--checked`)
- **`is-<state>`, `has-<state>`** - Global state modifiers
- **Utility classes** - No prefix, used for common styling patterns. All utility classes are defined in public/assets/styles/utility

### General Principles

**Backend:**

- Services are dependency-injected through the bus system
- All external dependencies are abstracted behind interfaces
- Database queries are centralized in query objects
- Handlers focus on HTTP concerns, business logic in services

**Frontend:**

- Page components are lazy-loaded for performance
- Shared components in `components/common/`
- Business logic in services, not components
- Type-safe API calls with proper error handling

## Build and Deployment

**Local Development:**

- `make watch` for development with hot reload
- Webpack dev server for fast UI rebuilds
- Air for Go server hot reload

**Production Build:**

- `make build` creates optimized binaries and assets
- SSR compilation for better SEO and performance
- Asset optimization and minification

This is a mature, production-ready feedback platform with comprehensive testing, i18n support, and a clean, maintainable architecture.
