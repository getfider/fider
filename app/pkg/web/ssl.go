package web

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"

	"golang.org/x/crypto/acme"
	"golang.org/x/net/idna"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/dbx"
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

func isValidHostName(ctx context.Context, host string) error {
	if host == "" {
		return errors.New("host cannot be empty.")
	}

	if env.IsSingleHostMode() {
		if env.Config.HostDomain == host {
			return nil
		}
		return errors.New("server name mismatch")
	}

	trx, err := dbx.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed start new transaction")
	}
	defer trx.MustCommit()

	isAvailable := &query.IsCNAMEAvailable{CNAME: host}
	newCtx := context.WithValue(ctx, app.TransactionCtxKey, trx)
	if err := bus.Dispatch(newCtx, isAvailable); err != nil {
		return errors.Wrap(err, "failed to find tenant by cname")
	}

	if isAvailable.Result {
		return errors.New("no tenants found with cname %s", host)
	}
	return nil
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
			Prompt:     autocert.AcceptTOS,
			Cache:      NewAutoCertCache(),
			Client:     acmeClient(),
			HostPolicy: isValidHostName,
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
		serverName, err := idna.Lookup.ToASCII(hello.ServerName)
		if err != nil {
			return nil, err
		}
		serverName = strings.Trim(serverName, ".")

		// If ServerName is empty or does't contain a dot, just return the certificate
		if serverName == "" || !strings.Contains(serverName, ".") {
			return &m.cert, nil
		}

		if env.IsSingleHostMode() || m.leaf.VerifyHostname(serverName) == nil {
			return &m.cert, nil
		}

		// throw an error if it doesn't match the leaf certificate but still ends with current hostname, example:
		// hostdomain is myserver.com and the certificate is *.myserver.com
		// serverName is something.else.myserver.com, it should throw an error
		if strings.HasSuffix(serverName, "."+env.Config.HostDomain) {
			return nil, errors.New("invalid ServerName used: %s", serverName)
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
