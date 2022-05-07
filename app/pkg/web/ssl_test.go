package web

import (
	"context"
	"crypto/tls"
	"testing"
	"time"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/blob/fs"

	. "github.com/getfider/fider/app/pkg/assert"
)

// this handler relies on real DNS records on feedbacktest.goenning.net
// the correct subdomains are defined below to simulate a configuration match for tests
func mockGetTenantWithCorrectSubdomains(ctx context.Context, q *query.GetTenantByDomain) error {
	if q.Domain == "feedbacktest.goenning.net" {
		q.Result = &entity.Tenant{Name: "Feedback for goenning.net", Subdomain: "goenning"}
		return nil
	} else {
		return app.ErrNotFound
	}
}

// this handler relies on real DNS records on feedbacktest.goenning.net
// wrong subdomains are defined below to simulate a mismatch in configuration for tests
func mockGetTenantWithIncorrectSubdomains(ctx context.Context, q *query.GetTenantByDomain) error {
	if q.Domain == "feedbacktest.goenning.net" {
		q.Result = &entity.Tenant{Name: "Feedback for goenning.net", Subdomain: "demo"}
		return nil
	} else {
		return app.ErrNotFound
	}
}

func TestUseAutoCert_WhenCNAMEAreRegistered(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})
	bus.AddHandler(mockGetTenantWithCorrectSubdomains)

	manager, err := NewCertificateManager(context.Background(), "", "")
	Expect(err).IsNil()

	cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
		ServerName: "feedbacktest.goenning.net",
	})
	Expect(err.Error()).ContainsSubstring(`acme/autocert: unable to satisfy`)
	Expect(err.Error()).ContainsSubstring(`for domain "feedbacktest.goenning.net": no viable challenge type found`)
	Expect(cert).IsNil()

	// GetCertificate starts a fire and forget go routine to delete items from cache, give it 2sec to complete it
	time.Sleep(2 * time.Second)
}

func TestGetCertificate(t *testing.T) {
	RegisterT(t)

	var testCases = []struct {
		mode       string
		cert       string
		serverName string
	}{
		{"multi", "all-test-fider-io", ""},
		{"multi", "all-test-fider-io", "fider"},
		{"multi", "all-test-fider-io", "feedback.test.fider.io"},
		{"multi", "all-test-fider-io", "FEEDBACK.test.fider.io"},
		{"multi", "all-test-fider-io", "44.194.119.243"},
		{"single", "test-fider-io", "test.fider.io"},
		{"single", "test-fider-io", "fider.io"},
		{"single", "test-fider-io", "44.194.119.243"},
	}

	for _, testCase := range testCases {
		env.Config.HostMode = testCase.mode
		certFile := env.Path("/app/pkg/web/testdata/" + testCase.cert + ".crt")
		keyFile := env.Path("/app/pkg/web/testdata/" + testCase.cert + ".key")
		wildcardCert, _ := tls.LoadX509KeyPair(certFile, keyFile)

		manager, err := NewCertificateManager(context.Background(), certFile, keyFile)
		Expect(err).IsNil()
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: testCase.serverName,
		})

		Expect(err).IsNil()
		Expect(cert.Certificate).Equals(wildcardCert.Certificate)
	}
}

func TestGetCertificate_WhenCNAMEAreNotConfigured(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})
	bus.AddHandler(mockGetTenantWithCorrectSubdomains)

	manager, err := NewCertificateManager(context.Background(), "", "")
	Expect(err).IsNil()

	invalidServerNames := []string{"feedback.heyworld.com", "ideas.app.com"}

	for _, serverName := range invalidServerNames {
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: serverName,
		})
		Expect(err.Error()).ContainsSubstring(`no tenant found with cname ` + serverName)
		Expect(cert).IsNil()
	}
}

func TestGetCertificate_WhenCNAMEDoesntMatch(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})

	bus.AddHandler(mockGetTenantWithIncorrectSubdomains)

	manager, err := NewCertificateManager(context.Background(), "", "")
	Expect(err).IsNil()

	cert, err := manager.GetCertificate(&tls.ClientHelloInfo{ServerName: "feedbacktest.goenning.net"})
	Expect(err.Error()).ContainsSubstring("cname goenning.test.fider.io. (from feedbacktest.goenning.net) doesn't match configured host demo.test.fider.io")
	Expect(cert).IsNil()
}

func TestGetCertificate_ServerNameMatchesCertificate_ShouldReturnIt(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})
	bus.AddHandler(mockGetTenantWithCorrectSubdomains)

	certFile := env.Etc("dev-fider-io.crt")
	certKey := env.Etc("dev-fider-io.key")
	manager, err := NewCertificateManager(context.Background(), certFile, certKey)
	Expect(err).IsNil()

	serverNames := []string{"dev.fider.io", "feedback.dev.fider.io", "anything.dev.fider.io", "IDEAS.DEV.fider.io", ".feedback.DEV.fider.io"}

	for _, serverName := range serverNames {
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: serverName,
		})
		Expect(err).IsNil()
		Expect(cert).IsNotNil()
	}
}

func TestGetCertificate_ServerNameDoesntMatchCertificate_ButEndsWithHostName_ShouldThrow(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})
	bus.AddHandler(mockGetTenantWithCorrectSubdomains)

	env.Config.HostDomain = "dev.fider.io"
	certFile := env.Etc("dev-fider-io.crt")
	certKey := env.Etc("dev-fider-io.key")
	manager, err := NewCertificateManager(context.Background(), certFile, certKey)
	Expect(err).IsNil()

	serverNames := []string{"sub.feedback.dev.fider.io"}

	for _, serverName := range serverNames {
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: serverName,
		})
		Expect(err.Error()).ContainsSubstring("invalid ServerName used: " + serverName)
		Expect(cert).IsNil()
	}
}
