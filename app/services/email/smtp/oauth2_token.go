package smtp

import (
	"context"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/getfider/fider/app/pkg/errors"
)

func splitCommaScopes(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func getClientCredentialsToken(ctx context.Context, tokenURL, clientID, clientSecret string, scopes []string) (*oauth2.Token, error) {
	if tokenURL == "" {
		return nil, errors.New("smtp: oauth token url is required")
	}
	if clientID == "" || clientSecret == "" {
		return nil, errors.New("smtp: oauth client id/secret are required")
	}

	conf := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
		Scopes:       scopes,
	}

	tok, err := conf.Token(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "smtp: failed to fetch oauth token")
	}
	return tok, nil
}
