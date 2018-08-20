package rand

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var chars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// String returns a random string of given length
func String(n int) string {
	if n <= 0 {
		return ""
	}

	bytes := make([]byte, n)
	for i := 0; i < n; i++ {
		bytes[i] = chars[rand.Intn(len(chars))]
	}

	return string(bytes)
}
