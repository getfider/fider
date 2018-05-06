package env_test

import (
	"os"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
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
	RegisterT(t)

	key := env.GetEnvOrDefault("UNKNOWN_KEY", "some value")
	Expect(key).Equals("some value")

	path := env.GetEnvOrDefault("PATH", "default path")
	Expect(path).NotEquals("default path")
}

func TestCurrent(t *testing.T) {
	RegisterT(t)

	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		Expect(env.Current()).Equals(testCase.env)
	}
}

func TestIsEnvironment(t *testing.T) {
	RegisterT(t)

	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		actual := testCase.isEnv()
		Expect(actual).IsTrue()
	}
}

func TestMustGet(t *testing.T) {
	RegisterT(t)

	Expect(func() {
		env.MustGet("THIS_DOES_NOT_EXIST")
	}).Panics()
}

func TestMultiTenantDomain(t *testing.T) {
	RegisterT(t)

	os.Setenv("AUTH_ENDPOINT", "https://login.test.fider.io:3000")
	Expect(env.MultiTenantDomain()).Equals(".test.fider.io")
	os.Setenv("AUTH_ENDPOINT", "https://login.test.fider.io")
	Expect(env.MultiTenantDomain()).Equals(".test.fider.io")
	os.Setenv("HOST_MODE", "single")
	Expect(env.MultiTenantDomain()).IsEmpty()
}

func TestGetPublicIP(t *testing.T) {
	RegisterT(t)

	os.Setenv("HOST_MODE", "multi")
	os.Setenv("AUTH_ENDPOINT", "https://login.dev.fider.io")
	ip, err := env.GetPublicIP()
	Expect(err).IsNil()
	Expect(ip).Equals("127.0.0.1")

	os.Setenv("AUTH_ENDPOINT", "https://login.fider.io")
	ip, err = env.GetPublicIP()
	Expect(err).IsNil()
	Expect(ip).Equals("174.138.111.44")

	os.Setenv("AUTH_ENDPOINT", "invalidinvalidinvalid.io")
	ip, err = env.GetPublicIP()
	Expect(err).IsNotNil()
	Expect(ip).IsEmpty()

	os.Setenv("HOST_MODE", "single")
	ip, err = env.GetPublicIP()
	Expect(err).IsNil()
	Expect(ip).Equals("")
}
