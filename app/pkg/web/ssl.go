package web

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/getfider/fider/app/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
)

//CertificateManager is used to manage SSL certificates
type CertificateManager struct {
	cert    tls.Certificate
	leaf    *x509.Certificate
	autossl autocert.Manager
}

//NewCertificateManager creates a new CertificateManager
func NewCertificateManager(certFile, keyFile, cacheDir string) (*CertificateManager, error) {
	manager := &CertificateManager{
		autossl: autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache(cacheDir),
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
		if hello.ServerName == "" || m.leaf.VerifyHostname(hello.ServerName) == nil {
			return &m.cert, nil
		}
	}
	return m.autossl.GetCertificate(hello)
}

//StartHTTPServer creates a new HTTP server on port 80 that is used for the ACME HTTP Challenge
func (m *CertificateManager) StartHTTPServer() {
	http.ListenAndServe(":80", m.autossl.HTTPHandler(nil))
}
