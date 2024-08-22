package helper

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str []byte) string {
	hash := md5.Sum(str)
	return hex.EncodeToString(hash[:])
}
