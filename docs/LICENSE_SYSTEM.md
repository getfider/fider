# Commercial License System

This document explains how to enable commercial features on your self-hosted Fider instance using a license key.

## Overview

Fider's commercial features (like content moderation) require a valid license key. When you purchase a Pro subscription from the hosted Fider platform, you'll receive a license key that you can use to unlock these features on your self-hosted instance.

## Enabling Commercial Features

### Step 1: Purchase a Pro Subscription

Visit the hosted Fider platform and subscribe to the Pro plan. After your payment is processed, you'll receive a commercial license key.

### Step 2: Configure Your Self-Hosted Instance

Add the following to your `.env` file:

```env
# Public key for license validation (provided by Fider)
LICENSE_PUBLIC_KEY=lPHzCZhOBBihIusKWs5lzXCgGxZEKBpCiplkmZSjGpU=

# Your commercial license key (received after purchase)
COMMERCIAL_KEY=FIDER-COMMERCIAL-xxx-xxx-xxx
```

### Step 3: Restart Fider

Restart your Fider instance. Commercial features (like content moderation) will now be enabled.

## How It Works

Your license key is validated locally when Fider starts:

1. Reads `COMMERCIAL_KEY` from your environment
2. Verifies the cryptographic signature using `LICENSE_PUBLIC_KEY`
3. If valid, enables commercial features

**No internet connection required** - validation happens entirely on your server.

## Security & Privacy

Fider uses **Ed25519 cryptographic signatures** to validate license keys:

- ✅ **Secure**: Modern elliptic curve cryptography
- ✅ **Private**: No phone-home or tracking
- ✅ **Offline**: Works without internet connection
- ✅ **Trustworthy**: Public key cryptography prevents forgery

## Troubleshooting

### Commercial features are not enabled

**Check the following:**

1. Both `COMMERCIAL_KEY` and `LICENSE_PUBLIC_KEY` are set in your `.env` file
2. The `LICENSE_PUBLIC_KEY` matches the one provided by Fider
3. Your license key is in the correct format: `FIDER-COMMERCIAL-{numbers}-{numbers}-{hex}`
4. You restarted Fider after adding the configuration
5. Check your logs for validation errors

### "LICENSE_PUBLIC_KEY environment variable must be set"

**Problem:** Fider cannot validate your license without the public key.

**Solution:** Add `LICENSE_PUBLIC_KEY` to your `.env` file (provided in your license email or documentation).

### "invalid signature: license key verification failed"

**Problem:** Your license key cannot be verified.

**Possible causes:**
- Wrong `LICENSE_PUBLIC_KEY` (make sure it matches what Fider provided)
- License key was incorrectly copied (check for typos or missing characters)
- License key is not valid for this installation

**Solution:** Double-check both the public key and license key match what you received from Fider. If issues persist, contact support.

## Support

If you have issues with your commercial license, please contact Fider support with:
- Your license key
- Fider version
- Any error messages from logs
