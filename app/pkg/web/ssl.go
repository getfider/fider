package web

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

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

// certFailureCooldown is how long a host is skipped after a failed certificate acquisition.
// Let's Encrypt allows roughly 5 failed validations per hostname per hour, so this keeps us
// to at most 4, leaving headroom for the rest of the account.
const certFailureCooldown = 15 * time.Minute

// hostFailureCache records hosts whose certificate acquisition recently failed, so that we
// stop asking Let's Encrypt for them until the cooldown elapses. Without it a single broken
// custom domain retries on every TLS handshake, which rate limits the whole ACME account.
type hostFailureCache struct {
	mu       sync.Mutex
	failures map[string]time.Time
}

func newHostFailureCache() *hostFailureCache {
	return &hostFailureCache{failures: make(map[string]time.Time)}
}

// normalizeHost matches how autocert canonicalizes the name it hands to the HostPolicy, so
// that failures recorded from a raw ServerName are found again on the next handshake.
func normalizeHost(host string) string {
	return strings.TrimSuffix(strings.ToLower(host), ".")
}

func (c *hostFailureCache) onCooldown(host string) bool {
	key := normalizeHost(host)

	c.mu.Lock()
	defer c.mu.Unlock()

	failedAt, ok := c.failures[key]
	if !ok {
		return false
	}

	if time.Since(failedAt) >= certFailureCooldown {
		delete(c.failures, key)
		return false
	}

	return true
}

func (c *hostFailureCache) recordFailure(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.failures[normalizeHost(host)] = time.Now()
}

func (c *hostFailureCache) clear(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.failures, normalizeHost(host))
}

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
		return errors.Wrap(errInvalidHostName, "failed to lookup CNAME")
	}

	if cname == "" {
		return errors.Wrap(errInvalidHostName, "no CNAME DNS record found for %s", host)
	}

	if strings.TrimSuffix(cname, ".") != getTenant.Result.Subdomain+env.MultiTenantDomain() {
		return errors.Wrap(errInvalidHostName, "cname %s (from %s) doesn't match configured host %s", cname, host, getTenant.Result.Subdomain+env.MultiTenantDomain())
	}

	return nil
}

// CertificateManager is used to manage SSL certificates
type CertificateManager struct {
	ctx      context.Context
	cert     tls.Certificate
	leaf     *x509.Certificate
	autotls  autocert.Manager
	failures *hostFailureCache
}

// NewCertificateManager creates a new CertificateManager
func NewCertificateManager(ctx context.Context, certFile, keyFile string) (*CertificateManager, error) {
	manager := &CertificateManager{
		ctx:      ctx,
		failures: newHostFailureCache(),
		autotls: autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  NewAutoCertCache(),
			Client: acmeClient(),
		},
	}
	manager.autotls.HostPolicy = manager.hostPolicy

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

// hostPolicy is the autocert HostPolicy. autocert calls it before looking up or issuing a
// certificate, but after serving TLS-ALPN challenge tokens, so refusing a host here stops us
// from contacting Let's Encrypt without interfering with an in-flight challenge validation.
func (m *CertificateManager) hostPolicy(ctx context.Context, host string) error {
	if m.failures.onCooldown(host) {
		return errors.Wrap(errInvalidHostName, "host %s is on cooldown after a recent certificate failure", host)
	}
	return isValidHostName(ctx, host)
}

// GetCertificate decides which certificate to use
// It first tries to use loaded certificate for incoming request if it's compatible
// Otherwise fallsback to a automatically generated certificate by Let's Encrypt
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
			// In multi-tenant mode, CNAME errors are expected customer DNS misconfigurations
			// In single-tenant mode, they indicate actual setup issues the admin should know about
			if env.IsSingleHostMode() {
				log.Warn(m.ctx, err.Error())
			} else {
				log.Debug(m.ctx, err.Error())
			}
		} else {
			// Put the host on cooldown so a domain that can't be issued doesn't retry on every
			// handshake and get the ACME account rate limited.
			m.failures.recordFailure(hello.ServerName)

			failure := errors.Wrap(err, "failed to get certificate for %s", hello.ServerName)
			// Same rationale as above: in multi-tenant mode a custom domain that fails ACME is
			// the customer's problem to fix, not a Fider bug, so keep it out of error alarms.
			if env.IsSingleHostMode() {
				log.Error(m.ctx, failure)
			} else {
				log.Warn(m.ctx, failure.Error())
			}
		}
	} else {
		m.failures.clear(hello.ServerName)
	}

	return cert, err
}

// StartHTTPServer creates a new HTTP server on port 80 that is used for the ACME HTTP Challenge
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
