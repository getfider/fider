package web

import (
	"mime"
	"strings"
)

// IsJSONContentType returns true if the media type of given Content-Type header
// value is exactly application/json (case-insensitive, parameters such as charset
// are allowed). Substring matching must not be used for this, as values such as
// "text/plain; x=application/json" are CORS-safelisted and can be sent cross-origin.
func IsJSONContentType(contentType string) bool {
	if contentType == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}
	return mediaType == JSONContentType
}

// ParseBearerToken extracts the token from an Authorization header value that uses
// the Bearer scheme (e.g. "Bearer <token>"). The scheme is matched case-insensitively.
// It returns false if the header does not use the Bearer scheme or the token is empty.
func ParseBearerToken(authorization string) (string, bool) {
	authorization = strings.TrimSpace(authorization)
	scheme, token, found := strings.Cut(authorization, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}
	return token, true
}
