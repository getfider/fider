package crypto

import (
	"crypto/sha512"

	"fmt"
)

//SHA512 returns the SHA512 hash of a given string
func SHA512(input string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(input)))
}
