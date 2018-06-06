package validate

import (
	"fmt"
	"regexp"

	"strings"

	"github.com/getfider/fider/app/storage"
)

var domainRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$")

//Subdomain validates given subdomain
func Subdomain(tenants storage.Tenant, subdomain string) *Result {
	subdomain = strings.ToLower(subdomain)

	if len(subdomain) <= 2 {
		return Failed([]string{"Subdomain must have more than 2 characters."})
	}

	if len(subdomain) > 40 {
		return Failed([]string{"Subdomain must have less than 40 characters."})
	}

	if !domainRegex.MatchString(subdomain) {
		return Failed([]string{"Subdomain contains invalid characters."})
	}

	switch subdomain {
	case
		"signup", "fider", "login", "customers", "admin", "setup", "about",
		"wecanhearyou", "dev", "mail", "billing", "www", "web", "translate",
		"help", "support", "status", "staging", "cdn", "assets":
		return Failed([]string{fmt.Sprintf("%s is a reserved subdomain.", subdomain)})
	}

	available, err := tenants.IsSubdomainAvailable(subdomain)
	if err != nil {
		return Error(err)
	}

	if !available {
		return Failed([]string{"This subdomain is not available anymore."})
	}

	return Success()
}
