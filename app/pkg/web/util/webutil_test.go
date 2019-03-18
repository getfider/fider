package webutil_test

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/getfider/fider/app/services/blob/s3"

	"github.com/getfider/fider/app/services/blob"

	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	webutil "github.com/getfider/fider/app/pkg/web/util"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func newContext(rawurl string) *web.Context {
	u, _ := url.Parse(rawurl)
	engine := web.New(&models.SystemSettings{})

	req := &http.Request{
		Host:       u.Host,
		RequestURI: u.RequestURI(),
	}
	if u.Scheme == "https" {
		req.TLS = &tls.ConnectionState{}
	}
	res := httptest.NewRecorder()
	return web.NewContext(engine, req, res, make(web.StringMap))
}

func TestGetOAuthBaseURL(t *testing.T) {
	RegisterT(t)

	ctx := newContext("https://mydomain.com/hello-world")

	env.Config.HostMode = "multi"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("https://login.test.fider.io")

	env.Config.HostMode = "single"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("https://mydomain.com")
}

func TestGetOAuthBaseURL_WithPort(t *testing.T) {
	RegisterT(t)

	ctx := newContext("http://demo.test.fider.io:3000/hello-world")

	env.Config.HostMode = "multi"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("http://login.test.fider.io:3000")

	env.Config.HostMode = "single"
	Expect(webutil.GetOAuthBaseURL(ctx)).Equals("http://demo.test.fider.io:3000")
}

func TestProcessImageUpload(t *testing.T) {
	RegisterT(t)
	bus.Init(s3.Service{})

	img, _ := ioutil.ReadFile(env.Path("/app/pkg/img/testdata/logo3-200w.gif"))
	ctx := newContext("http://demo.test.fider.io:3000/hello-world")

	upload := &models.ImageUpload{
		Upload: &models.ImageUploadData{
			Content:     img,
			ContentType: "image/gif",
		},
	}
	err := webutil.ProcessImageUpload(ctx, upload, "attachments")
	Expect(err).IsNil()
	Expect(upload.BlobKey).ContainsSubstring("attachments/")

	retrieveCmd := &blob.RetrieveBlob{Key: upload.BlobKey}
	err = bus.Dispatch(context.Background(), retrieveCmd)
	Expect(err).IsNil()
	Expect(retrieveCmd.Blob.ContentType).Equals("image/gif")
	Expect(retrieveCmd.Blob.Content).Equals(img)
	Expect(int(retrieveCmd.Blob.Size)).Equals(len(img))
}

func TestMultiProcessImageUpload(t *testing.T) {
	RegisterT(t)
	bus.Init(s3.Service{})

	img, _ := ioutil.ReadFile(env.Path("/app/pkg/img/testdata/logo3-200w.gif"))
	ctx := newContext("http://demo.test.fider.io:3000/hello-world")

	upload1 := &models.ImageUpload{
		Upload: &models.ImageUploadData{
			Content:     img,
			ContentType: "image/gif",
		},
	}
	upload2 := &models.ImageUpload{
		Upload: &models.ImageUploadData{
			Content:     img,
			ContentType: "image/gif",
		},
	}
	err := webutil.ProcessMultiImageUpload(ctx, []*models.ImageUpload{upload1, upload2}, "attachments")
	Expect(err).IsNil()
	Expect(upload1.BlobKey).ContainsSubstring("attachments/")
	Expect(upload2.BlobKey).ContainsSubstring("attachments/")
}
