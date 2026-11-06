package cmd

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/env"
)

func TestRunPing_AutoTLSUsesBaseURLHostname(t *testing.T) {
	RegisterT(t)

	originalConfig := env.Config
	transport := http.DefaultTransport.(*http.Transport)
	originalTLSConfig := transport.TLSClientConfig
	defer func() {
		env.Config = originalConfig
		transport.TLSClientConfig = originalTLSConfig
	}()

	var serverNames []string
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	srv.TLS = &tls.Config{}
	// Like autocert, only hand out a certificate for a real domain name
	srv.TLS.GetCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		serverNames = append(serverNames, hello.ServerName)
		if !strings.Contains(hello.ServerName, ".") {
			return nil, errors.New("acme/autocert: server name component count invalid")
		}
		return &srv.TLS.Certificates[0], nil
	}
	srv.StartTLS()
	defer srv.Close()

	_, port, err := net.SplitHostPort(srv.Listener.Addr().String())
	Expect(err).IsNil()

	env.Config.HostMode = "single"
	env.Config.BaseURL = "https://feedback.example.com"
	env.Config.Port = port
	env.Config.TLS.Automatic = true
	env.Config.TLS.Certificate = ""

	Expect(RunPing()).Equals(0)
	Expect(serverNames).Equals([]string{"feedback.example.com"})
}
