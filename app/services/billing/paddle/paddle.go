package paddle

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
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
}

func getApiBasePath() string {
	if env.Config.Paddle.IsSandbox {
		return "https://sandbox-vendors.paddle.com"
	}
	return "https://vendors.paddle.com"
}

func generateCheckoutLink(ctx context.Context, c *cmd.GenerateCheckoutLink) error {
	params := url.Values{}
	params.Set("vendor_id", env.Config.Paddle.VendorID)
	params.Set("vendor_auth_code", env.Config.Paddle.VendorAuthCode)
	params.Set("product_id", env.Config.Paddle.PlanID)

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

	c.URL = jsonq.New(string(res.Response)).String("url")
	return nil
}
