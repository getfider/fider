package env

import (
	"fmt"
	"os"

	"path"
)

// GetEnvOrDefault retrieves the value of the environment variable named by the key.
// It returns the value if available, otherwise returns defaultValue
func GetEnvOrDefault(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// MustGet returns environment variable or panic
func MustGet(name string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Errorf("Could not find environment variable named '%s'", name))
	}
	return value
}

// Mode returns HOST_MODE or its default value
func Mode() string {
	return GetEnvOrDefault("HOST_MODE", "single")
}

// IsSingleHostMode returns true if host mode is set to single tenant
func IsSingleHostMode() bool {
	return Mode() == "single"
}

// Current returns current Fider environment
func Current() string {
	env := os.Getenv("GO_ENV")
	switch env {
	case "test":
		return "test"
	case "production":
		return "production"
	}
	return "development"
}

// IsProduction returns true on Fider production environment
func IsProduction() bool {
	return Current() == "production"
}

// IsTest returns true on Fider test environment
func IsTest() bool {
	return Current() == "test"
}

// Path returns root path of project + given path
func Path(p ...string) string {
	root := "./"
	if IsTest() {
		root = os.Getenv("GOPATH") + "/src/github.com/getfider/fider/"
	}

	elems := append([]string{root}, p...)
	return path.Join(elems...)
}

// IsDevelopment returns true on Fider production environment
func IsDevelopment() bool {
	return Current() == "development"
}
