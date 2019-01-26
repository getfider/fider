package webutil_test

import (
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/di"

	"github.com/getfider/fider/app/models"
	"github.com/getfider/fider/app/pkg/env"
	webutil "github.com/getfider/fider/app/pkg/web/util"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/web"
)

func newContext(rawurl string) web.Context {
	url, _ := url.Parse(rawurl)

	return web.Context{
		Request: web.Request{
			URL: url,
		},
	}
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

	img, _ := ioutil.ReadFile(env.Path("/app/pkg/img/testdata/logo3-200w.gif"))
	services := &app.Services{
		Blobs: di.NewBlobStorage(nil),
	}
	ctx := newContext("http://demo.test.fider.io:3000/hello-world")
	ctx.SetServices(services)

	upload := &models.ImageUpload{
		Upload: &models.ImageUploadData{
			Content:     img,
			ContentType: "image/gif",
		},
	}
	err := webutil.ProcessImageUpload(ctx, upload, "attachments")
	Expect(err).IsNil()
	Expect(upload.BlobKey).ContainsSubstring("attachments/")

	blob, err := services.Blobs.Get(upload.BlobKey)
	Expect(err).IsNil()
	Expect(blob.ContentType).Equals("image/gif")
	Expect(blob.Object).Equals(img)
	Expect(int(blob.Size)).Equals(len(img))
}

func TestMultiProcessImageUpload(t *testing.T) {
	RegisterT(t)

	img, _ := ioutil.ReadFile(env.Path("/app/pkg/img/testdata/logo3-200w.gif"))
	services := &app.Services{
		Blobs: di.NewBlobStorage(nil),
	}
	ctx := newContext("http://demo.test.fider.io:3000/hello-world")
	ctx.SetServices(services)

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
