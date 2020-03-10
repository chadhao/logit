package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

var phoneCodeDigitalTable = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GetRandomCode(max int) string {
	b := make([]byte, max)
	io.ReadAtLeast(rand.Reader, b, max)
	for i := 0; i < max; i++ {
		b[i] = phoneCodeDigitalTable[int(b[i])%len(phoneCodeDigitalTable)]
	}
	return string(b)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
