# Internal: License System Setup

**This document is for internal use only** - setting up the Fider hosted instance to generate commercial licenses.

## Initial Setup (One-Time)

### 1. Generate License Key Pair

Run the key generation script:

```bash
go run scripts/generate_license_keys.go
```

This outputs:
```
=== Ed25519 License Key Pair Generated ===

PUBLIC KEY (for self-hosted instances):
lPHzCZhOBBihIusKWs5lzXCgGxZEKBpCiplkmZSjGpU=

PRIVATE KEY (for hosted server only - KEEP SECRET!):
B0Srq+DtsBI6MwY4sljCywXNLWAizfCdiJ6emGx+xIqU8fMJmE4EGKEi6wpazmXNcKAbFkQoGkKKmWSZlKMalQ==
```

### 2. Configure Production Environment

Add **only the private key** to the hosted instance's `.env`:

```env
LICENSE_PRIVATE_KEY=B0Srq+DtsBI6MwY4sljCywXNLWAizfCdiJ6emGx+xIqU8fMJmE4EGKEi6wpazmXNcKAbFkQoGkKKmWSZlKMalQ==
```

**Security:**
- Never commit this to git
- Store in secure secrets management (AWS Secrets Manager, Vault, etc.)
- Restrict access to production environment

### 3. Distribute Public Key

The **public key** is safe to share publicly. Include it in:
- Self-hosted installation documentation
- License purchase confirmation emails
- `.example.env` file (already done)
- Customer support materials

Current public key for distribution:
```
LICENSE_PUBLIC_KEY=lPHzCZhOBBihIusKWs5lzXCgGxZEKBpCiplkmZSjGpU=
```

## How License Generation Works

### Automatic Generation via Stripe

When a customer purchases Pro via Stripe webhook (`app/handlers/webhooks/stripe.go`):

1. **Checkout completed** (`handleCheckoutSessionCompleted`):
   - Calls `license.GenerateKey(tenantID)`
   - Stores key in database via `ActivateStripeSubscription`
   - Key is displayed to customer in billing page

2. **Subscription cancelled** (`handleSubscriptionDeleted`):
   - Calls `CancelStripeSubscription`
   - Sets tenant back to free plan

### Manual License Generation (if needed)

If you need to manually generate a license for a customer:

```go
// In Go code or via admin endpoint
licenseKey := license.GenerateKey(tenantID)
// Returns: FIDER-COMMERCIAL-{tenantID}-{timestamp}-{signature}
```

The tenant ID should match the customer's tenant in your database.

## Code Reference

### Key Files

- **Generation**: `app/services/license/validator.go` - `GenerateKey()` function
- **Webhook**: `app/handlers/webhooks/stripe.go` - Automatic generation on purchase
- **Storage**: `app/services/sqlstore/postgres/billing.go` - Database operations
- **Validation**: `app/services/license/validator.go` - `ValidateKey()` function (used by self-hosted)

### Database Tables

- `tenants.is_pro` - Boolean flag indicating Pro subscription
- `tenants_billing` - Stores Stripe subscription ID and license key

## Testing

### Test Environment

The `.test.env` file contains test keys for development:

```env
LICENSE_PUBLIC_KEY=9o3OnNs9tnnD780fZcTYas8uhcA90MWtin1rWGhXVPE=
LICENSE_PRIVATE_KEY=fwRNcFEiusrH0Hx1aBlOobE+6hm2/CPXUkUTWXj3cQn2jc6c2z22ecPvzR9lxNhqzy6FwD3Qxa2KfWtYaFdU8Q==
```

These are safe for testing and are different from production keys.

### Running Tests

```bash
make test-server
# or specifically:
godotenv -f .test.env go test -v ./app/services/license/...
```

All tests should pass:
- `TestGenerateAndValidateKey_Ed25519`
- `TestValidateKey_Ed25519_InvalidSignature`
- `TestValidateKey_Ed25519_WrongPublicKey`
- `TestValidateKey_EmptyKey`
- `TestValidateKey_InvalidFormat`
- `TestValidateKey_NoPublicKeyConfigured`

## Key Rotation (Emergency Procedure)

If the private key is ever compromised:

1. **Generate new key pair**:
   ```bash
   go run scripts/generate_license_keys.go
   ```

2. **Update production environment** with new `LICENSE_PRIVATE_KEY`

3. **Notify all self-hosted customers**:
   - Send new `LICENSE_PUBLIC_KEY`
   - They must update their `.env` file
   - Previous licenses will still work with new public key

4. **Update documentation** with new public key

5. **Consider**: Implement key versioning in license format for smoother transitions

## Monitoring

### What to Monitor

- License generation failures (check Stripe webhook logs)
- Invalid license validation attempts on hosted instance
- Customer support tickets about license issues

### Logs to Watch

```bash
# License generation
grep "GenerateKey" /var/log/fider/app.log

# Stripe webhook events
grep "stripe" /var/log/fider/app.log | grep "checkout\|subscription"
```

## Support Scenarios

### Customer Reports License Not Working

1. Verify they have correct `LICENSE_PUBLIC_KEY` in their `.env`
2. Check their license key format is correct
3. Test their license key validation locally:
   ```bash
   # Set up test environment with their key
   LICENSE_PUBLIC_KEY=<current_public_key> \
   go run scripts/validate_license.go "FIDER-COMMERCIAL-xxx-xxx-xxx"
   ```
4. Check if license was successfully generated in database
5. If needed, manually regenerate license for their tenant ID

### Customer Wants to Transfer License

Licenses are tied to tenant ID. To transfer:
- Generate new license for new tenant ID
- Cancel old subscription/license
- Update billing records

## Security Best Practices

1. **Never** share or commit the private key
2. **Always** use environment variables for production
3. **Rotate** keys if any compromise is suspected
4. **Audit** license generation logs periodically
5. **Monitor** for unusual patterns (same license on multiple IPs)
6. **Backup** production environment configuration securely

## Production Checklist

Before deploying to production:

- [ ] Production `LICENSE_PRIVATE_KEY` is set and secured
- [ ] Public key is distributed in documentation
- [ ] Stripe webhook is configured and tested
- [ ] All tests pass
- [ ] Monitoring/logging is in place
- [ ] Support team knows how to help customers with license issues
- [ ] Emergency key rotation procedure is documented
