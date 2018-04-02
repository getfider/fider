package errors

import (
	"runtime"
	"strings"
)

var (
	prefixSize int
	basePath   string
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if file == "?" {
		return
	}
	if ok {
		size := len(file)
		suffix := len("app/pkg/errors/path.go")
		basePath = file[:size-suffix]
		prefixSize = len(basePath)
	}
}

func trimBasePath(filename string) string {
	if strings.HasPrefix(filename, basePath) {
		return filename[prefixSize:]
	}
	return filename
}
