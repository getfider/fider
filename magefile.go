// +build mage

package main

import (
	"os"
	"runtime"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Aliases = map[string]interface{}{
	"watch": Watch.All,
}

func init() {
	os.Setenv("MAGEFILE_VERBOSE", "true")
}

var buildTime = time.Now().Format("2006.01.02.150405")
var buildNumber = os.Getenv("CIRCLE_BUILD_NUM")

func Run() error {
	return sh.Run("godotenv", "-f", ".env", "./fider")
}

func Lint() error {
	return sh.Run("npx", "tslint", "-c", "tslint.json", "'public/**/*.{ts,tsx}'", "'tests/**/*.{ts,tsx}'")
}

func Clean() error {
	os.RemoveAll("./dist")
	return os.Mkdir("./dist", 0777)
}

type Watch mg.Namespace

func (Watch) All() {
	Clean()
	mg.Deps(Watch.Server, Watch.UI)
}

func (Watch) UI() error {
	return sh.Run("npx", "webpack", "-w")
}

func (Watch) Server() error {
	return sh.Run("air", "-c", "air.conf")
}

type Build mg.Namespace

func (Build) All() {
	mg.Deps(Build.Server, Build.UI)
}

func (Build) Docker() error {
	mg.Deps(Build.UI)
	if err := buildServer(map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	}); err != nil {
		return err
	}

	return sh.Run("docker", "build", "-t", "getfider/fider", ".")
}

func (Build) Server() error {
	return buildServer(map[string]string{
		"GOOS":   runtime.GOOS,
		"GOARCH": runtime.GOARCH,
	})
}

func (Build) UI() error {
	mg.Deps(Clean)
	env := map[string]string{"NODE_ENV": "production"}
	return sh.RunWith(env, "npx", "webpack", "-p")
}

type Test mg.Namespace

func (Test) All() {
	mg.Deps(Test.Server, Test.UI)
}

func (Test) Coverage() error {
	mg.Deps(Build.Server)
	sh.Run("godotenv", "-f", ".test.env", "./fider", "migrate")
	return sh.Run("godotenv", "-f", ".test.env", "go", "test", "./...", "-coverprofile=cover.out", "-coverpkg=all", "-race")
}

func (Test) Server() error {
	mg.Deps(Build.Server)
	sh.Run("godotenv", "-f", ".test.env", "./fider", "migrate")
	return sh.Run("godotenv", "-f", ".test.env", "go", "test", "./...", "-race")
}

func (Test) UI() error {
	env := map[string]string{"TZ": "GMT"}
	return sh.RunWith(env, "npx", "jest", "./public")
}

// Utils
func buildServer(env map[string]string) error {
	ldflags := "-s -w -X main.buildtime=" + buildTime + " -X main.buildnumber=" + buildNumber
	return sh.RunWith(env, "go", "build", "-ldflags='"+ldflags+"'", "-o", "fider", ".")
}
