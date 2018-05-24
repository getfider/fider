package env

import (
	"fmt"
	"net"
	"os"
	"strings"

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

// IsDefined returns true if given environment variable is defined
func IsDefined(name string) bool {
	value := os.Getenv(name)
	return value != ""
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

// MultiTenantDomain returns domain name of current instance for multi tenant hosts
func MultiTenantDomain() string {
	if !IsSingleHostMode() {
		endpoint := MustGet("AUTH_ENDPOINT")
		firstIdx := strings.Index(endpoint, ".")
		if firstIdx != -1 {
			lastIdx := strings.LastIndex(endpoint, ":")
			if lastIdx != -1 && lastIdx > firstIdx {
				return endpoint[firstIdx:lastIdx]
			}
			return endpoint[firstIdx:]
		}
		panic(fmt.Sprintf("Could not extract domain from %s", endpoint))
	}
	return ""
}

var publicIP = make(map[string]string, 0)

// GetPublicIP returns the public IP of current hosting server
func GetPublicIP() (string, error) {
	if IsSingleHostMode() {
		return "", nil
	}

	domain := MultiTenantDomain()[1:]
	if _, ok := publicIP[domain]; !ok {
		addr, err := net.LookupIP(domain)
		if err == nil {
			publicIP[domain] = addr[0].String()
		} else {
			return "", err
		}
	}
	return publicIP[domain], nil
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
	} else {
		panic(Current())
	}

	elems := append([]string{root}, p...)
	return path.Join(elems...)
}

// Etc returns a path to a folder or file inside the /etc/ folder
func Etc(p ...string) string {
	paths := append([]string{"etc"}, p...)
	return Path(paths...)
}

// IsDevelopment returns true on Fider production environment
func IsDevelopment() bool {
	return Current() == "development"
}
