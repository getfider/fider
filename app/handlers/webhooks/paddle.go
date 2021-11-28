package webhooks

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"sort"
	"time"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/log"
	"github.com/getfider/fider/app/pkg/web"
)

// IncomingPaddleWebhook handles all incoming requests from Paddle Webhooks
func IncomingPaddleWebhook() web.HandlerFunc {
	return func(c *web.Context) error {
		params, err := url.ParseQuery(c.Request.Body)
		if err != nil {
			return c.Failure(err)
		}

		err = verifyPaddleSig(params, env.Config.Paddle.PublicKey)
		if err != nil {
			return c.Failure(errors.Wrap(err, "failed to verity paddle signature"))
		}

		action := params.Get("alert_name")
		switch action {
		case "subscription_created":
			return handlePaddleSubscriptionCreated(c, params)
		case "subscription_cancelled":
			return handlePaddleSubscriptionCancelled(c, params)
		default:
			log.Warnf(c, "Unsupported Paddle webhook action: '@{Action}'", dto.Props{
				"Action": action,
			})
			return c.Ok(web.Map{})
		}
	}
}

func handlePaddleSubscriptionCreated(c *web.Context, params url.Values) error {
	passthrough := dto.PaddlePassthrough{}
	if err := json.Unmarshal([]byte(params.Get("passthrough")), &passthrough); err != nil {
		return c.Failure(err)
	}

	activate := &cmd.ActivateBillingSubscription{
		TenantID:       passthrough.TenantID,
		SubscriptionID: params.Get("subscription_id"),
		PlanID:         params.Get("subscription_plan_id"),
	}

	if err := bus.Dispatch(c, activate); err != nil {
		return c.Failure(err)
	}

	return c.Ok(web.Map{})
}

func handlePaddleSubscriptionCancelled(c *web.Context, params url.Values) error {
	passthrough := dto.PaddlePassthrough{}
	if err := json.Unmarshal([]byte(params.Get("passthrough")), &passthrough); err != nil {
		return c.Failure(err)
	}

	cancellationEffectiveDate := params.Get("cancellation_effective_date")
	subscriptionEndsAt, err := time.Parse("2006-01-02", cancellationEffectiveDate)
	if err != nil {
		log.Error(c, errors.Wrap(err, "failed to parse date '%s'", cancellationEffectiveDate))
		subscriptionEndsAt = time.Now().AddDate(0, 0, 30)
	}

	cancel := &cmd.CancelBillingSubscription{
		TenantID:           passthrough.TenantID,
		SubscriptionEndsAt: subscriptionEndsAt,
	}

	if err := bus.Dispatch(c, cancel); err != nil {
		return c.Failure(err)
	}

	return c.Ok(web.Map{})
}

// verifyPaddleSig verifies the p_signature parameter sent
// in Paddle webhooks. 'values' is the decoded form values sent
// in the webhook response body.
func verifyPaddleSig(values url.Values, publicKeyPEM string) error {
	der, _ := pem.Decode([]byte(publicKeyPEM))
	if der == nil {
		return errors.New("Could not parse public key pem")
	}

	pub, err := x509.ParsePKIXPublicKey(der.Bytes)
	if err != nil {
		return errors.New("Could not parse public key pem der")
	}

	signingKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("Not the correct key format")
	}

	sig, err := base64.StdEncoding.DecodeString(values.Get("p_signature"))
	if err != nil {
		return err
	}

	values.Del("p_signature")
	sha1Sum := sha1.Sum(phpserialize(values))
	err = rsa.VerifyPKCS1v15(signingKey, crypto.SHA1, sha1Sum[:], sig)
	if err != nil {
		return err
	}

	return nil
}

func phpserialize(form url.Values) []byte {
	var keys []string
	for k := range form {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	serialized := fmt.Sprintf("a:%d:{", len(keys))
	for _, k := range keys {
		serialized += fmt.Sprintf("s:%d:\"%s\";s:%d:\"%s\";", len(k), k, len(form.Get(k)), form.Get(k))
	}
	serialized += "}"

	return []byte(serialized)
}
