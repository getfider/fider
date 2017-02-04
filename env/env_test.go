package env_test

import (
	"os"
	"testing"

	"github.com/WeCanHearYou/wchy/env"
	. "github.com/onsi/gomega"
)

var envs = []struct {
	go_env string
	domain string
	env    string
}{
	{"test", "test.canhearyou.com", "test"},
	{"staging", "staging.canhearyou.com", "staging"},
	{"development", "dev.canhearyou.com", "development"},
	{"production", "canhearyou.com", "production"},
	{"anything", "canhearyou.com", "development"},
}

func TestGetEnvOrDefault(t *testing.T) {
	RegisterTestingT(t)

	key := env.GetEnvOrDefault("UNKNOWN_KEY", "some value")
	Expect(key).To(Equal("some value"))

	path := env.GetEnvOrDefault("PATH", "default path")
	Expect(path).NotTo(Equal("default path"))
}

func TestGetCurrentDomain(t *testing.T) {
	RegisterTestingT(t)

	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		actual := env.GetCurrentDomain()
		Expect(actual).To(Equal(testCase.domain))
	}
}

func TestCurrent(t *testing.T) {
	RegisterTestingT(t)

	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		actual := env.Current()
		Expect(actual).To(Equal(testCase.env))
	}
}
