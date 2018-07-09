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

var hasLegal *bool

// HasLegal returns true if current instance contains legal documents: privacy.md and terms.md
func HasLegal() bool {
	if hasLegal == nil {
		_, err1 := os.Stat(Etc("privacy.md"))
		_, err2 := os.Stat(Etc("terms.md"))
		exists := err1 == nil && err2 == nil
		hasLegal = &exists
	}
	return *hasLegal
}

// MultiTenantDomain returns domain name of current instance for multi tenant hosts
func MultiTenantDomain() string {
	if !IsSingleHostMode() {
		return "."+MustGet("HOST_DOMAIN")
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

// Subdomain returns the Fider subdomain (if available) from given host
func Subdomain(host string) string {
	if IsSingleHostMode() {
		return ""
	}

	domain := MultiTenantDomain()
	if domain != "" && strings.Contains(host, domain) {
		return strings.Replace(host, domain, "", -1)
	}

	if IsDefined("CDN_HOST") {
		domain = MustGet("CDN_HOST")
		parts := strings.Split(domain, ":")
		if parts[0] != "" && strings.Contains(host, "."+parts[0]) {
			return strings.Replace(host, "."+parts[0], "", -1)
		}
	}

	return ""
}
