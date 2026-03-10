# Quick Test Guide

Fast setup to test sub-path hosting locally with HTTPS.

## 1. Setup hosts file

**Windows (run as Administrator):**
```powershell
.\setup-hosts.ps1
```

**Linux/Mac:**
```bash
chmod +x setup-hosts.sh
sudo ./setup-hosts.sh
```

Or manually add to your hosts file:
```
127.0.0.1 fider.local
127.0.0.1 app.local
127.0.0.1 multi.local
127.0.0.1 tenant1.multi.local
127.0.0.1 tenant2.multi.local
```

## 2. Pull the test image

```bash
docker pull ghcr.io/3rg0n/fider:subpath-fix
```

If the image is private, authenticate first:
```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

## 3. Start the test environment

```bash
docker-compose -f docker-compose-test.yml up -d
```

Wait ~10 seconds for all services to start, then check:
```bash
docker-compose -f docker-compose-test.yml ps
```

All services should show "Up" status.

## 4. Run the tests

### Test Scenario 1: Single-host (no sub-path)
Open: https://fider.local

Accept the certificate warning, complete setup wizard.

**Quick checks:**
- Homepage at `/`
- Create a post, verify redirect to `/posts/{id}/...`
- Delete post, verify redirect to `/`

### Test Scenario 2: Single-host WITH sub-path (THE FIX)
Open: https://app.local

You should see a landing page. Click "Go to Feedback".

Open: https://app.local/feedback

**Critical checks:**
- Homepage at `/feedback` (not `/`)
- Create post → redirects to `/feedback/posts/{id}/...`
- Delete post → redirects to `/feedback/` (not `/`)
- All links include `/feedback` prefix
- DevTools Network tab: API calls go to `/feedback/api/v1/...`

### Test Scenario 3: Multi-host (no sub-path)
Open: https://multi.local

Create tenant with subdomain `tenant1`.

Open: https://tenant1.multi.local

**Quick checks:**
- Homepage at `/` (tenant subdomain)
- Create/delete posts work with root-relative paths
- No `/feedback` prefix anywhere

## 5. View logs

```bash
docker-compose -f docker-compose-test.yml logs -f
```

## 6. Cleanup

```bash
docker-compose -f docker-compose-test.yml down -v
```

---

## Troubleshooting

**Certificate warnings:**
Normal for self-signed certs. Click "Advanced" → "Proceed".

**Port conflicts:**
If ports 80/443 are in use:
```bash
# Stop conflicting services
# On Windows: Stop IIS or other web servers
# On Linux: sudo systemctl stop nginx/apache2
```

**Services not starting:**
```bash
# Check logs
docker-compose -f docker-compose-test.yml logs postgres
docker-compose -f docker-compose-test.yml logs caddy
```

**Database connection errors:**
Wait for postgres health check:
```bash
docker-compose -f docker-compose-test.yml ps postgres
# Should show "(healthy)"
```

**Image pull fails:**
Make the package public at:
https://github.com/users/3rg0n/packages/container/fider/settings

Or authenticate with GITHUB_TOKEN that has `read:packages` scope.

---

For detailed test checklists, see: **TEST-SCENARIOS.md**
