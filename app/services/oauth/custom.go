package oauth

import (
	"net/url"
	"strings"

	"github.com/getfider/fider/app/models/entity"
)

var providerParams = map[string]map[string]string{
	"id.twitch.tv": {
		"claims": `{"userinfo":{"preferred_username":null,"email":null,"email_verified":null}}`,
	},
}

func getProviderInitialParams(u *url.URL, config *entity.OAuthConfig) url.Values {
	v := url.Values{}
	hostname := u.Hostname()
	
	// Add static provider-specific parameters
	if params, ok := providerParams[hostname]; ok {
		for key, value := range params {
			v.Add(key, value)
		}
	}
	
	// Add dynamic parameters based on provider and scope
	// Apple Sign-In requires response_mode=form_post when name or email scopes are requested
	if hostname == "appleid.apple.com" && config != nil {
		scope := strings.ToLower(config.Scope)
		if strings.Contains(scope, "name") || strings.Contains(scope, "email") {
			v.Add("response_mode", "form_post")
		}
	}
	
	return v
}
