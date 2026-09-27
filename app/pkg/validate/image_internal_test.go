package validate

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"testing"

	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/goenning/imagic"
)

func TestSafeApply_LimitsConcurrentDecodes(t *testing.T) {
	RegisterT(t)

	// Occupy every decode slot
	for i := 0; i < maxConcurrentDecodes; i++ {
		decodeSemaphore <- struct{}{}
	}
	defer func() {
		for i := 0; i < maxConcurrentDecodes; i++ {
			<-decodeSemaphore
		}
	}()

	// A new decode has to wait for a slot, and gives up when the request is cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var buf bytes.Buffer
	Expect(png.Encode(&buf, image.NewGray(image.Rect(0, 0, 100, 100)))).IsNil()
	result, err := SafeApply(ctx, buf.Bytes(), imagic.Resize(50))
	Expect(err).Equals(context.Canceled)
	Expect(result).IsNil()
}
