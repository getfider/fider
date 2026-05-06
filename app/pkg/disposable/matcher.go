package disposable

import "strings"

// IsBlocked reports whether the email's domain should be rejected.
// Allow rules win over deny rules and over the bundled list (admin's escape hatch).
// `deny` and `allow` are tenant-scoped domain lists; both nil is fine.
func IsBlocked(email string, deny, allow []string) bool {
	domain := domainFromEmail(email)
	if domain == "" {
		return false
	}
	if matchesAny(domain, allow) {
		return false
	}
	if matchesAny(domain, deny) {
		return true
	}
	return IsBundled(domain)
}

func domainFromEmail(email string) string {
	at := strings.LastIndex(email, "@")
	if at < 0 || at == len(email)-1 {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(email[at+1:]))
}

func matchesAny(domain string, rules []string) bool {
	for _, r := range rules {
		r = strings.TrimSpace(strings.ToLower(r))
		if r == "" {
			continue
		}
		if domain == r || strings.HasSuffix(domain, "."+r) {
			return true
		}
	}
	return false
}
