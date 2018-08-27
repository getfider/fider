package validate

import (
	"fmt"
	"regexp"

	"strings"

	"github.com/getfider/fider/app/storage"
)

var domainRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$")

//Subdomain validates given subdomain
func Subdomain(tenants storage.Tenant, subdomain string) ([]string, error) {
	subdomain = strings.ToLower(subdomain)

	if len(subdomain) <= 2 {
		return []string{"Subdomain must have more than 2 characters."}, nil
	}

	if len(subdomain) > 40 {
		return []string{"Subdomain must have less than 40 characters."}, nil
	}

	if !domainRegex.MatchString(subdomain) {
		return []string{"Subdomain contains invalid characters."}, nil
	}

	switch subdomain {
	case
		"signup", "fider", "login", "customers", "admin", "setup", "about",
		"wecanhearyou", "dev", "mail", "billing", "www", "web", "translate",
		"help", "support", "status", "staging", "cdn", "assets", "live",
		"manage", "mgmt", "platform", "production", "development":
		return []string{fmt.Sprintf("%s is a reserved subdomain.", subdomain)}, nil
	}

	available, err := tenants.IsSubdomainAvailable(subdomain)
	if err != nil {
		return nil, err
	}

	if !available {
		return []string{"This subdomain is not available anymore."}, nil
	}

	return []string{}, nil
}
