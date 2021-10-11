package mailgun

import (
	"context"
	"fmt"
	"strings"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/log"
)

// Known base URLs
// Should Mailgun add other regions we'll just need to add their URLs here
// Use upper case keys - incoming env var values are normalized before being used
var baseURLs = map[string]string{
	"US": "https://api.mailgun.net/v3/%s",
	"EU": "https://api.eu.mailgun.net/v3/%s",
}

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "Mailgun"
}

func (s Service) Category() string {
	return "email"
}

func (s Service) Enabled() bool {
	return env.Config.Email.Mailgun.APIKey != ""
}

func (s Service) Init() {
	bus.AddListener(sendMail)
	bus.AddHandler(fetchRecentSupressions)
}

// Try getting the URL of the Mailgun API using Environment vars and the Sender's domain
// Fall back to the URL for the US region if that fails to maintain compatibility with older installs
func getEndpoint(ctx context.Context, domain, path string) string {
	var regionCode = env.Config.Email.Mailgun.Region
	regionCode = strings.ToUpper(regionCode)

	// Default to the US domain if no region code was provided (ENV not set)
	// or if the provided code isn't valid
	if len(regionCode) < 1 {
		regionCode = "US"
	} else if len(baseURLs[regionCode]) < 1 {
		log.Warnf(ctx,
			"Unknown Mailgun region code '@{Code}' configured - falling back to 'US'",
			dto.Props{
				"Code": env.Config.Email.Mailgun.Region,
			},
		)

		regionCode = "US"
	}

	return fmt.Sprintf(baseURLs[regionCode], domain) + path
}
