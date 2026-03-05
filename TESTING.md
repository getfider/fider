# Testing Matrix for Sub-Path Hosting Support

This document covers all supported installation types and their test status for PRs #1454 and #1455.

## Installation Types

Fider supports two host modes:
1. **Single-host mode** (`HOST_MODE=single`) - One Fider instance for one tenant
2. **Multi-host mode** (`HOST_MODE=multi`) - One Fider instance serving multiple tenants via subdomains

Each mode can be deployed with or without a sub-path.

---

## Test Matrix

| Mode | Configuration | Example BASE_URL | Test Coverage | Status |
|------|---------------|-----------------|---------------|--------|
| **Single-host** | No sub-path | `https://example.com` | Unit test: `TestBaseURL_SingleHostMode` | ✅ PASS |
| **Single-host** | With sub-path | `https://example.com/feedback` | Unit test: `TestBaseURL_SingleHostMode_WithPath`<br>Manual: Deployed to production | ✅ PASS<br>✅ VERIFIED |
| **Multi-host** | No sub-path | `HOST_DOMAIN=.fider.io`<br>(tenants at `tenant1.fider.io`) | Unit tests: `TestBaseURL`, `TestBaseURL_HTTPS`, `TestBaseURL_HTTPS_Proxy`, `TestTenantURL` | ✅ PASS |
| **Multi-host** | With sub-path | Not a standard configuration | N/A - unsupported scenario | ⚠️ NOT APPLICABLE |

---

## Unit Test Results

All `BaseURL()` tests pass:

```
=== RUN   TestBaseURL
--- PASS: TestBaseURL (0.00s)
=== RUN   TestBaseURL_HTTPS
--- PASS: TestBaseURL_HTTPS (0.00s)
=== RUN   TestBaseURL_HTTPS_Proxy
--- PASS: TestBaseURL_HTTPS_Proxy (0.00s)
=== RUN   TestBaseURL_SingleHostMode
--- PASS: TestBaseURL_SingleHostMode (0.00s)
=== RUN   TestBaseURL_SingleHostMode_WithPath
--- PASS: TestBaseURL_SingleHostMode_WithPath (0.00s)
=== RUN   TestTenantURL
--- PASS: TestTenantURL (0.00s)
=== RUN   TestTenantURL_WithCNAME
--- PASS: TestTenantURL_WithCNAME (0.00s)
=== RUN   TestTenantURL_SingleHostMode
--- PASS: TestTenantURL_SingleHostMode (0.00s)
```

---

## Manual Testing

### Single-host with sub-path (Primary use case)

**Configuration:**
```bash
HOST_MODE=single
BASE_URL=https://okratlas.cisco.com/feedback
```

**Reverse proxy setup:**
- Next.js app at `okratlas.cisco.com/`
- Fider at `okratlas.cisco.com/feedback/`
- Proxy routes `/feedback/*` → Fider backend
- Proxy routes `/api/v1/*` → Fider API

**Test scenarios verified:**
- ✅ Homepage loads at `/feedback/`
- ✅ Navigation links include `/feedback` prefix
- ✅ API calls go to `/feedback/api/v1/...` (not `/api/v1/...`)
- ✅ Create post redirects to `/feedback/posts/{id}/{slug}`
- ✅ Delete post redirects to `/feedback/` (not `/`)
- ✅ Sign-in/sign-out redirects include `/feedback`
- ✅ Admin panel navigation includes `/feedback/admin/*`
- ✅ OAuth redirects use full BASE_URL with path
- ✅ Atom feed links include `/feedback` prefix

**Result:** All navigation and API calls respect the configured sub-path.

---

## Architecture Notes

### Single-host mode
In single-host mode, `Context.BaseURL()` returns `env.Config.BaseURL` directly:
```go
func (c *Context) BaseURL() string {
    if env.IsSingleHostMode() {
        return env.Config.BaseURL  // Returns full URL including path
    }
    return c.Request.BaseURL()
}
```

This value is passed to the frontend as `Fider.settings.baseURL`, which the `basePath()` utility extracts the path component from:
```typescript
export function basePath(): string {
  try {
    const p = new URL(Fider.settings.baseURL).pathname.replace(/\/$/, "")
    return p === "/" ? "" : p  // Returns "/feedback" or ""
  } catch {
    return ""
  }
}
```

### Multi-host mode
In multi-host mode, `Context.BaseURL()` calls `Request.BaseURL()` which derives the URL from the incoming HTTP request (tenant subdomain):
```go
func (r *Request) BaseURL() string {
    address := r.URL.Scheme + "://" + r.URL.Hostname()
    if r.URL.Port() != "" {
        address += ":" + r.URL.Port()
    }
    return address  // Returns only scheme://host:port (no path)
}
```

Each tenant gets their own subdomain (e.g., `tenant1.fider.io`, `tenant2.fider.io`), so sub-path routing doesn't apply to multi-tenant mode.

---

## Related Issues

- #1452 - Core `Context.BaseURL()` bug (fixes path component being stripped)
- #1453 - Full sub-path hosting support (umbrella issue for all hardcoded paths)

---

## Summary

✅ **Single-host without sub-path** - Fully tested and working (existing functionality preserved)
✅ **Single-host with sub-path** - Fully tested and verified in production
✅ **Multi-host without sub-path** - Fully tested and working (existing functionality preserved)
⚠️ **Multi-host with sub-path** - Not a supported configuration (architectural mismatch)
