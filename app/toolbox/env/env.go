package env

import "os"
import "fmt"

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

// GetCurrentDomain returns Wechy domain based on current environment variables
func GetCurrentDomain() string {
	env := os.Getenv("GO_ENV")
	switch env {
	case "test":
		return "test.canhearyou.com"
	case "development":
		return "dev.canhearyou.com"
	case "staging":
		return "staging.canhearyou.com"
	}
	return "canhearyou.com"

}

// Current returns current Wechy environment
func Current() string {
	env := os.Getenv("GO_ENV")
	switch env {
	case "test":
		return "test"
	case "staging":
		return "staging"
	case "production":
		return "production"
	}
	return "development"
}

// IsProduction returns true on Wechy production environment
func IsProduction() bool {
	return Current() == "production"
}

// IsStaging returns true on Wechy production environment
func IsStaging() bool {
	return Current() == "staging"
}

// IsTest returns true on Wechy production environment
func IsTest() bool {
	return Current() == "test"
}

// IsDevelopment returns true on Wechy production environment
func IsDevelopment() bool {
	return Current() == "development"
}
