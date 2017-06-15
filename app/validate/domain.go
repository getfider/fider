package validate

import (
	"fmt"
)

//Subdomain validates given subdomain
func Subdomain(subdomain string) (bool, []string) {
	if len(subdomain) <= 2 {
		return false, []string{"Subdomain must be more than 2 characters"}
	}

	if len(subdomain) > 40 {
		return false, []string{"Subdomain must be less than 40 characters"}
	}

	switch subdomain {
	case
		"signup", "fider", "admin", "setup", "about", "wecanhearyou":
		return false, []string{fmt.Sprintf("%s is not a valid subdomain name", subdomain)}
	}

	return true, make([]string, 0)
}
