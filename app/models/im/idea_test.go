package im_test

import (
	"testing"

	"github.com/getfider/fider/app/models/im"
	. "github.com/onsi/gomega"
)

func TestInvalidIdeaTitles(t *testing.T) {
	RegisterTestingT(t)

	for _, title := range []string{
		"me",
		"",
		"  ",
		"signup",
		"my company",
		"my@company",
		"my.company",
		"my+company",
		"1234567890123456789012345678901234567890ABC",
	} {
		idea := im.Idea{Title: title}
		result := idea.Validate(services)
		ExpectFailed(result, "title")
	}
}

func TestValidIdeaTitles(t *testing.T) {
	RegisterTestingT(t)

	for _, title := range []string{
		"this is my new idea",
		"this idea is very descriptive",
	} {
		idea := im.Idea{Title: title}
		result := idea.Validate(services)
		ExpectSuccess(result)
	}
}
