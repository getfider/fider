package blob

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/getfider/fider/app"
	"github.com/gosimple/slug"
)

// ErrNotFound is returned when given blob is not found
var ErrNotFound = errors.New("Blob not found")

// ErrInvalidKeyFormat is returned when blob key is in invalid format
var ErrInvalidKeyFormat = errors.New("Blob key is in invalid format")

// SanitizeFileName replaces invalid characters from given filename
func SanitizeFileName(fileName string) string {
	fileName = strings.TrimSpace(fileName)
	ext := filepath.Ext(fileName)
	if ext != "" {
		return slug.Make(fileName[0:len(fileName)-len(ext)]) + ext
	}
	return slug.Make(fileName)
}

// ValidateKey checks if key is is valid format
func ValidateKey(key string) error {
	if len(key) == 0 || len(key) > 512 || strings.Contains(key, " ") {
		return ErrInvalidKeyFormat
	}
	if strings.HasPrefix(key, "/") || strings.HasSuffix(key, "/") {
		return ErrInvalidKeyFormat
	}
	return nil
}

// EnsureAuthorizedPrefix panics if the path is invalid under the given context
func EnsureAuthorizedPrefix(ctx context.Context, path string) {
	// if it's running under the context of a tenant, any prefix is valid
	if ctx.Value(app.TenantCtxKey) != nil {
		return
	}

	// 'tenants' prefix is not valid when running outside a tenant context
	if strings.HasPrefix(path, "tenants") {
		panic(errors.New("Unauthorized access to 'tenants' path."))
	}
}