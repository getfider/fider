package log

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/pkg/color"
)

var placeholderFinder = regexp.MustCompile("@{.*?}")

// Parse is used to merge props into format and return a text message
func Parse(format string, props dto.Props, colorize bool) string {
	if props == nil || len(props) == 0 {
		return format
	}

	for {
		indexes := placeholderFinder.FindSubmatchIndex([]byte(format))
		if len(indexes) == 0 {
			return format
		}

		ph := format[indexes[0]:indexes[1]]
		phContent := ph[2 : len(ph)-1]
		phSeparatorIdx := strings.Index(phContent, ":")
		value := props[phContent]
		if phSeparatorIdx >= 0 {
			phName := phContent[:phSeparatorIdx]
			phColor := phContent[phSeparatorIdx+1:]
			value = props[phName]
			if colorize {
				value = color.FromName(phColor, value)
			}
		}
		format = fmt.Sprintf("%s%v%s", format[:indexes[0]], value, format[indexes[1]:])
	}
}
