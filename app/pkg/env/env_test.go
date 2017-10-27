package env_test

import (
	"os"
	"testing"

	"github.com/getfider/fider/app/pkg/env"
	. "github.com/onsi/gomega"
)

var envs = []struct {
	go_env string
	env    string
	isEnv  func() bool
}{
	{"test", "test", env.IsTest},
	{"development", "development", env.IsDevelopment},
	{"production", "production", env.IsProduction},
	{"anything", "development", env.IsDevelopment},
}

func TestGetEnvOrDefault(t *testing.T) {
	RegisterTestingT(t)

	key := env.GetEnvOrDefault("UNKNOWN_KEY", "some value")
	Expect(key).To(Equal("some value"))

	path := env.GetEnvOrDefault("PATH", "default path")
	Expect(path).NotTo(Equal("default path"))
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

func TestMustGet(t *testing.T) {
	RegisterTestingT(t)
	Expect(func() {
		env.MustGet("THIS_DOES_NOT_EXIST")
	}).To(Panic())
}

func TestMultiTenantDomain(t *testing.T) {
	RegisterTestingT(t)

	os.Setenv("AUTH_ENDPOINT", "https://login.test.fider.io:3000")
	Expect(env.MultiTenantDomain()).To(Equal(".test.fider.io"))

	os.Setenv("AUTH_ENDPOINT", "https://login.test.fider.io")
	Expect(env.MultiTenantDomain()).To(Equal(".test.fider.io"))

	os.Setenv("HOST_MODE", "single")
	Expect(env.MultiTenantDomain()).To(Equal(""))
}
