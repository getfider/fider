package email

import (
	"regexp"
	"strings"

	"github.com/getfider/fider/app/pkg/env"
)

// NoReply is the default 'from' address
var NoReply = env.Config.Email.NoReply

var whitelist = env.Config.Email.Whitelist
var whitelistRegex = regexp.MustCompile(whitelist)
var blacklist = env.Config.Email.Blacklist
var blacklistRegex = regexp.MustCompile(blacklist)

// SetWhitelist can be used to change email whitelist during runtime
func SetWhitelist(s string) {
	whitelist = s
	whitelistRegex = regexp.MustCompile(whitelist)
}

// SetBlacklist can be used to change email blacklist during runtime
func SetBlacklist(s string) {
	blacklist = s
	blacklistRegex = regexp.MustCompile(blacklist)
}

// CanSendTo returns true if Fider is allowed to send email to given address
func CanSendTo(address string) bool {
	if strings.TrimSpace(address) == "" {
		return false
	}

	if whitelist != "" {
		return whitelistRegex.MatchString(address)
	}

	if blacklist != "" {
		return !blacklistRegex.MatchString(address)
	}

	return true
}
