package webutil_test

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/services/blob/fs"

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

func TestProcessImageUpload(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})

	img, _ := ioutil.ReadFile(env.Path("/app/pkg/web/testdata/logo3-200w.gif"))
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

	q := &query.GetBlobByKey{Key: upload.BlobKey}
	err = bus.Dispatch(context.Background(), q)
	Expect(err).IsNil()
	Expect(q.Result.ContentType).Equals("image/gif")
	Expect(q.Result.Content).Equals(img)
	Expect(int(q.Result.Size)).Equals(len(img))
}

func TestMultiProcessImageUpload(t *testing.T) {
	RegisterT(t)
	bus.Init(fs.Service{})

	img, _ := ioutil.ReadFile(env.Path("/app/pkg/web/testdata/logo3-200w.gif"))
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
