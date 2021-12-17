package paddle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jsonq"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Paddle"
}

func (s Service) Category() string {
	return "billing"
}

func (s Service) Enabled() bool {
	return env.IsBillingEnabled()
}

func (s Service) Init() {
	bus.AddHandler(generateCheckoutLink)
	bus.AddHandler(getBillingSubscription)
}

func getApiBasePath() string {
	if env.Config.Paddle.IsSandbox {
		return "https://sandbox-vendors.paddle.com"
	}
	return "https://vendors.paddle.com"
}

// generateCheckoutLink generates a checkout link using Paddle API
// It's a 2-steps process:
// 1. Create a "Custom Checkout Create URL
// 2. Use that URL to generate a Checkout Link and return it
func generateCheckoutLink(ctx context.Context, c *cmd.GenerateCheckoutLink) error {
	passthrough, err := json.Marshal(c.Passthrough)
	if err != nil {
		return errors.Wrap(err, "failed to marshal Passthrough object")
	}

	params := url.Values{}
	params.Set("vendor_id", env.Config.Paddle.VendorID)
	params.Set("vendor_auth_code", env.Config.Paddle.VendorAuthCode)
	params.Set("product_id", env.Config.Paddle.PlanID)
	params.Set("passthrough", string(passthrough))

	req := &cmd.HTTPRequest{
		URL:    fmt.Sprintf("%s/api/2.0/product/generate_pay_link", getApiBasePath()),
		Body:   strings.NewReader(params.Encode()),
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}

	if err := bus.Dispatch(ctx, req); err != nil {
		return errors.Wrap(err, "failed to generate paddle checkout link")
	}

	if req.ResponseStatusCode >= 300 {
		return errors.New("unexpected status code while generating a paddle checkout link: %d", req.ResponseStatusCode)
	}

	res := &PaddleResponse{}
	if err := json.Unmarshal(req.ResponseBody, &res); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	if !res.IsSuccess {
		return errors.New("failed to generate paddle checkout link with '%s (%d)'", res.Error.Message, res.Error.Code)
	}

	// the "url" here is the "Custom Checkout Create URL", we need to append some extra QueryString parameters
	u, err := url.Parse(jsonq.New(string(res.Response)).String("url"))
	if err != nil {
		return errors.Wrap(err, "generated paddle checkout url is invalid")
	}

	// parent_url and parentURL must match one of the approved domains in Paddle.com
	q := u.Query()
	q.Set("parent_url", fmt.Sprintf("https://%s", env.Config.HostDomain))
	q.Set("parentURL", fmt.Sprintf("https://%s", env.Config.HostDomain))
	u.RawQuery = q.Encode()

	req = &cmd.HTTPRequest{
		URL:    u.String(),
		Method: http.MethodGet,
	}

	if err := bus.Dispatch(ctx, req); err != nil {
		return errors.Wrap(err, "failed to generate paddle checkout link")
	}

	if req.ResponseStatusCode <= 299 || req.ResponseStatusCode >= 400 {
		return errors.New("unexpected status code while generating a paddle checkout link: %d", req.ResponseStatusCode)
	}

	c.URL = req.ResponseHeader.Get("location")
	if c.URL == "" {
		return errors.New("response is missing 'location' header")
	}
	return nil
}

func getBillingSubscription(ctx context.Context, q *query.GetBillingSubscription) error {
	params := url.Values{}
	params.Set("vendor_id", env.Config.Paddle.VendorID)
	params.Set("vendor_auth_code", env.Config.Paddle.VendorAuthCode)
	params.Set("subscription_id", q.SubscriptionID)

	req := &cmd.HTTPRequest{
		URL:    fmt.Sprintf("%s/api/2.0/subscription/users", getApiBasePath()),
		Body:   strings.NewReader(params.Encode()),
		Method: http.MethodPost,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}

	if err := bus.Dispatch(ctx, req); err != nil {
		return errors.Wrap(err, "failed to get paddle subscription details")
	}

	if req.ResponseStatusCode >= 300 {
		return errors.New("unexpected status code while fetching paddle subscription details: %d", req.ResponseStatusCode)
	}

	res := &PaddleResponse{}
	if err := json.Unmarshal(req.ResponseBody, &res); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	if !res.IsSuccess {
		return errors.New("failed to fetch paddle subscription details with '%s (%d)'", res.Error.Message, res.Error.Code)
	}

	sub := []PaddleSubscriptionItem{}
	if err := json.Unmarshal(res.Response, &sub); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	if len(sub) > 0 {
		q.Result = &entity.BillingSubscription{
			CancelURL: sub[0].CancelURL,
			UpdateURL: sub[0].UpdateURL,
			PaymentInformation: entity.BillingPaymentInformation{
				PaymentMethod:  sub[0].PaymentInformation.PaymentMethod,
				CardType:       sub[0].PaymentInformation.CardType,
				LastFourDigits: sub[0].PaymentInformation.LastFourDigits,
				ExpiryDate:     sub[0].PaymentInformation.ExpiryDate,
			},
			LastPayment: entity.BillingLastPayment{
				Amount:   sub[0].LastPayment.Amount,
				Currency: sub[0].LastPayment.Currency,
				Date:     sub[0].LastPayment.Date,
			},
		}
	}
	return nil
}
