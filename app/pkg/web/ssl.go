package web

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
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
	"github.com/getfider/fider/app/pkg/log"
	"golang.org/x/crypto/acme/autocert"
)

func getDefaultTLSConfig(autoSSL bool) *tls.Config {
	nextProtos := []string{"h2", "http/1.1"}
	if autoSSL {
		nextProtos = append(nextProtos, acme.ALPNProto)
	}

	return &tls.Config{
		NextProtos: nextProtos,
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

var errInvalidHostName = errors.New("autotls: invalid hostname")

func isValidHostName(ctx context.Context, host string) error {
	// In this context, host can only be custom domains, not a subdomain of fider.io

	if env.IsSingleHostMode() {
		return nil
	}

	if host == "" {
		return errors.Wrap(errInvalidHostName, "host cannot be empty.")
	}

	trx, err := dbx.BeginTx(ctx)
	if err != nil {
		return errors.Wrap(err, "failed start new transaction")
	}
	defer trx.MustCommit()
	dbCtx := context.WithValue(ctx, app.TransactionCtxKey, trx)

	getTenant := &query.GetTenantByDomain{Domain: host}
	err = bus.Dispatch(dbCtx, getTenant)
	if err != nil {
		if errors.Cause(err) == app.ErrNotFound {
			return errors.Wrap(errInvalidHostName, "no tenant found with cname %s", host)
		}
		return errors.Wrap(err, "failed to get tenant by cname")
	}

	cname, err := net.DefaultResolver.LookupCNAME(ctx, host)
	if err != nil {
		return errors.Wrap(err, "failed to lookup CNAME")
	}

	if cname == "" {
		return errors.Wrap(errInvalidHostName, "no CNAME DNS record found for %s", host)
	}

	if strings.TrimSuffix(cname, ".") != getTenant.Result.Subdomain+env.MultiTenantDomain() {
		return errors.Wrap(errInvalidHostName, "cname %s (from %s) doesn't match configured host %s", cname, host, getTenant.Result.Subdomain+env.MultiTenantDomain())
	}

	return nil
}

//CertificateManager is used to manage SSL certificates
type CertificateManager struct {
	ctx     context.Context
	cert    tls.Certificate
	leaf    *x509.Certificate
	autotls autocert.Manager
}

//NewCertificateManager creates a new CertificateManager
func NewCertificateManager(ctx context.Context, certFile, keyFile string) (*CertificateManager, error) {
	manager := &CertificateManager{
		ctx: ctx,
		autotls: autocert.Manager{
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

		// If it's an IP address, just return the cert we have
		if net.ParseIP(serverName) != nil {
			return &m.cert, nil
		}

		// throw an error if it doesn't match the leaf certificate but still ends with current hostname, example:
		// hostdomain is myserver.com and the certificate is *.myserver.com
		// serverName is something.else.myserver.com, it should throw an error
		if strings.HasSuffix(serverName, "."+env.Config.HostDomain) {
			return nil, errors.New("invalid ServerName used: %s", serverName)
		}
	}

	//TODO: consider recovering from a possible panic here
	cert, err := m.autotls.GetCertificate(hello)
	if err != nil {
		if errors.Cause(err) == errInvalidHostName {
			log.Warn(m.ctx, err.Error())
		} else {
			log.Error(m.ctx, errors.Wrap(err, "failed to get certificate for %s", hello.ServerName))
		}
	}

	return cert, err
}

//StartHTTPServer creates a new HTTP server on port 80 that is used for the ACME HTTP Challenge
func (m *CertificateManager) StartHTTPServer() {
	err := http.ListenAndServe(":80", m.autotls.HTTPHandler(nil))
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
