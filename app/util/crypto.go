package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

//MD5 encrypt content with MD5 algorithm
func MD5(content string) string {
	hasher := md5.New()
	hasher.Write([]byte(content))
	return fmt.Sprintf("%v", hex.EncodeToString(hasher.Sum(nil)))
}
