package validate

import (
	"fmt"
	"regexp"

	"github.com/getfider/fider/app/storage"
)

var r, _ = regexp.Compile("^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$")

//Subdomain validates given subdomain
func Subdomain(tenants storage.Tenant, subdomain string) (bool, []string, error) {
	if len(subdomain) <= 2 {
		return false, []string{"Subdomain must be more than 2 characters"}, nil
	}

	if len(subdomain) > 40 {
		return false, []string{"Subdomain must be less than 40 characters"}, nil
	}

	if !r.MatchString(subdomain) {
		return false, []string{"Subdomain contains invalid characters"}, nil
	}

	switch subdomain {
	case
		"signup", "fider", "admin", "setup", "about", "wecanhearyou":
		return false, []string{fmt.Sprintf("%s is not a valid subdomain name", subdomain)}, nil
	}

	available, err := tenants.IsSubdomainAvailable(subdomain)
	if err != nil {
		return false, make([]string, 0), err
	}

	if !available {
		return false, []string{"This subdomain is not available anymore"}, nil
	}

	return true, make([]string, 0), nil
}
