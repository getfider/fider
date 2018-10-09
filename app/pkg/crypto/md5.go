package crypto

import (
	"crypto/md5"
	"fmt"
)

//MD5 returns the MD5 hash of a given string
func MD5(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}
