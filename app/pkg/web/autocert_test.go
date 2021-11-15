package web

import (
	"context"
	"crypto/tls"
	"strings"
	"testing"
	"time"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/services/blob/fs"

	. "github.com/getfider/fider/app/pkg/assert"
)

func mockIsCNAMEAvailable(ctx context.Context, q *query.IsCNAMEAvailable) error {
	if q.CNAME == "ideas.app.com" || q.CNAME == "feedback.mysite.com" {
		q.Result = false
	} else {
		q.Result = true
	}
	return nil
}

func TestUseAutoCert_WhenCNAMEAreRegistered(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})
	bus.AddHandler(mockIsCNAMEAvailable)

	manager, err := NewCertificateManager(context.Background(), "", "")
	Expect(err).IsNil()

	invalidServerNames := []string{"ideas.app.com", "FEEDBACK.mysite.COM"}

	for _, serverName := range invalidServerNames {
		cert, err := manager.GetCertificate(&tls.ClientHelloInfo{
			ServerName: serverName,
		})
		Expect(err.Error()).ContainsSubstring(`acme/autocert: unable to satisfy`)
		Expect(err.Error()).ContainsSubstring(`for domain "` + strings.ToLower(serverName) + `": no viable challenge type found`)
		Expect(cert).IsNil()
	}

	// GetCertificate starts a fire and forget go routine to delete items from cache, give it 2sec to complete it
	time.Sleep(2 * time.Second)
}
