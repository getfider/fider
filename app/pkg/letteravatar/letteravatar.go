package letteravatar

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Options contains the options for letter avatar generation
type Options struct {
	PaletteKey string
}

// Extract returns the first letter of each word in the given name
func Extract(name string) string {
	if name == "" {
		return "?"
	}

	words := strings.Fields(name)
	if len(words) == 0 {
		return "?"
	}

	// Get first rune from first word
	r, _ := utf8.DecodeRuneInString(words[0])
	if r == utf8.RuneError {
		return "?"
	}

	return string(r)
}

// Draw generates a letter avatar image with the given size and text
func Draw(size int, text string, opts *Options) (image.Image, error) {
	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Generate a background color based on the text
	bg := generateColor(text, opts.PaletteKey)

	// Draw the background
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Try to load Arial font first
	fontBytes, err := os.ReadFile("/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf")
	if err != nil {
		// Fallback to built-in Go font if Arial is not available
		fontBytes = goregular.TTF
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	// Calculate font size (40% of image size)
	fontSize := float64(size) * 0.4

	// Create a new context for drawing text
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)
	c.SetHinting(font.HintingFull)

	// Get text bounds
	bounds, err := c.DrawString(text, freetype.Pt(0, 0))
	if err != nil {
		return nil, err
	}

	// Calculate text position to center it
	textWidth := int(bounds.X.Round())
	textHeight := int(c.PointToFixed(fontSize).Round())
	x := (size - textWidth) / 2
	y := (size + textHeight) / 2

	// Draw the text
	_, err = c.DrawString(text, freetype.Pt(x, y))
	if err != nil {
		return nil, err
	}

	return img, nil
}

// generateColor generates a color based on the text
func generateColor(text, key string) color.Color {
	// Simple hash function
	hash := 0
	for _, r := range text + key {
		hash = int(r) + ((hash << 5) - hash)
	}

	// Convert hash to color
	hue := float64(hash%360) / 360.0
	saturation := 0.5
	value := 0.95

	// Convert HSV to RGB
	h := hue * 6.0
	c := value * saturation
	x := c * (1.0 - math.Abs(math.Mod(h, 2.0)-1.0))
	m := value - c

	var r, g, b float64
	switch int(h) {
	case 0:
		r, g, b = c, x, 0
	case 1:
		r, g, b = x, c, 0
	case 2:
		r, g, b = 0, c, x
	case 3:
		r, g, b = 0, x, c
	case 4:
		r, g, b = x, 0, c
	default:
		r, g, b = c, 0, x
	}

	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
} 