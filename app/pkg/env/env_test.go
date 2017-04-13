package env_test

import (
	"os"
	"testing"

	"github.com/WeCanHearYou/wechy/app/pkg/env"
	. "github.com/onsi/gomega"
)

var envs = []struct {
	go_env string
	domain string
	env    string
	isEnv  func() bool
}{
	{"test", "test.canhearyou.com", "test", env.IsTest},
	{"staging", "staging.canhearyou.com", "staging", env.IsStaging},
	{"development", "dev.canhearyou.com", "development", env.IsDevelopment},
	{"production", "canhearyou.com", "production", env.IsProduction},
	{"anything", "canhearyou.com", "development", env.IsDevelopment},
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

func TestIsEnvironment(t *testing.T) {
	RegisterTestingT(t)

	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		actual := testCase.isEnv()
		Expect(actual).To(BeTrue())
	}
}
