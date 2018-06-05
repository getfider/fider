// Package letteravatar generates letter-avatars.
package letteravatar

import (
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

// Options are letter-avatar parameters.
type Options struct {
	Font        *truetype.Font
	Palette     []color.Color
	LetterColor color.Color

	// PaletteKey is used to pick the background color from the Palette.
	// Using the same PaletteKey leads to the same background color being picked.
	// If PaletteKey is empty (default) the background color is picked randomly.
	PaletteKey string
}

var defaultLetterColor = color.RGBA{0xf0, 0xf0, 0xf0, 0xf0}

// Extract two letters from given name
func Extract(name string) string {
	name = strings.Trim(name, " ")
	if name == "" {
		return "?"
	}

	slices := splitBySpace(name)
	if len(slices) < 2 {
		slices = splitByCamelCase(name)
	}

	if len(slices) >= 2 {
		firstletter := first(slices[0])
		secondLetter := ""
		for i := len(slices) - 1; i > 0; i-- {
			secondLetter = first(slices[i])
			if secondLetter != "" {
				break
			}
		}
		return strings.ToUpper(firstletter + secondLetter)
	}

	return strings.ToUpper(first(name))
}

func first(s string) string {
	firstRune, _ := utf8.DecodeRuneInString(s)
	if unicode.IsLetter(firstRune) || unicode.IsNumber(firstRune) {
		return string(firstRune)
	}
	return ""
}

func splitBySpace(src string) (entries []string) {
	return strings.FieldsFunc(src, func(r rune) bool {
		return r == '.' || r == ',' || r == ' '
	})
}

func splitByCamelCase(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return
}

// Draw generates a new letter-avatar image of the given size using the given letter
// with the given options. Default parameters are used if a nil *Options is passed.
func Draw(size int, letters string, options *Options) (image.Image, error) {
	font := defaultFont
	if options != nil && options.Font != nil {
		font = options.Font
	}

	palette := defaultPalette
	if options != nil && options.Palette != nil {
		palette = options.Palette
	}

	var letterColor color.Color = defaultLetterColor
	if options != nil && options.LetterColor != nil {
		letterColor = options.LetterColor
	}

	var bgColor color.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	if len(palette) > 0 {
		if options != nil && len(options.PaletteKey) > 0 {
			bgColor = palette[keyindex(len(palette), options.PaletteKey)]
		} else {
			bgColor = palette[randint(len(palette))]
		}
	}

	return drawAvatar(bgColor, letterColor, font, size, letters)
}

func drawAvatar(bgColor, fgColor color.Color, font *truetype.Font, size int, letters string) (image.Image, error) {
	dst := newRGBA(size, size, bgColor)
	fontSize := float64(size) * 0.6
	if len(letters) > 1 {
		fontSize = fontSize * 0.7
	}
	src, err := drawString(bgColor, fgColor, font, fontSize, letters)
	if err != nil {
		return nil, err
	}

	r := src.Bounds().Add(dst.Bounds().Size().Div(2)).Sub(src.Bounds().Size().Div(2))
	draw.Draw(dst, r, src, src.Bounds().Min, draw.Src)

	return dst, nil
}

func drawString(bgColor, fgColor color.Color, font *truetype.Font, fontSize float64, str string) (image.Image, error) {
	c := freetype.NewContext()
	c.SetDPI(72)

	bb := font.Bounds(c.PointToFixed(fontSize))
	w := bb.Max.X.Ceil() - bb.Min.X.Floor()
	h := bb.Max.Y.Ceil() - bb.Min.Y.Floor()

	dst := newRGBA(w, h, bgColor)
	src := image.NewUniform(fgColor)

	c.SetDst(dst)
	c.SetSrc(src)
	c.SetClip(dst.Bounds())
	c.SetFontSize(fontSize)
	c.SetFont(font)

	p, err := c.DrawString(str, fixed.Point26_6{X: 0, Y: bb.Max.Y})
	if err != nil {
		return nil, err
	}

	return dst.SubImage(image.Rect(0, 0, p.X.Ceil(), h)), nil
}

func newRGBA(w, h int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.SetRGBA(x, y, rgba)
		}
	}
	return img
}

func keyindex(n int, key string) int {
	var index int64
	for _, r := range key {
		index = (index + int64(r)) % int64(n)
	}
	return int(index)
}

var (
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rngMu = new(sync.Mutex)
)

func randint(n int) int {
	rngMu.Lock()
	defer rngMu.Unlock()
	return rng.Intn(n)
}
