package web

import (
	"errors"
	"mime"
	"strings"
)

// IsJSONContentType returns true if the media type of given Content-Type header
// value is application/json (case-insensitive, parameters such as charset are
// allowed). Substring matching must not be used for this, as values such as
// "text/plain; x=application/json" are CORS-safelisted and can be sent cross-origin.
//
// Malformed or duplicate parameters (e.g. "application/json; charset") are
// tolerated for backwards compatibility, as long as the media type itself is
// application/json. This is still CSRF-safe: a CORS simple request can never
// have an application/json essence.
func IsJSONContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil && !errors.Is(err, mime.ErrInvalidMediaParameter) {
		return false
	}
	return mediaType == JSONContentType
}

// ParseBearerToken extracts the token from an Authorization header value that uses
// the Bearer scheme. The header must start with "Bearer" (case-insensitive); the
// remainder, trimmed of any whitespace, is the token and must not be empty.
//
// This is deliberately lenient (e.g. "Bearer\t<key>" and "Bearer<key>" are
// accepted) to stay compatible with how API keys were historically parsed.
// Other schemes, including "Basic Bearer <key>", are rejected.
func ParseBearerToken(authorization string) (string, bool) {
	const scheme = "Bearer"
	authorization = strings.TrimSpace(authorization)
	if len(authorization) < len(scheme) || !strings.EqualFold(authorization[:len(scheme)], scheme) {
		return "", false
	}
	token := strings.TrimSpace(authorization[len(scheme):])
	if token == "" {
		return "", false
	}
	return token, true
}
