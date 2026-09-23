# Portal Directory — Design

Date: 2026-08-21
Status: approved, ready for implementation planning

## Problem

In multi-tenant mode (`HOST_MODE=multi`) a Fider instance resolves tenants purely from the
request hostname (`app/middlewares/tenant.go:48`). Hitting the root domain with no subdomain
leaves `c.Tenant()` nil, so `middlewares.RequireTenant` returns 404 for every route except
`/signup`. There is no query that lists tenants — only `GetFirstTenant` and
`GetTenantByDomain` in `app/models/query/tenant.go` — and no page that lists them.

The consequence: the only way to reach a portal is to already know its subdomain. Portals on
the same instance are undiscoverable.

## Solution

Serve a public portal directory at the root domain of a multi-tenant instance: a browsable,
client-filterable list of every active, public portal, each entry linking to that portal's own
URL.

The directory is off by default behind a new environment flag, so upgrading an existing
instance never starts enumerating its tenants without the operator opting in.

## Scope decisions

Settled during brainstorming; each rejected option is additive and can be revisited later.

| Decision | Choice | Rejected |
| --- | --- | --- |
| Audience | Public — anyone hitting the root domain | Signed-in "my portals"; operator-only view |
| Eligibility | Automatic: active + non-private | Per-portal opt-in setting (needs a migration + admin UI) |
| Placement | Root `/` of the root domain | Dedicated `/portals` path |
| Default | Off, behind an env flag | Always on |
| Entry content | Logo, name, host | Activity stats (post counts, last activity) |
| Scale | Whole list in page props + client-side filter | Server-side search and pagination |

Explicitly out of scope: opt-in setting, activity stats, pagination, a JSON API endpoint, a
signup call to action on the directory.

## Architecture

Five units, each independently testable:

1. **Config flag** (`app/pkg/env`) — is the directory enabled at all?
2. **Query + store** (`app/models/query`, `app/services/sqlstore/postgres`) — which portals are
   eligible?
3. **Route wiring** (`app/middlewares`, `app/cmd/routes.go`) — when does the root domain serve
   the directory instead of 404ing?
4. **Handler** (`app/handlers`) — assemble the view data.
5. **Page** (`public/pages/PortalDirectory`) — render and filter.

### 1. Config flag

Add to the `config` struct in `app/pkg/env/env.go`, alongside the other feature flags:

```go
PortalDirectoryEnabled bool `env:"PORTAL_DIRECTORY_ENABLED,default=false"`
```

And a helper beside `IsMultiHostMode` (`env.go:255`):

```go
// IsPortalDirectoryEnabled returns true when the public portal directory should be served
// on the root domain. Single-host instances have exactly one portal, so it never applies.
func IsPortalDirectoryEnabled() bool {
	return IsMultiHostMode() && Config.PortalDirectoryEnabled
}
```

Document the variable in `.example.env`, following the commented-out-with-explanation style
already used there for `ALLOW_PRIVATE_NETWORK_TARGETS`:

```
# PORTAL_DIRECTORY_ENABLED=true
# Serves a public directory of this instance's portals at the root domain. Multi-tenant
# (HOST_MODE=multi) only. Lists every active, non-private portal by name and links to it.
```

Folding the host-mode check into the helper means no caller has to remember it, and a
single-tenant instance behaves identically whether the flag is set or not.

### 2. Query and store

`app/models/query/tenant.go`:

```go
// GetPublicTenants returns the tenants eligible for the public portal directory:
// active, non-private, and not inside a deletion grace window. Ordered by name.
type GetPublicTenants struct {
	// Output
	Result []*entity.Tenant
}
```

`app/services/sqlstore/postgres/tenant.go` — `getPublicTenants`, reusing the column list and
`tenants_billing` join from `getFirstTenant` (`tenant.go:258`) so `dbEntities.Tenant.ToModel`
is fully populated:

```sql
-- Column list is copied verbatim from getFirstTenant, ending with
-- t.scheduled_deletion_at and the computed has_paddle_subscription.
SELECT t.id, t.name, t.subdomain, ... , t.scheduled_deletion_at,
	(b.paddle_subscription_id IS NOT NULL AND b.stripe_subscription_id IS NULL) AS has_paddle_subscription
FROM tenants t
LEFT JOIN tenants_billing b ON b.tenant_id = t.id
WHERE t.status = $1
  AND t.is_private = false
  AND t.scheduled_deletion_at IS NULL
ORDER BY t.name
```

The full column list is not restated here to avoid two copies drifting apart; take it from
`getFirstTenant` and change only the `WHERE` and `ORDER BY`. Note `trx.Select` (plural) rather
than `trx.Get`, and an empty result is not an error — unlike `getFirstTenant`, zero rows is a
valid answer.

`$1` is `enum.TenantActive`. Registered with `bus.AddHandler(getPublicTenants)` in
`postgres.go` next to the existing tenant handlers (`postgres.go:109`).

**Eligibility rationale.** `status = TenantActive` excludes pending signups, disabled tenants,
and locked (payment-lapsed) tenants. `scheduled_deletion_at IS NULL` excludes sites the owner
has asked to delete but whose grace window has not elapsed. `is_private = false` matters most:
`CheckTenantPrivacy` already blocks anonymous access to private portals, and listing their
names would leak exactly what that setting protects.

`prevent_indexing` is deliberately **not** part of the filter. `CreateTenant` inserts it as
`true` for every new tenant (`tenant.go:246`), so honouring it would leave the directory
permanently empty. It governs search-engine indexing, not membership of an on-instance list.

### 3. Route wiring

`app/cmd/routes.go:141` registers `r.Get("/", handlers.Index())` after
`r.Use(middlewares.RequireTenant())`. The engine is `julienschmidt/httprouter`
(`app/pkg/web/engine.go:80`), which panics on a duplicate method+path registration, so a
second `"/"` route is not an option.

Instead, a fallthrough middleware inserted immediately before `RequireTenant` — the exact
point where the root-domain 404 originates today:

```go
// RootDomainFallback serves the given handler for the root path of a multi-tenant
// instance's root domain, where no tenant resolves and RequireTenant would 404.
func RootDomainFallback(handler web.HandlerFunc) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			isRootDomain := c.Request.URL.Hostname() == env.Config.HostDomain

			if c.Tenant() == nil && isRootDomain && c.Request.URL.Path == "/" {
				return handler(c)
			}
			return next(c)
		}
	}
}
```

**The host check is load-bearing.** "No tenant resolved" is not the same as "this is the root
domain": it is equally true of an unknown subdomain, a *disabled* tenant's subdomain (the
`MultiTenant` middleware deliberately leaves the tenant unset for those,
`middlewares/tenant.go:63`), and a custom domain pointed here that matches no tenant. Guarding
only on `tenant == nil && path == "/"` turns every one of those 404s into a 200 serving the
directory — confirmed by running the app, where `disabled.localhost/` and `nosuch.localhost/`
both rendered `PortalDirectory/PortalDirectory.page`.

Comparing against `env.Config.HostDomain` is the exact test. `env.Subdomain(host) == ""` is
*not* a valid substitute: it also returns `""` for any host that does not end in the
multi-tenant domain, so unmatched custom domains would still slip through
(`env.go:306-324`).

In `routes.go`, immediately above `r.Use(middlewares.RequireTenant())`:

```go
// The root domain of a multi-tenant instance resolves no tenant, so RequireTenant would
// 404. Serve the portal directory there instead when the operator has enabled it.
r.Use(middlewares.RootDomainFallback(handlers.PortalDirectory()))
r.Use(middlewares.RequireTenant())
```

The middleware takes the handler as a `web.HandlerFunc` parameter, so `middlewares` does not
import `handlers`; `routes.go` wires the two, as it already does for everything else.

The flag check lives in the handler rather than the middleware, so the middleware stays a
pure "is this the root domain's root path" decision and the handler owns its own
enable/disable behaviour. With the flag off the handler returns 404 — the same response as
today.

Alternatives rejected:

- **Move `"/"` into an early group** — the group would have to re-apply `RequireTenant`,
  `BlockPendingTenants`, and `CheckTenantPrivacy` for the tenant case. Duplicated middleware
  wiring that will drift from the real chain.
- **Teach `RequireTenant` to render the directory** — smallest diff, but it forces
  `middlewares` to import `handlers` and overloads a middleware whose single job is a guard.

### 4. Handler and view data

`app/models/dto/portal.go`:

```go
// PortalSummary is one entry in the public portal directory.
type PortalSummary struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Host    string `json:"host"`
	LogoURL string `json:"logoURL,omitempty"`
}
```

`app/handlers/portals.go`:

```go
// PortalDirectory lists the public portals hosted on this instance. Only reachable on the
// root domain of a multi-tenant instance, and only when the operator has enabled it.
func PortalDirectory() web.HandlerFunc {
	return func(c *web.Context) error {
		if !env.IsPortalDirectoryEnabled() {
			return c.NotFound()
		}

		q := &query.GetPublicTenants{}
		if err := bus.Dispatch(c, q); err != nil {
			return c.Failure(err)
		}

		portals := make([]dto.PortalSummary, 0, len(q.Result))
		for _, tenant := range q.Result {
			url := web.TenantBaseURL(c, tenant)
			portals = append(portals, dto.PortalSummary{
				Name:    tenant.Name,
				URL:     url,
				Host:    strings.TrimPrefix(strings.TrimPrefix(url, "https://"), "http://"),
				LogoURL: web.TenantLogoURL(c, tenant),
			})
		}

		return c.Page(http.StatusOK, web.Props{
			Page:        "PortalDirectory/PortalDirectory.page",
			Title:       "Portals",
			Description: "Browse the feedback portals hosted here.",
			Data:        web.Map{"portals": portals},
		})
	}
}
```

`web.Props.Data` is a `web.Map` reaching the page as `fider.session.props` — the same
mechanism `LegalPage` uses (`app/handlers/common.go:44`). The renderer already tolerates a nil
tenant (`app/pkg/web/renderer.go:150,180`), which is why `/signup` works on the root domain.

**Absolute logo URLs.** `/static/images/*bkey` is registered inside the `tenantAssets` group,
after `RequireTenant` (`routes.go:96`), so it only resolves on a host that maps to a tenant.
The client-side `uploadedImageURL` helper (`public/services/utils.ts:104`) builds URLs from
`Fider.settings.assetsURL`, which on the root domain has no tenant and would 404. So the
handler emits absolute per-portal logo URLs instead, via a new helper beside the existing
`LogoURL` (`app/pkg/web/context.go:599`):

```go
// TenantLogoURL returns an absolute URL to the given tenant's logo, or "" when it has none.
// Unlike LogoURL it takes the tenant explicitly, so it works on the root domain where no
// tenant is in context.
func TenantLogoURL(ctx context.Context, tenant *entity.Tenant) string {
	if tenant.LogoBlobKey == "" {
		return ""
	}
	if env.Config.CDN.Host != "" {
		request := ctx.Value(app.RequestCtxKey).(Request)
		return request.URL.Scheme + "://" + tenant.Subdomain + "." + env.Config.CDN.Host +
			"/static/images/" + tenant.LogoBlobKey + "?size=200"
	}
	return TenantBaseURL(ctx, tenant) + "/static/images/" + tenant.LogoBlobKey + "?size=200"
}
```

The CDN branch mirrors `AssetsURL` (`context.go:591`). It returns `""` rather than the Fider
fallback logo that `LogoURL` uses, letting the page render its own initials placeholder.

### 5. Page

`public/pages/PortalDirectory/` — `PortalDirectory.page.tsx`, `PortalDirectory.scss`,
`index.ts`. Pages are resolved by name at runtime through `AsyncPage`
(`public/AsyncPages.tsx`), so no registration step is needed beyond the directory existing.

```tsx
interface PortalSummary {
  name: string
  url: string
  host: string
  logoURL?: string
}

interface PortalDirectoryPageProps {
  portals: PortalSummary[]
}
```

Structure:

- Renders bare content, not `<Page>`/`<Header>` — both read `fider.session.tenant`, which is
  nil here. `SignUp.page.tsx` is the precedent for a tenant-less page.
- A filter `<Input>` over `useState`, matching case-insensitively on `name` and `host`. Pure
  client-side; no network calls.
- Each portal is a card-shaped `<a href={portal.url}>`: logo `<img>` when `logoURL` is set,
  otherwise a CSS circle with the name's first letter (the letter-avatar endpoint at
  `routes.go:97` is also tenant-gated, so it cannot be used here); name; host as subtitle.
- Two distinct empty states: no portals at all, and no portals matching the filter.
- Utility classes from `public/assets/styles/utility/` first, per CLAUDE.md; page-specific
  BEM (`#p-portal-directory`, `.c-portal-card__*`) only for the card grid.
- User-facing copy wrapped in `<Trans>`, with ids added to `locale/en/client.json` (the
  committed source of truth; `locale/**/*.js` is generated and gitignored). The search
  placeholder uses `i18n._({id, message})`, the pattern already used in
  `CompleteSignInProfile.page.tsx`.

**SSR registration.** `public/ssr.tsx` maps page names to modules in a *static* table —
esbuild cannot do dynamic imports — and `ssrRender` throws `Page not found` for anything
missing from it. Server-side rendering only runs for crawlers
(`renderer.go:242`, gated on `Request.IsCrawler()`), and a throw there is logged and degraded
to the client-rendered shell rather than failing the response. A public directory page that
crawlers cannot read defeats the point, so `PortalDirectory/PortalDirectory.page` must be
added to that table. This is what makes the page's `export default` mandatory: `ssrRender`
reads `pages[...]?.default`.

## Data flow

```
GET / on root domain (multi-tenant)
  → middlewares.Tenant → MultiTenant: hostname matches no subdomain/cname, tenant stays nil
  → middlewares.RootDomainFallback: tenant == nil && path == "/" → handlers.PortalDirectory()
      → flag off?  → 404 (today's behaviour)
      → flag on    → bus.Dispatch(query.GetPublicTenants)
                   → postgres.getPublicTenants
                   → []dto.PortalSummary with absolute URL + logo URL per portal
                   → c.Page(props.Data["portals"])
  → PortalDirectory.page renders cards; filtering is local state
  → click → navigates to that portal's own host, where MultiTenant resolves it normally
```

Any other path on the root domain, and every path on a tenant host, falls through to
`RequireTenant` unchanged.

## Error handling

- Flag off → `c.NotFound()`. Identical to today, so upgrades are invisible until opted in.
- Query failure → `c.Failure(err)`, the codebase's standard 500 path.
- Zero eligible portals → 200 with an empty-state page, not an error. A fresh instance with
  no portals yet is a normal state.
- A portal that becomes private or disabled between page render and click → its own host
  applies `CheckTenantPrivacy` / the disabled-tenant check as usual. The directory is a
  pointer, never an authorisation decision.
- `RootDomainFallback` guards on `c.Tenant() == nil`, so a tenant host serving `/` is never
  intercepted even with the flag on.

## Testing

**Go — `app/handlers/portals_test.go`** (`mock.NewServer()` is the multi-tenant harness,
`app/pkg/mock/setup.go:34`):

- flag off → 404
- flag on, no tenant in context → 200, and the rendered props carry exactly the eligible
  portals
- each entry's `url` points at that portal's own host; `logoURL` is absolute and
  tenant-hosted; a logo-less portal yields an empty `logoURL`

**Go — `app/services/sqlstore/postgres/tenant_test.go`** (file exists), for
`getPublicTenants`:

- active public tenants are returned, ordered by name
- private, pending, locked, and disabled tenants are excluded
- a tenant with `scheduled_deletion_at` set is excluded

**Go — `app/middlewares/tenant_test.go`**, for `RootDomainFallback`:

- tenant nil and path `/` → the fallback handler runs
- tenant nil and path `/posts/1` → falls through to `next`
- tenant present and path `/` → falls through to `next`

**Jest — `public/pages/PortalDirectory/PortalDirectory.page.spec.tsx`:**

- renders one card per portal
- filtering narrows the list case-insensitively, by name and by host
- an empty `portals` prop shows the no-portals state, and hides the filter box
- a logo-less portal renders the initial placeholder rather than a broken `<img>`

The two empty states are asserted by element (`.c-portal-directory__empty`,
`.c-portal-directory__nomatch`) rather than by their copy. The lingui macro hoists `<Trans>`
children into a `message` prop, and `public/jest.setup.tsx` mocks `@lingui/react`'s `Trans` to
render `children` — so translated text is absent from the rendered output under test.

**Go — SSR path, in `portals_test.go`:** a request carrying a crawler User-Agent must come
back containing `c-portal-card` markup. That substring can only appear if the page rendered
server-side, so this fails if the page is ever dropped from the `public/ssr.tsx` table. It
needs a current `ssr.js`, which `make test-server` guarantees by depending on `build-ssr`.

Verification: `make lint` and `make test`.

## Files touched

New:

- `app/handlers/portals.go`, `app/handlers/portals_test.go`
- `app/models/dto/portal.go`
- `public/pages/PortalDirectory/{PortalDirectory.page.tsx,PortalDirectory.page.scss,index.ts,PortalDirectory.page.spec.tsx}`
- `docs/superpowers/specs/2026-08-21-portal-directory-design.md` (this file)

Modified:

- `public/ssr.tsx` — register the page in the static SSR table
- `locale/en/client.json` — the five new message ids
- `app/pkg/web/context_test.go` — `TenantLogoURL` tests
- `app/pkg/env/env_test.go` — flag test

- `app/pkg/env/env.go` — flag + `IsPortalDirectoryEnabled`
- `.example.env` — document `PORTAL_DIRECTORY_ENABLED`
- `app/models/query/tenant.go` — `GetPublicTenants`
- `app/services/sqlstore/postgres/tenant.go` — `getPublicTenants`
- `app/services/sqlstore/postgres/postgres.go` — register the handler
- `app/services/sqlstore/postgres/tenant_test.go` — store tests
- `app/middlewares/tenant.go` — `RootDomainFallback`
- `app/middlewares/tenant_test.go` — middleware tests
- `app/cmd/routes.go` — wire the fallback before `RequireTenant`
- `app/pkg/web/context.go` — `TenantLogoURL`

No database migration. No changes to existing tenant behaviour.
