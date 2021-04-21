package web

import (
	"crypto/tls"
	"testing"
	"time"

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
		Expect(err.Error()).ContainsSubstring(`acme/autocert: unable to satisfy`)
		Expect(err.Error()).ContainsSubstring(`for domain "` + serverName + `": no viable challenge type found`)
		Expect(cert).IsNil()
	}

	// GetCertificate starts a fire and forget go routine to delete items from cache, give it 2sec to complete it
	time.Sleep(2*time.Second)
}
