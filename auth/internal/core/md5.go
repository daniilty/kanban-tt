package core

import (
	"crypto/md5"
	"fmt"
)

func getMD5Sum(s string) string {
	data := []byte(s)

	return fmt.Sprintf("%x", md5.Sum(data))
}
