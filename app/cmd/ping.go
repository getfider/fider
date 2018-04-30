package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/getfider/fider/app/pkg/env"
)

//RunPing checks if Fider Server is running and is healthy
//Returns an exitcode, 0 for OK and 1 for ERROR
func RunPing() int {
	var client = &http.Client{}

	protocol := "http://"
	if env.IsDefined("SSL_CERT") || env.IsDefined("SSL_AUTO") {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		protocol = "https://"
	}

	resp, err := client.Get(protocol + "localhost:3000/-/health")
	if err != nil {
		fmt.Printf("Request failed with: %s\n", err)
		return 1
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return 1
	}

	return 0
}
