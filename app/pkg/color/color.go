package color

import (
	"fmt"
	"strings"
)

var (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[36m"
	styleBold    = "\033[1m"
	styleReverse = "\033[7m"
	resetAll     = "\033[0m"
)

// FromName return content with given color name
func FromName(colorName string, content any) string {
	switch strings.ToUpper(colorName) {
	case "RED":
		return Red(content)
	case "GREEN":
		return Green(content)
	case "YELLOW":
		return Yellow(content)
	case "BLUE":
		return Blue(content)
	case "MAGENTA":
		return Magenta(content)
	case "BOLD":
		return Bold(content)
	case "REVERSE":
		return Reverse(content)
	default:
		return fmt.Sprintf("%v", content)
	}
}

// Red return content with red color
func Red(content any) string {
	return fmt.Sprintf("%s%v%s", colorRed, content, resetAll)
}

// Green return content with green color
func Green(content any) string {
	return fmt.Sprintf("%s%v%s", colorGreen, content, resetAll)
}

// Yellow return content with yellow color
func Yellow(content any) string {
	return fmt.Sprintf("%s%v%s", colorYellow, content, resetAll)
}

// Blue return content with blue color
func Blue(content any) string {
	return fmt.Sprintf("%s%v%s", colorBlue, content, resetAll)
}

// Magenta return content with magenta color
func Magenta(content any) string {
	return fmt.Sprintf("%s%v%s", colorMagenta, content, resetAll)
}

// Bold return content with bold text
func Bold(content any) string {
	return fmt.Sprintf("%s%v%s", styleBold, content, resetAll)
}

// Reverse return content with reverse colors
func Reverse(content any) string {
	return fmt.Sprintf("%s%v%s", styleReverse, content, resetAll)
}
