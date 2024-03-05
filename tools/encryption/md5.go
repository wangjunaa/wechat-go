package encryption

import (
	"crypto/md5"
	"fmt"
)

func Encode(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func IsEqualAfterEncode(str string, encodeStr string) bool {
	return Encode(str) == encodeStr
}
