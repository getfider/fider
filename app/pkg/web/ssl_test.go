package web

import (
	"crypto/tls"
	"testing"

	"github.com/getfider/fider/app/pkg/env"

	. "github.com/getfider/fider/app/pkg/assert"
)

func Test_GetCertificate(t *testing.T) {
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

		manager, err := NewCertificateManager(certFile, keyFile)
		Expect(err).IsNil()
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: testCase.serverName,
		})

		Expect(err).IsNil()
		Expect(cert.Certificate).Equals(wildcardCert.Certificate)
	}
}
