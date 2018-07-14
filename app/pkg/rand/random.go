package rand

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []byte("123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Int returns a random integer
func Int(min, max int) int {
	return min + rand.Intn(max-min)
}

// String returns a random string of given length
func String(lenght int) string {
	if lenght <= 0 {
		return ""
	}

	bytes := make([]byte, lenght)
	for i := 0; i < lenght; i++ {
		bytes[i] = letterRunes[Int(0, len(letterRunes))]
	}

	return string(bytes)
}
