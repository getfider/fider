package md5

import (
	"crypto/md5"
	"fmt"
)

//Hash returns the MD5 hash of a given string
func Hash(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}
