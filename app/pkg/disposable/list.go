// Package disposable provides matching against a bundled list of known
// disposable / throwaway email providers, plus per-tenant deny/allow rules.
package disposable

import (
	_ "embed"
	"strings"
)

//go:embed domains.txt
var bundledRaw string

var bundled = parseList(bundledRaw)

func parseList(raw string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(strings.ToLower(line))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		out[line] = struct{}{}
	}
	return out
}

// IsBundled reports whether the given hostname matches a domain in the
// bundled disposable list, either exactly or as any subdomain.
func IsBundled(domain string) bool {
	domain = strings.TrimSpace(strings.ToLower(domain))
	domain = strings.TrimRight(domain, ".")
	if domain == "" || strings.HasPrefix(domain, ".") {
		return false
	}
	if strings.Contains(domain, "..") {
		return false
	}
	if _, ok := bundled[domain]; ok {
		return true
	}
	parts := strings.Split(domain, ".")
	for i := 1; i < len(parts); i++ {
		if _, ok := bundled[strings.Join(parts[i:], ".")]; ok {
			return true
		}
	}
	return false
}
