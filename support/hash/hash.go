package hash

import (
	sha12 "crypto/sha1"
	"encoding/hex"
)

func Sha1String(target string) string {
	return hex.EncodeToString(Sha1(target))
}

func Sha1(target string) []byte {
	sha1Hash := sha12.New()
	sha1Hash.Write([]byte(target))
	hash := sha1Hash.Sum(nil)
	sha1Hash.Reset()
	return hash
}
