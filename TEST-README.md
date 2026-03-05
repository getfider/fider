# Sub-Path Hosting Test Suite

Complete test environment for verifying sub-path hosting support across all Fider deployment modes.

## Files in this test suite

- **TEST-SCENARIOS.md** - Detailed test scenarios and checklists for each deployment mode
- **QUICK-TEST.md** - Fast setup guide to get testing quickly
- **docker-compose-test.yml** - Multi-scenario test environment with Caddy reverse proxy
- **Caddyfile.test** - Caddy configuration for HTTPS with self-signed certificates
- **setup-hosts.sh** - Linux/Mac script to configure hosts file
- **setup-hosts.ps1** - Windows script to configure hosts file
- **TESTING.md** - Test matrix documentation for PR review

## What gets tested

### ✅ Scenario 1: Single-host without sub-path
- **URL:** https://fider.local
- **Config:** `HOST_MODE=single`, `BASE_URL=https://fider.local`
- **Tests:** Standard Fider at domain root (existing functionality)

### ✅ Scenario 2: Single-host WITH sub-path (THE PRIMARY FIX)
- **URL:** https://app.local/feedback
- **Config:** `HOST_MODE=single`, `BASE_URL=https://app.local/feedback`
- **Tests:** Fider under sub-path with root landing page (new functionality)

### ✅ Scenario 3: Multi-host without sub-path
- **URLs:** https://tenant1.multi.local, https://tenant2.multi.local
- **Config:** `HOST_MODE=multi`, `HOST_DOMAIN=multi.local`
- **Tests:** Multiple tenants via subdomains (existing functionality)

## Quick Start

```bash
# 1. Setup hosts file (Windows: run PowerShell as Admin)
.\setup-hosts.ps1          # Windows
sudo ./setup-hosts.sh      # Linux/Mac

# 2. Pull test image
docker pull ghcr.io/3rg0n/fider:subpath-fix

# 3. Start test environment
docker-compose -f docker-compose-test.yml up -d

# 4. Test each scenario
# - https://fider.local (single-host, no sub-path)
# - https://app.local/feedback (single-host WITH sub-path)
# - https://multi.local (multi-host setup)

# 5. Cleanup
docker-compose -f docker-compose-test.yml down -v
```

## Making the test image public

If you get authentication errors pulling the image, make it public:

1. Go to: https://github.com/users/3rg0n/packages/container/fider/settings
2. Scroll to "Danger Zone"
3. Click "Change visibility" → "Public"
4. Confirm the change

Alternatively, authenticate:
```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u 3rg0n --password-stdin
```

## Test Checklist

After testing all scenarios, document results in TESTING.md:

- [ ] Single-host without sub-path: All navigation and API calls work
- [ ] Single-host WITH sub-path: All paths include `/feedback` prefix
- [ ] Multi-host: Each tenant isolated with root-relative paths
- [ ] HTTPS works with self-signed certificates
- [ ] No JavaScript console errors
- [ ] No 404 errors in Network tab
- [ ] OAuth redirects use correct BASE_URL (if configured)

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Caddy (HTTPS)                        │
│              Automatic Self-Signed Certs                │
└──────────────────┬──────────────────────────────────────┘
                   │
         ┌─────────┼──────────┐
         │         │          │
    fider.local  app.local  *.multi.local
         │         │          │
         ▼         ▼          ▼
    ┌────────┐  ┌─────────┐  ┌──────────┐
    │Fider   │  │Landing +│  │Fider     │
    │Single  │  │Fider at │  │Multi     │
    │No path │  │/feedback│  │Subdomains│
    └────────┘  └─────────┘  └──────────┘
         │         │          │
         └─────────┴──────────┘
                   │
              PostgreSQL
```

## Related PRs

- **#1454** - Core Context.BaseURL() fix
- **#1455** - Full sub-path hosting support

## Maintainer Review

This test suite addresses the maintainer's request to:
> test against all supported types of installation [...] single host as is now, and on a sub-path, and the same for multi-host

Results documented in TESTING.md after running all scenarios.
