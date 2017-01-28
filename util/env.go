package util

import "os"

// GetEnvOrDefault retrieves the value of the environment variable named by the key.
// It returns the value if available, otherwise returns defaultValue
func GetEnvOrDefault(name string, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetCurrentDomain returns WCHY domain based on current environment variables
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
