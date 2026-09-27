package validate

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
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

// encodeJPEG returns a small JPEG. Go's encoder always writes baseline (SOF0) JPEGs.
func encodeJPEG() []byte {
	var buf bytes.Buffer
	Expect(jpeg.Encode(&buf, image.NewYCbCr(image.Rect(0, 0, 16, 16), image.YCbCrSubsampleRatio444), nil)).IsNil()
	return buf.Bytes()
}

// withSOF returns a copy of the given JPEG with its SOF0 marker replaced by the given marker
func withSOF(content []byte, marker byte) []byte {
	content = append([]byte{}, content...)
	sof := bytes.Index(content, []byte{0xFF, 0xC0})
	Expect(sof > 0).IsTrue()
	content[sof+1] = marker
	return content
}

func TestIsProgressiveJPEG(t *testing.T) {
	RegisterT(t)

	baseline := encodeJPEG()
	Expect(isProgressiveJPEG(baseline)).IsFalse()
	Expect(isProgressiveJPEG(withSOF(baseline, 0xC1))).IsFalse() // extended sequential

	// Progressive and other non-sequential SOF markers
	for _, marker := range []byte{0xC2, 0xC3, 0xC6, 0xCA, 0xCE} {
		Expect(isProgressiveJPEG(withSOF(baseline, marker))).IsTrue()
	}

	// Go reads the header of a patched SOF2 JPEG just fine
	_, format, err := image.DecodeConfig(bytes.NewReader(withSOF(baseline, 0xC2)))
	Expect(err).IsNil()
	Expect(format).Equals("jpeg")

	// Fill bytes, standalone markers (RST0, TEM) and APPn segments before the SOF are skipped
	withExtras := []byte{0xFF, 0xD8, 0xFF, 0xFF, 0xD0, 0xFF, 0x01, 0xFF, 0xE1, 0x00, 0x06, 'E', 'x', 'i', 'f'}
	Expect(isProgressiveJPEG(append(withExtras, baseline[2:]...))).IsFalse()
	withExtras = []byte{0xFF, 0xD8, 0xFF, 0xFF, 0xFF}
	Expect(isProgressiveJPEG(append(withExtras, withSOF(baseline, 0xC2)[2:]...))).IsTrue()

	// Truncated, malformed or garbage input is treated as progressive (worst case)
	sof := bytes.Index(baseline, []byte{0xFF, 0xC0})
	for n := 0; n <= sof+1; n++ {
		Expect(isProgressiveJPEG(baseline[:n])).IsTrue()
	}
	Expect(isProgressiveJPEG(nil)).IsTrue()
	Expect(isProgressiveJPEG([]byte("not a jpeg at all"))).IsTrue()
	Expect(isProgressiveJPEG([]byte{0xFF, 0xD8, 0x00, 0xFF, 0xC0})).IsTrue()                   // data where a marker should be
	Expect(isProgressiveJPEG([]byte{0xFF, 0xD8, 0xFF, 0x00, 0xFF, 0xC0})).IsTrue()             // stuffed zero
	Expect(isProgressiveJPEG([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x01, 0xFF, 0xC0})).IsTrue() // invalid segment length
	Expect(isProgressiveJPEG([]byte{0xFF, 0xD8, 0xFF, 0xDA, 0x00, 0x02, 0xFF, 0xC0})).IsTrue() // SOS before SOF
	Expect(isProgressiveJPEG([]byte{0xFF, 0xD8, 0xFF, 0xD9})).IsTrue()                         // EOI before SOF
	Expect(isProgressiveJPEG([]byte{0xFF, 0xD8, 0xFF, 0xE0, 0xFF, 0xFF, 0xFF, 0xC0})).IsTrue() // segment runs past the end

	// SOF beyond the scan limit
	padded := []byte{0xFF, 0xD8}
	for len(padded) < maxJPEGHeaderScan {
		padded = append(padded, 0xFF, 0xE1, 0xFF, 0xFF)
		padded = append(padded, make([]byte, 0xFFFF-2)...)
	}
	padded = append(padded, baseline[2:]...)
	Expect(isProgressiveJPEG(padded)).IsTrue()
}
