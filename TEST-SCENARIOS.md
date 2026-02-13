# Local HTTPS Testing Scenarios

This document describes how to test all sub-path hosting scenarios locally with HTTPS using self-signed certificates.

## Prerequisites

1. Docker and Docker Compose installed
2. Hosts file configured (add these entries to `/etc/hosts` on Linux/Mac or `C:\Windows\System32\drivers\etc\hosts` on Windows):

```
127.0.0.1 fider.local
127.0.0.1 app.local
127.0.0.1 multi.local
127.0.0.1 tenant1.multi.local
127.0.0.1 tenant2.multi.local
```

3. Pull the test image:
```bash
docker pull ghcr.io/3rg0n/fider:subpath-fix
```

## Starting the Test Environment

```bash
# Start all services
docker-compose -f docker-compose-test.yml up -d

# Watch logs
docker-compose -f docker-compose-test.yml logs -f

# Stop all services
docker-compose -f docker-compose-test.yml down

# Clean up (removes databases and volumes)
docker-compose -f docker-compose-test.yml down -v
```

## Test Scenarios

### Scenario 1: Single-host WITHOUT sub-path

**URL:** https://fider.local

**Configuration:**
- `HOST_MODE=single`
- `BASE_URL=https://fider.local`

**Initial Setup:**
1. Visit https://fider.local (accept self-signed certificate warning)
2. Complete the initial setup wizard
3. Create an admin account

**Tests to perform:**

- [ ] Homepage loads at `/` (not `/feedback`)
- [ ] Create a new post
  - [ ] After creation, redirected to `/posts/{id}/{slug}` (no `/feedback` prefix)
- [ ] Delete a post
  - [ ] After deletion, redirected to `/` (not `/feedback/`)
- [ ] Navigation links
  - [ ] Header logo links to `/`
  - [ ] Admin link goes to `/admin`
  - [ ] Settings link goes to `/settings`
- [ ] Sign out
  - [ ] Redirected to `/`
- [ ] API calls
  - [ ] Open browser DevTools Network tab
  - [ ] Verify API calls go to `/api/v1/...` (no prefix)
- [ ] OAuth (if configured)
  - [ ] OAuth redirect URLs use `https://fider.local`

**Expected:** All paths are root-relative with no sub-path prefix.

---

### Scenario 2: Single-host WITH sub-path (PRIMARY FIX)

**URL:** https://app.local/feedback

**Configuration:**
- `HOST_MODE=single`
- `BASE_URL=https://app.local/feedback`

**Initial Setup:**
1. Visit https://app.local (landing page)
2. Click "Go to Feedback (Fider)"
3. Visit https://app.local/feedback
4. Complete the initial setup wizard
5. Create an admin account

**Tests to perform:**

- [ ] Homepage loads at `/feedback` (not `/`)
- [ ] Landing page at `/` is NOT Fider (shows "App Landing Page")
- [ ] Create a new post
  - [ ] After creation, redirected to `/feedback/posts/{id}/{slug}` (WITH prefix)
- [ ] Delete a post
  - [ ] After deletion, redirected to `/feedback/` (not `/`)
- [ ] Navigation links
  - [ ] Header logo links to `/feedback/`
  - [ ] Admin link goes to `/feedback/admin`
  - [ ] Settings link goes to `/feedback/settings`
  - [ ] Sign out link goes to `/feedback/signout`
- [ ] Sign out
  - [ ] Redirected to `/feedback/`
- [ ] API calls
  - [ ] Open browser DevTools Network tab
  - [ ] Verify API calls go to `/feedback/api/v1/...` (WITH prefix)
  - [ ] Create/vote/delete actions should use `/feedback/api/v1/...`
- [ ] OAuth (if configured)
  - [ ] OAuth redirect URLs use `https://app.local/feedback`
- [ ] Browser URL bar
  - [ ] All navigation keeps `/feedback` prefix
  - [ ] Refreshing any page works (no 404s)
- [ ] Admin panel
  - [ ] All admin links include `/feedback/admin/*`
  - [ ] Export links include `/feedback/admin/export/*`
  - [ ] Billing links include `/feedback/admin/billing`

**Expected:** All paths include `/feedback` prefix. No redirects to `/` without the prefix.

---

### Scenario 3: Multi-host WITHOUT sub-path

**URLs:**
- https://multi.local (initial setup)
- https://tenant1.multi.local (tenant 1)
- https://tenant2.multi.local (tenant 2)

**Configuration:**
- `HOST_MODE=multi`
- `HOST_DOMAIN=multi.local`

**Initial Setup:**
1. Visit https://multi.local
2. Complete initial setup for tenant 1 (subdomain: `tenant1`)
3. Visit https://tenant1.multi.local
4. Create admin account for tenant 1
5. Visit https://multi.local again
6. Create tenant 2 (subdomain: `tenant2`)
7. Visit https://tenant2.multi.local
8. Create admin account for tenant 2

**Tests to perform:**

**For each tenant (tenant1.multi.local and tenant2.multi.local):**

- [ ] Homepage loads at tenant's subdomain root `/`
- [ ] Create a new post
  - [ ] After creation, redirected to `/posts/{id}/{slug}` (no sub-path)
- [ ] Delete a post
  - [ ] After deletion, redirected to `/` (tenant subdomain root)
- [ ] Navigation links
  - [ ] All links are root-relative (no `/feedback` or other prefix)
  - [ ] Admin link goes to `/admin`
- [ ] API calls
  - [ ] Verify API calls go to `/api/v1/...` at the tenant subdomain
  - [ ] No sub-path prefix in API URLs
- [ ] Tenant isolation
  - [ ] Posts created in tenant1 don't appear in tenant2
  - [ ] Users in tenant1 can't access tenant2 admin
- [ ] CNAME support (if available)
  - [ ] Set CNAME in admin settings
  - [ ] Verify access via custom domain works

**Expected:** Each tenant operates independently at their subdomain with root-relative paths. No sub-path prefixes.

---

## Debugging

### View logs
```bash
# All services
docker-compose -f docker-compose-test.yml logs -f

# Specific service
docker-compose -f docker-compose-test.yml logs -f fider-subpath
docker-compose -f docker-compose-test.yml logs -f caddy
```

### Check service health
```bash
docker-compose -f docker-compose-test.yml ps
```

### Access database
```bash
docker-compose -f docker-compose-test.yml exec postgres psql -U fider -d fider
docker-compose -f docker-compose-test.yml exec postgres psql -U fider -d fider_subpath
docker-compose -f docker-compose-test.yml exec postgres psql -U fider -d fider_multi
```

### View MailHog
Access captured emails at: http://localhost:8025

### Restart specific service
```bash
docker-compose -f docker-compose-test.yml restart fider-subpath
```

---

## Certificate Warnings

Caddy generates self-signed certificates for local testing. Your browser will show a certificate warning. To proceed:

- **Chrome/Edge:** Click "Advanced" → "Proceed to [domain] (unsafe)"
- **Firefox:** Click "Advanced" → "Accept the Risk and Continue"

Alternatively, install Caddy's root CA certificate in your system trust store:

```bash
# Export Caddy's root CA (while services are running)
docker-compose -f docker-compose-test.yml exec caddy cat /data/caddy/pki/authorities/local/root.crt > caddy-root.crt

# Then install caddy-root.crt to your system's trusted certificates
```

---

## Test Results Template

Copy this template for documenting test results:

```markdown
## Test Results - [Date]

### Scenario 1: Single-host without sub-path
- [ ] All tests passed
- Issues found: [None / List issues]

### Scenario 2: Single-host with sub-path
- [ ] All tests passed
- Issues found: [None / List issues]

### Scenario 3: Multi-host without sub-path
- [ ] All tests passed
- Issues found: [None / List issues]

### Browser tested: [Chrome/Firefox/Safari/Edge version]
### OS: [Windows/macOS/Linux]
```

---

## Cleanup

```bash
# Stop and remove all containers, networks, and volumes
docker-compose -f docker-compose-test.yml down -v

# Remove test image (optional)
docker rmi ghcr.io/3rg0n/fider:subpath-fix
```
