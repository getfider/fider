package web

import (
	"crypto/tls"
	"testing"

	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/services/blob/fs"

	. "github.com/getfider/fider/app/pkg/assert"
)

func Test_UseAutoCert(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})

	manager, err := NewCertificateManager("", "")
	Expect(err).IsNil()

	invalidServerNames := []string{"ideas.app.com", "feedback.mysite.com"}

	for _, serverName := range invalidServerNames {
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: serverName,
		})
		Expect(err.Error()).ContainsSubstring(`acme/autocert: unable to authorize "` + serverName + `"; challenge "tls-alpn-01" failed with error: acme: authorization error`)
		Expect(cert).IsNil()
	}
}
