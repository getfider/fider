package validate

import (
	"strings"
)

//Idea validates given idea
func Idea(title string, description string) (bool, []string, error) {

	if strings.Trim(title, " ") == "" {
		return false, []string{"Title is required."}, nil
	}

	if len(title) < 10 || len(strings.Split(title, " ")) < 3 {
		return false, []string{"Title needs to be more descriptive."}, nil
	}

	return true, make([]string, 0), nil
}
