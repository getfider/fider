package web

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme"

	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
)

func getDefaultTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		},
	}
}

//CertificateManager is used to manage SSL certificates
type CertificateManager struct {
	cert    tls.Certificate
	leaf    *x509.Certificate
	autossl autocert.Manager
}

//NewCertificateManager creates a new CertificateManager
func NewCertificateManager(certFile, keyFile string) (*CertificateManager, error) {
	manager := &CertificateManager{
		autossl: autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  NewAutoCertCache(),
			Client: acmeClient(),
		},
	}

	if certFile != "" && keyFile != "" {
		var err error
		manager.cert, err = tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load X509KeyPair for %s and %s", certFile, keyFile)
		}

		manager.leaf, err = x509.ParseCertificate(manager.cert.Certificate[0])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse x509 certificate")
		}
	}

	return manager, nil
}

//GetCertificate decides which certificate to use
//It first tries to use loaded certificate for incoming request if it's compatible
//Otherwise fallsback to a automatically generated certificate by Let's Encrypt
func (m *CertificateManager) GetCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if m.leaf != nil {
		serverName := strings.Trim(strings.ToLower(hello.ServerName), ".")

		// If ServerName is empty or does't contain a dot, just return the certificate
		if serverName == "" || !strings.Contains(serverName, ".") {
			return &m.cert, nil
		}

		if env.IsSingleHostMode() || m.leaf.VerifyHostname(serverName) == nil {
			return &m.cert, nil
		}
	}

	return m.autossl.GetCertificate(hello)
}

//StartHTTPServer creates a new HTTP server on port 80 that is used for the ACME HTTP Challenge
func (m *CertificateManager) StartHTTPServer() {
	err := http.ListenAndServe(":80", m.autossl.HTTPHandler(nil))
	if err != nil {
		panic(err)
	}
}

func acmeClient() *acme.Client {
	if env.IsTest() {
		return &acme.Client{
			DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
		}
	}
	return nil
}
