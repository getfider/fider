// Package netguard contains the shared SSRF protection used when Fider makes
// outbound HTTP requests to user-configurable URLs (webhooks, custom OAuth
// token/profile URLs).
//
// There are two layers:
//   - validate.WebhookURL performs a friendly preflight check (it resolves the
//     hostname and rejects private targets with a readable message).
//   - The client returned by Client/ClientFor enforces the same IsBlockedIP
//     check at dial time, on the exact IP address being connected to. This is
//     the real enforcement: it cannot be bypassed by DNS rebinding (a hostname
//     that resolves to a public IP during validation and to a private IP when
//     the request is made). The exception is requests sent through an HTTP(S)
//     proxy when SSRF_GUARD_USE_PROXY=true; see NewClient.
package netguard

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
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

// ProxyFunc selects the proxy for a request, like http.Transport.Proxy.
// A nil ProxyFunc means requests are always sent directly.
type ProxyFunc func(*http.Request) (*url.URL, error)

// proxyAddrKey is the context key under which guardedTransport stores the
// "host:port" of the proxy chosen for a request.
type proxyAddrKey struct{}

// guardedTransport wraps the real transport to tell the dialer which address,
// if any, is the proxy for this particular request.
//
// With a proxy, net/http dials the proxy rather than the target, so a
// dial-time check can only ever see the proxy's address. The proxy itself is
// operator-configured and may well live on a private network, so that one
// connection must be allowed. Everything else, in particular requests for
// which the proxy func returns nil (NO_PROXY match, loopback targets, or no
// proxy configured), must still be checked. To do that without guessing, the
// wrapper evaluates the proxy func itself and stores the resulting proxy
// address in the request context; the dialer skips the IP check only when the
// address it is asked to dial equals that value exactly. If the two ever
// disagree (e.g. the proxy environment changed between calls), the dial is
// checked as normal, i.e. the guard fails closed.
type guardedTransport struct {
	proxy ProxyFunc
	base  *http.Transport
}

func (t *guardedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.proxy != nil {
		proxyURL, err := t.proxy(req)
		if err != nil {
			return nil, err
		}
		if proxyURL != nil {
			ctx := context.WithValue(req.Context(), proxyAddrKey{}, proxyDialAddr(proxyURL))
			req = req.WithContext(ctx)
		}
	}
	return t.base.RoundTrip(req)
}

// proxyDialAddr returns the address net/http dials for proxyURL (mirrors
// net/http's canonicalAddr: default port by scheme).
func proxyDialAddr(proxyURL *url.URL) string {
	port := proxyURL.Port()
	if port == "" {
		switch strings.ToLower(proxyURL.Scheme) {
		case "https":
			port = "443"
		case "socks5", "socks5h":
			port = "1080"
		default:
			port = "80"
		}
	}
	return net.JoinHostPort(proxyURL.Hostname(), port)
}

// NewClient returns an http.Client that refuses to connect to any address for
// which IsBlockedIP returns true. It mirrors Fider's default client settings:
// 30s timeout and redirects are not followed.
//
// If proxy is nil, requests are always sent directly and every connection is
// checked. If proxy is set, requests it routes through a proxy are NOT covered
// by the dial-time check (the proxy resolves and connects to the target), so
// for those only the validate.WebhookURL preflight applies; this is weaker
// against DNS rebinding. Requests the proxy func sends directly are still
// fully checked. See guardedTransport.
func NewClient(proxy ProxyFunc) *http.Client {
	guarded := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Control:   dialControl,
	}
	toProxy := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	base := http.DefaultTransport.(*http.Transport).Clone()
	base.Proxy = proxy
	base.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if proxyAddr, _ := ctx.Value(proxyAddrKey{}).(string); proxyAddr != "" && addr == proxyAddr {
			return toProxy.DialContext(ctx, network, addr)
		}
		return guarded.DialContext(ctx, network, addr)
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: &guardedTransport{proxy: proxy, base: base},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

var (
	directOnce    sync.Once
	directClient  *http.Client
	proxiedOnce   sync.Once
	proxiedClient *http.Client
)

// Client returns the shared guarded client (safe for concurrent use).
// It is chosen when called, not at package init, so it reflects the loaded
// env config:
//   - SSRF_GUARD_USE_PROXY=false (default): HTTP(S)_PROXY is ignored and every
//     connection is checked at dial time.
//   - SSRF_GUARD_USE_PROXY=true: http.ProxyFromEnvironment is used; proxied
//     requests rely on the preflight only, direct ones are checked at dial time.
func Client() *http.Client {
	if env.Config.SSRFGuardUseProxy {
		proxiedOnce.Do(func() { proxiedClient = NewClient(http.ProxyFromEnvironment) })
		return proxiedClient
	}
	directOnce.Do(func() { directClient = NewClient(nil) })
	return directClient
}

// ClientFor returns the HTTP client to use for an outbound request. When
// blockPrivateNetworkTargets is true it returns the guarded Client(), unless
// the instance opted out via ALLOW_PRIVATE_NETWORK_TARGETS=true, in which case
// (and when blockPrivateNetworkTargets is false) it returns http.DefaultClient.
// This is the single place where the escape hatch is applied.
func ClientFor(blockPrivateNetworkTargets bool) *http.Client {
	if blockPrivateNetworkTargets && !env.Config.AllowPrivateNetworkTargets {
		return Client()
	}
	return http.DefaultClient
}
