package netguard_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/netguard"
)

func TestIsBlockedIP(t *testing.T) {
	RegisterT(t)

	testCases := []struct {
		ip      string
		blocked bool
	}{
		// IPv4 private / internal
		{"0.0.0.0", true},
		{"0.1.2.3", true},
		{"10.0.0.1", true},
		{"100.64.0.1", true},
		{"100.127.255.255", true},
		{"127.0.0.1", true},
		{"127.255.255.254", true},
		{"169.254.169.254", true},
		{"172.16.0.1", true},
		{"172.31.255.255", true},
		{"192.168.1.1", true},
		{"168.63.129.16", true}, // Azure WireServer
		{"192.0.0.0", true},
		{"192.0.0.1", true},
		{"192.0.0.7", true},
		{"192.0.0.8", true},   // IPv4 dummy address
		{"192.0.0.11", true},  // IETF protocol assignment, not globally reachable
		{"192.0.0.170", true}, // NAT64/DNS64 discovery
		{"192.0.0.171", true}, // NAT64/DNS64 discovery
		{"192.0.0.192", true}, // Oracle Cloud metadata
		{"192.0.0.255", true},
		{"::ffff:192.0.0.192", true},
		{"198.18.0.1", true},
		{"198.19.255.255", true},
		{"224.0.0.1", true},
		{"239.255.255.250", true},
		{"240.0.0.1", true},
		{"255.255.255.255", true},

		// IPv4-mapped IPv6
		{"::ffff:127.0.0.1", true},
		{"::ffff:169.254.169.254", true},
		{"::ffff:8.8.8.8", false},

		// IPv6 private / internal
		{"::", true},
		{"::1", true},
		{"fc00::1", true},
		{"fd12:3456::1", true},
		{"fe80::1", true},
		{"fec0::1", true},
		{"feff:ffff::1", true},
		{"ff02::1", true},
		{"ff0e::1", true},
		{"64:ff9b:1::1", true},
		{"100::1", true},                   // discard-only
		{"100::ffff:ffff:ffff:ffff", true}, // discard-only
		{"2001:db8::1", true},              // documentation
		{"2001:db8:ffff::1", true},         // documentation

		// NAT64 64:ff9b::/96
		{"64:ff9b::a9fe:a9fe", true},
		{"64:ff9b::7f00:1", true},
		{"64:ff9b::a00:1", true},
		{"64:ff9b::808:808", false},

		// 6to4 2002::/16
		{"2002:a9fe:a9fe::1", true},
		{"2002:7f00:1::1", true},
		{"2002:c0a8:101::1", true},
		{"2002:0808:0808::1", false},

		// Teredo 2001:0000::/32, client IPv4 = last 4 bytes XOR 0xff
		{"2001:0:4136:e378:8000:63bf:80ff:fffe", true},  // client 127.0.0.1
		{"2001:0:4136:e378:8000:63bf:5601:5601", true},  // client 169.254.169.254
		{"2001:0:4136:e378:8000:63bf:f7f7:f7f7", false}, // client 8.8.8.8

		// Deprecated IPv4-compatible
		{"::7f00:1", true},
		{"::a9fe:a9fe", true},

		// IPv4-translated / SIIT ::ffff:0:0:0/96
		{"::ffff:0:a9fe:a9fe", true}, // 169.254.169.254
		{"::ffff:0:7f00:1", true},    // 127.0.0.1
		{"::ffff:0:c000:c0", true},   // 192.0.0.192
		{"::ffff:0:808:808", false},  // 8.8.8.8

		// Public
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"203.0.113.1", false},
		{"100.63.255.255", false},
		{"100.128.0.0", false},
		{"198.20.0.1", false},
		{"192.0.0.9", false},  // PCP anycast (globally reachable)
		{"192.0.0.10", false}, // TURN anycast (globally reachable)
		{"::ffff:192.0.0.9", false},
		{"168.63.129.15", false},
		{"168.63.129.17", false},
		{"100:0:0:1::1", false}, // just outside 100::/64
		{"2001:db9::1", false},  // just outside 2001:db8::/32
		{"192.0.1.1", false},
		{"192.0.2.1", false},
		{"191.255.255.255", false},
		{"2001:4860:4860::8888", false},
		{"2606:4700:4700::1111", false}}

	for _, tc := range testCases {
		ip := net.ParseIP(tc.ip)
		if ip == nil {
			t.Fatalf("invalid test IP %s", tc.ip)
		}
		if got := netguard.IsBlockedIP(ip); got != tc.blocked {
			t.Errorf("IsBlockedIP(%s) = %v, expected %v", tc.ip, got, tc.blocked)
		}
	}

	Expect(netguard.IsBlockedIP(nil)).IsTrue()
}

func newLocalServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestClient_RejectsBlockedTargetsAtDialTime(t *testing.T) {
	RegisterT(t)

	server := newLocalServer(t)
	_, port, _ := net.SplitHostPort(server.Listener.Addr().String())

	for _, rawurl := range []string{
		server.URL,                       // literal 127.0.0.1
		"http://localhost:" + port + "/", // hostname resolving to loopback
	} {
		req, _ := http.NewRequestWithContext(context.Background(), "GET", rawurl, nil)
		res, err := netguard.Client.Do(req)
		if res != nil {
			_ = res.Body.Close()
		}
		if err == nil || !errors.Is(err, netguard.ErrBlockedAddress) {
			t.Errorf("%s: expected ErrBlockedAddress, got %v", rawurl, err)
		}
	}
}

func TestClientFor(t *testing.T) {
	RegisterT(t)

	original := env.Config.AllowPrivateNetworkTargets
	t.Cleanup(func() { env.Config.AllowPrivateNetworkTargets = original })

	env.Config.AllowPrivateNetworkTargets = false
	Expect(netguard.ClientFor(true) == netguard.Client).IsTrue()
	Expect(netguard.ClientFor(false) == http.DefaultClient).IsTrue()

	env.Config.AllowPrivateNetworkTargets = true
	Expect(netguard.ClientFor(true) == http.DefaultClient).IsTrue()
	Expect(netguard.ClientFor(false) == http.DefaultClient).IsTrue()
}
