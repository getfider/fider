// Package netguard contains the shared SSRF protection used when Fider makes
// outbound HTTP requests to user-configurable URLs (webhooks, custom OAuth
// token/profile URLs).
//
// There are two layers:
//   - validate.WebhookURL performs a friendly preflight check (it resolves the
//     hostname and rejects private targets with a readable message).
//   - The client returned by NewClient enforces the same IsBlockedIP check at
//     dial time, on the exact IP address being connected to. This is the real
//     enforcement: it cannot be bypassed by DNS rebinding (a hostname that
//     resolves to a public IP during validation and to a private IP when the
//     request is made).
package netguard

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/getfider/fider/app/pkg/env"
)

// ErrBlockedAddress is returned when a connection to a private/internal
// network address is attempted through a guarded client.
var ErrBlockedAddress = errors.New("connection to a private or internal network address is not allowed")

func mustParseCIDRs(cidrs ...string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		_, n, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		nets = append(nets, n)
	}
	return nets
}

var blockedIPv4 = mustParseCIDRs(
	"0.0.0.0/8",      // "this" network
	"10.0.0.0/8",     // private
	"100.64.0.0/10",  // carrier-grade NAT
	"127.0.0.0/8",    // loopback
	"169.254.0.0/16", // link-local (incl. cloud metadata 169.254.169.254)
	"172.16.0.0/12",  // private
	"192.0.0.0/24",   // IETF protocol assignments (incl. Oracle Cloud metadata 192.0.0.192)
	"192.168.0.0/16", // private
	"198.18.0.0/15",  // benchmarking
	"224.0.0.0/4",    // multicast
	"240.0.0.0/4",    // reserved, incl. broadcast 255.255.255.255
)

var blockedIPv6 = mustParseCIDRs(
	"::/128",         // unspecified
	"::1/128",        // loopback
	"fc00::/7",       // unique local
	"fe80::/10",      // link-local
	"fec0::/10",      // site-local (deprecated)
	"ff00::/8",       // multicast
	"64:ff9b:1::/48", // NAT64 local-use
)

var (
	nat64Prefix     = mustParseCIDRs("64:ff9b::/96")[0]
	sixToFourPrefix = mustParseCIDRs("2002::/16")[0]
	teredoPrefix    = mustParseCIDRs("2001::/32")[0]
	ipv4Compatible  = mustParseCIDRs("::/96")[0]           // deprecated IPv4-compatible IPv6 (::a.b.c.d)
	ipv4Translated  = mustParseCIDRs("::ffff:0:0:0/96")[0] // IPv4-translated / SIIT (::ffff:0:a.b.c.d)
)

// IsBlockedIP reports whether ip belongs to a private, internal, or otherwise
// non-public range that outbound requests to user-configured URLs must not reach.
// IPv6 transition addresses (NAT64, 6to4, Teredo, IPv4-compatible, SIIT) are
// unwrapped and the embedded IPv4 address is checked as well.
// A nil/invalid IP is treated as blocked.
func IsBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	// Covers plain IPv4 and IPv4-mapped IPv6 (::ffff:a.b.c.d).
	if v4 := ip.To4(); v4 != nil {
		for _, n := range blockedIPv4 {
			if n.Contains(v4) {
				return true
			}
		}
		return false
	}

	ip16 := ip.To16()
	if ip16 == nil {
		return true
	}

	for _, n := range blockedIPv6 {
		if n.Contains(ip16) {
			return true
		}
	}

	if embedded := embeddedIPv4(ip16); embedded != nil {
		return IsBlockedIP(embedded)
	}

	return false
}

// embeddedIPv4 extracts the IPv4 address embedded in IPv6 transition
// addresses, or returns nil when ip is not such an address.
func embeddedIPv4(ip net.IP) net.IP {
	switch {
	case nat64Prefix.Contains(ip):
		// 64:ff9b::a.b.c.d
		return net.IPv4(ip[12], ip[13], ip[14], ip[15])
	case sixToFourPrefix.Contains(ip):
		// 2002:AABB:CCDD::/48
		return net.IPv4(ip[2], ip[3], ip[4], ip[5])
	case teredoPrefix.Contains(ip):
		// 2001:0000:<server v4>:<flags>:<port>:<client v4 XOR 0xffffffff>
		return net.IPv4(ip[12]^0xff, ip[13]^0xff, ip[14]^0xff, ip[15]^0xff)
	case ipv4Compatible.Contains(ip):
		// ::a.b.c.d (:: and ::1 are already handled by blockedIPv6)
		return net.IPv4(ip[12], ip[13], ip[14], ip[15])
	case ipv4Translated.Contains(ip):
		// ::ffff:0:a.b.c.d
		return net.IPv4(ip[12], ip[13], ip[14], ip[15])
	}
	return nil
}

// dialControl runs after DNS resolution, immediately before each connect()
// syscall, with the literal IP:port being connected to. It therefore also
// covers every address tried by happy-eyeballs / fallback dialing.
func dialControl(network, address string, _ syscall.RawConn) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrBlockedAddress, address)
	}
	if IsBlockedIP(net.ParseIP(host)) {
		return fmt.Errorf("%w: %s", ErrBlockedAddress, host)
	}
	return nil
}

// NewClient returns an http.Client that refuses to connect to any address for
// which IsBlockedIP returns true. It mirrors Fider's default client settings:
// 30s timeout and redirects are not followed.
//
// Proxies are intentionally disabled (Proxy: nil, i.e. HTTP(S)_PROXY is
// ignored). With a proxy, the dial would go to the proxy rather than to the
// target, so the dial-time check would inspect the proxy's address and either
// block every request (proxy on a private network) or protect nothing (the
// proxy resolves and connects to the target itself). Guarded requests are
// therefore always sent directly.
func NewClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Control:   dialControl,
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = dialer.DialContext

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// Client is a shared guarded client, safe for concurrent use.
var Client = NewClient()

// ClientFor returns the HTTP client to use for an outbound request. When
// blockPrivateNetworkTargets is true it returns the guarded Client, unless the
// instance opted out via ALLOW_PRIVATE_NETWORK_TARGETS=true, in which case (and
// when blockPrivateNetworkTargets is false) it returns http.DefaultClient.
// This is the single place where the escape hatch is applied.
func ClientFor(blockPrivateNetworkTargets bool) *http.Client {
	if blockPrivateNetworkTargets && !env.Config.AllowPrivateNetworkTargets {
		return Client
	}
	return http.DefaultClient
}
