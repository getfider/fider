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

func getClientCredentialsToken(ctx context.Context, tokenURL, clientID, clientSecret string, scopes []string) (string, error) {
	if tokenURL == "" {
		return "", errors.New("smtp: oauth token url is required")
	}
	if clientID == "" || clientSecret == "" {
		return "", errors.New("smtp: oauth client id/secret are required")
	}

	key := tokenSourceKey(tokenURL, clientID, scopes)

	tokenSourceMu.Lock()
	tokenSource, ok := tokenSourceByKey[key]
	if !ok {
		conf := clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			TokenURL:     tokenURL,
			Scopes:       scopes,
		}

		base := conf.TokenSource(ctx)
		tokenSource = oauth2.ReuseTokenSource(nil, base)
		tokenSourceByKey[key] = tokenSource
	}
	tokenSourceMu.Unlock()

	tok, err := tokenSource.Token()
	if err != nil {
		return "", errors.Wrap(err, "smtp: failed to fetch oauth token")
	}
	if tok == nil || tok.AccessToken == "" {
		return "", errors.New("smtp: oauth returned empty access token")
	}
	return tok.AccessToken, nil
}
