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

	current := env.Current()
	defer func() {
		os.Setenv("GO_ENV", current)
	}()

	for _, testCase := range envs {
		os.Setenv("GO_ENV", testCase.go_env)
		Expect(env.Current()).Equals(testCase.env)
	}
}

func TestPath(t *testing.T) {
	RegisterT(t)

	Expect(env.Path("etc/deep/file.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/etc/deep/file.txt")
	Expect(env.Path("etc/file.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/etc/file.txt")
	Expect(env.Path("///etc/file.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/etc/file.txt")
	Expect(env.Path("/etc/file.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/etc/file.txt")
	Expect(env.Path("./etc/file.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/etc/file.txt")
	Expect(env.Path("file.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/file.txt")
	Expect(env.Path("/file.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/file.txt")
	Expect(env.Path("")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider")

	Expect(env.Etc("a.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/etc/a.txt")
	Expect(env.Etc("b.txt")).Equals(os.Getenv("GOPATH") + "/src/github.com/getfider/fider/etc/b.txt")
}

func TestIsEnvironment(t *testing.T) {
	RegisterT(t)

	current := env.Current()
	defer func() {
		os.Setenv("GO_ENV", current)
	}()

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

func TestHasLegal(t *testing.T) {
	RegisterT(t)

	Expect(env.HasLegal()).IsTrue()
}

func TestMultiTenantDomain(t *testing.T) {
	RegisterT(t)

	os.Setenv("HOST_DOMAIN", "test.fider.io")
	Expect(env.MultiTenantDomain()).Equals(".test.fider.io")
	os.Setenv("HOST_DOMAIN", "dev.fider.io")
	Expect(env.MultiTenantDomain()).Equals(".dev.fider.io")
	os.Setenv("HOST_DOMAIN", "fider.io")
	Expect(env.MultiTenantDomain()).Equals(".fider.io")
}

func TestGetPublicIP(t *testing.T) {
	RegisterT(t)

	os.Setenv("HOST_MODE", "multi")
	os.Setenv("HOST_DOMAIN", "dev.fider.io")
	ip, err := env.GetPublicIP()
	Expect(err).IsNil()
	Expect(ip).Equals("127.0.0.1")

	os.Setenv("HOST_DOMAIN", "fider.io")
	ip, err = env.GetPublicIP()
	Expect(err).IsNil()
	Expect(ip).Equals("174.138.111.44")

	os.Setenv("HOST_DOMAIN", "invalidinvalidinvalid.io")
	ip, err = env.GetPublicIP()
	Expect(err).IsNotNil()
	Expect(ip).IsEmpty()

	os.Setenv("HOST_MODE", "single")
	ip, err = env.GetPublicIP()
	Expect(err).IsNil()
	Expect(ip).Equals("")
}

func TestSubdomain(t *testing.T) {
	RegisterT(t)

	Expect(env.Subdomain("demo.test.assets-fider.io")).Equals("")

	os.Setenv("CDN_HOST", "test.assets-fider.io:3000")

	Expect(env.Subdomain("demo.test.fider.io")).Equals("demo")
	Expect(env.Subdomain("demo.test.assets-fider.io")).Equals("demo")
	Expect(env.Subdomain("test.fider.io")).Equals("")
	Expect(env.Subdomain("test.assets-fider.io")).Equals("")
	Expect(env.Subdomain("helloworld.com")).Equals("")

	os.Setenv("HOST_MODE", "single")

	Expect(env.Subdomain("demo.test.fider.io")).Equals("")
	Expect(env.Subdomain("demo.test.assets-fider.io")).Equals("")
	Expect(env.Subdomain("test.fider.io")).Equals("")
	Expect(env.Subdomain("test.assets-fider.io")).Equals("")
	Expect(env.Subdomain("helloworld.com")).Equals("")
}
