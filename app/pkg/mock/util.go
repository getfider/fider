package mock

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// RandomString returns a random string of given length
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
