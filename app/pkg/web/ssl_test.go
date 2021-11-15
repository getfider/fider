package web

import (
	"context"
	"crypto/tls"
	"testing"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/services/blob/fs"

	. "github.com/getfider/fider/app/pkg/assert"
)

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
		{"single", "test-fider-io", "test.fider.io"},
		{"single", "test-fider-io", "fider.io"},
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

func TestGetCertificate_WhenCNAMEAreInvalid(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})
	bus.AddHandler(mockIsCNAMEAvailable)

	manager, err := NewCertificateManager(context.Background(), "", "")
	Expect(err).IsNil()

	invalidServerNames := []string{"feedback.heyworld.com", "2.2.2.2"}

	for _, serverName := range invalidServerNames {
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: serverName,
		})
		Expect(err.Error()).ContainsSubstring(`no tenants found with cname ` + serverName)
		Expect(cert).IsNil()
	}
}

func TestGetCertificate_ServerNameMatchesCertificate_ShouldReturnIt(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})
	bus.AddHandler(mockIsCNAMEAvailable)

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
	bus.AddHandler(mockIsCNAMEAvailable)

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
