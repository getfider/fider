package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/getfider/fider/app/pkg/env"
)

//RunPing checks if Fider Server is running and is healthy
//Returns an exitcode, 0 for OK and 1 for ERROR
func RunPing() int {

	client := &http.Client{}

	host := "localhost:" + env.Config.Port
	protocol := "http://"
	if env.Config.TLS.Certificate != "" || env.Config.TLS.Automatic {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		// Use the site's hostname for SNI instead of "localhost". With SSL_AUTO the server can
		// only get certificates for real domain names, so a handshake for "localhost" fails with
		// "acme/autocert: server name component count invalid".
		if env.IsSingleHostMode() {
			if u, err := url.Parse(env.Config.BaseURL); err == nil {
				tlsConfig.ServerName = u.Hostname()
			}
		}
		http.DefaultTransport.(*http.Transport).TLSClientConfig = tlsConfig
		protocol = "https://"
	}

	req, err := http.NewRequest("GET", protocol+host+"/_health", nil)
	if err != nil {
		fmt.Printf("Setting up request failed with: %s", err)
		return 1
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed with: %s\n", err)
		return 1
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 400 {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return 1
	}

	return 0
}
