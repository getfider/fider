package validate

import (
	"fmt"
	"regexp"
)

var r, _ = regexp.Compile("^[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]$")

//Subdomain validates given subdomain
func Subdomain(subdomain string) (bool, []string) {
	if len(subdomain) <= 2 {
		return false, []string{"Subdomain must be more than 2 characters"}
	}

	if len(subdomain) > 40 {
		return false, []string{"Subdomain must be less than 40 characters"}
	}

	if !r.MatchString(subdomain) {
		return false, []string{"Subdomain contains invalid characters"}
	}

	switch subdomain {
	case
		"signup", "fider", "admin", "setup", "about", "wecanhearyou":
		return false, []string{fmt.Sprintf("%s is not a valid subdomain name", subdomain)}
	}

	return true, make([]string, 0)
}
