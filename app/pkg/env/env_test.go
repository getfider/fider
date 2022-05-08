package env_test

import (
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
	{"anything", "production", env.IsProduction},
}

func TestIsEnvironment(t *testing.T) {
	RegisterT(t)

	for _, testCase := range envs {
		env.Config.Environment = testCase.go_env
		actual := testCase.isEnv()
		Expect(actual).IsTrue()
	}
}

func TestHasLegal(t *testing.T) {
	RegisterT(t)

	Expect(env.HasLegal()).IsTrue()
}

func TestMultiTenantDomain(t *testing.T) {
	RegisterT(t)

	env.Config.HostDomain = "test.fider.io"
	Expect(env.MultiTenantDomain()).Equals(".test.fider.io")
	env.Config.HostDomain = "dev.fider.io"
	Expect(env.MultiTenantDomain()).Equals(".dev.fider.io")
	env.Config.HostDomain = "fider.io"
	Expect(env.MultiTenantDomain()).Equals(".fider.io")
}

func TestSubdomain(t *testing.T) {
	RegisterT(t)

	Expect(env.Subdomain("demo.test.fidercdn.com")).Equals("")

	env.Config.CDN.Host = "test.fidercdn.com:3000"

	Expect(env.Subdomain("demo.test.fider.io")).Equals("demo")
	Expect(env.Subdomain("demo.test.fidercdn.com")).Equals("demo")
	Expect(env.Subdomain("test.fider.io")).Equals("")
	Expect(env.Subdomain("test.fidercdn.com")).Equals("")
	Expect(env.Subdomain("helloworld.com")).Equals("")

	env.Config.HostMode = "single"

	Expect(env.Subdomain("demo.test.fider.io")).Equals("")
	Expect(env.Subdomain("demo.test.fidercdn.com")).Equals("")
	Expect(env.Subdomain("test.fider.io")).Equals("")
	Expect(env.Subdomain("test.fidercdn.com")).Equals("")
	Expect(env.Subdomain("helloworld.com")).Equals("")
}
