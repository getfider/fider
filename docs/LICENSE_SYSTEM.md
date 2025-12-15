# Commercial License System

This document explains how to enable commercial features on your self-hosted Fider instance using a license key.

## Overview

Fider's commercial features (like content moderation) require a valid license key. When you purchase a Pro subscription from the hosted Fider platform, you'll receive a license key that you can use to unlock these features on your self-hosted instance.

## Enabling Commercial Features

### Step 1: Purchase a Pro Subscription

Visit the hosted Fider platform and subscribe to the Pro plan. After your payment is processed, you'll receive a commercial license key.

### Step 2: Configure Your Self-Hosted Instance

Add your license key to your `.env` file:

```env
# Your commercial license key (received after purchase)
COMMERCIAL_KEY=FIDER-COMMERCIAL-xxx-xxx-xxx
```

**That's it!** The public key for validation is already embedded in Fider - no additional configuration needed.

### Step 3: Restart Fider

Restart your Fider instance. Commercial features (like content moderation) will now be enabled.

## How It Works

Your license key is validated locally when Fider starts:

1. Reads `COMMERCIAL_KEY` from your environment
2. Verifies the cryptographic signature using the embedded public key
3. If valid, enables commercial features

**No internet connection required** - validation happens entirely on your server.

### Advanced: Overriding the Public Key

The public key is embedded in Fider's code, so it updates automatically when you upgrade Fider. However, you can override it if needed:

```env
LICENSE_PUBLIC_KEY=your-custom-public-key-here
```

You typically **don't need** to set this unless:
- You're testing with a custom license system
- Fider support specifically instructs you to update it

## Security & Privacy

Fider uses **Ed25519 cryptographic signatures** to validate license keys:

- ✅ **Secure**: Modern elliptic curve cryptography
- ✅ **Private**: No phone-home or tracking
- ✅ **Offline**: Works without internet connection
- ✅ **Trustworthy**: Public key cryptography prevents forgery

## Troubleshooting

### Commercial features are not enabled

**Check the following:**

1. `COMMERCIAL_KEY` is set in your `.env` file
2. Your license key is in the correct format: `FIDER-COMMERCIAL-{numbers}-{numbers}-{hex}`
3. You restarted Fider after adding the configuration
4. You're running a recent version of Fider (with the embedded public key)
5. Check your logs for validation errors

### "invalid signature: license key verification failed"

**Problem:** Your license key cannot be verified.

**Possible causes:**
- License key was incorrectly copied (check for typos or missing characters)
- You're running an outdated version of Fider (update to get the latest public key)
- License key is not valid

**Solution:**
1. Verify your license key is copied correctly
2. Update to the latest Fider version
3. If issues persist, contact support

## Support

If you have issues with your commercial license, please contact Fider support with:
- Your license key
- Fider version
- Any error messages from logs
