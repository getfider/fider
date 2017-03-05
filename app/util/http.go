package util

import "strings"

//StripPort removes port from host and returns it
func StripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	return hostport[:colon]
}
