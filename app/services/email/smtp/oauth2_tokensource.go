package smtp

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"
	"sync"

	"golang.org/x/oauth2"
)

var (
	tokenSourceMu    sync.Mutex
	tokenSourceByKey = map[string]oauth2.TokenSource{}
)

func tokenSourceKey(tokenURL, clientID string, scopes []string) string {
	normalized := make([]string, 0, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope != "" {
			normalized = append(normalized, scope)
		}
	}
	sort.Strings(normalized)

	sum := sha256.Sum256([]byte(tokenURL + "|" + clientID + "|" + strings.Join(normalized, ",")))
	return hex.EncodeToString(sum[:])
}
