package rand

import (
	"crypto/rand"
	"math/big"
)

var chars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// String returns a random string of given length
func String(n int) string {
	if n <= 0 {
		return ""
	}

	bytes := make([]byte, n)
	charsetLen := big.NewInt(int64(len(chars)))
	for i := 0; i < n; i++ {
		c, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic(err)
		}
		bytes[i] = chars[c.Int64()]
	}

	return string(bytes)
}
