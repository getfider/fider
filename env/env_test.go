package env_test

import (
	"os"

	"github.com/WeCanHearYou/wchy/env"
	. "github.com/onsi/ginkgo"
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

var _ = Describe("GetEnvOrDefault", func() {
	It("should return null when key is unknown", func() {
		value := env.GetEnvOrDefault("UNKNOWN_KEY", "some value")
		Expect(value).To(Equal("some value"))
	})

	It("should return env variable when key is known", func() {
		value := env.GetEnvOrDefault("PATH", "default path")
		Expect(value).NotTo(Equal("default path"))
	})
})

var _ = Describe("GetCurrentDomain", func() {
	It("should return correct domain based on env", func() {
		for _, testCase := range envs {
			os.Setenv("GO_ENV", testCase.go_env)
			actual := env.GetCurrentDomain()
			Expect(actual).To(Equal(testCase.domain))
		}
	})
})

var _ = Describe("Current", func() {
	It("should return correct environment", func() {
		for _, testCase := range envs {
			os.Setenv("GO_ENV", testCase.go_env)
			actual := env.Current()
			Expect(actual).To(Equal(testCase.env))
		}
	})
})
