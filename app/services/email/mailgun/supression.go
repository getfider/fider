package mailgun

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
)

type EventsResponse struct {
	Items []struct {
		ID        string `json:"id"`
		Recipient string `json:"recipient"`
	} `json:"items"`
	Paging struct {
		Next string `json:"next"`
	} `json:"paging"`
}

func fetchRecentSupressions(ctx context.Context, q *query.FetchRecentSupressions) error {
	supressedAddresses := make(map[string]bool)

	params := url.Values{}
	params.Set("event", "failed")
	params.Set("limit", "300")
	params.Set("severity", "permanent")
	params.Set("begin", strconv.FormatInt(q.StartTime.Unix(), 10))

	endpointURL := fmt.Sprintf("%s?%s", getEndpoint(ctx, env.Config.Email.Mailgun.Domain, "/events"), params.Encode())
	for {
		res, err := fetchEvents(ctx, endpointURL)
		if err != nil {
			return err
		}

		for _, item := range res.Items {
			supressedAddresses[strings.TrimSpace(item.Recipient)] = true
		}

		// No more emails to fetch
		if len(res.Items) == 0 {
			q.EmailAddresses = make([]string, 0)
			for email := range supressedAddresses {
				q.EmailAddresses = append(q.EmailAddresses, email)
			}
			return nil
		}

		// continue on the next page
		endpointURL = res.Paging.Next
	}
}

func fetchEvents(ctx context.Context, endpoint string) (*EventsResponse, error) {
	req := &cmd.HTTPRequest{
		Method: "GET",
		URL:    endpoint,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		BasicAuth: &dto.BasicAuth{
			User:     "api",
			Password: env.Config.Email.Mailgun.APIKey,
		},
	}

	if err := bus.Dispatch(ctx, req); err != nil {
		return nil, errors.Wrap(err, "failed to fetch failed events")
	}

	if req.ResponseStatusCode >= 300 {
		return nil, errors.New("unexpected status code while fetching failed events: %d", req.ResponseStatusCode)
	}

	res := &EventsResponse{}
	if err := json.Unmarshal(req.ResponseBody, &res); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	return res, nil
}
