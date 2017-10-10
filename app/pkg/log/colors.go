package log

import "fmt"

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

// Red return content with red color
func Red(content interface{}) string {
	return fmt.Sprintf("%s%v%s", colorRed, content, resetAll)
}

// Green return content with green color
func Green(content interface{}) string {
	return fmt.Sprintf("%s%v%s", colorGreen, content, resetAll)
}

// Yellow return content with yellow color
func Yellow(content interface{}) string {
	return fmt.Sprintf("%s%v%s", colorYellow, content, resetAll)
}

// Blue return content with blue color
func Blue(content interface{}) string {
	return fmt.Sprintf("%s%v%s", colorBlue, content, resetAll)
}

// Magenta return content with magenta color
func Magenta(content interface{}) string {
	return fmt.Sprintf("%s%v%s", colorMagenta, content, resetAll)
}

// Bold return content with bold text
func Bold(content interface{}) string {
	return fmt.Sprintf("%s%v%s", styleBold, content, resetAll)
}

// Reverse return content with reverse colors
func Reverse(content interface{}) string {
	return fmt.Sprintf("%s%v%s", styleReverse, content, resetAll)
}
