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
