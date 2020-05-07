package hash

import (
	"crypto/md5"
	sha12 "crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

var (
	sha1Hash   = sha12.New()
	md5Hash    = md5.New()
	sha256Hash = sha256.New()
)

func Sha1(target string) []byte {
	sha1Hash.Write([]byte(target))
	hash := sha1Hash.Sum(nil)
	sha1Hash.Reset()
	return hash
}

func Sha1String(target string) string {
	return hex.EncodeToString(Sha1(target))
}

func Md5(target string) []byte {
	md5Hash.Write([]byte(target))
	hash := md5Hash.Sum(nil)
	md5Hash.Reset()
	return hash
}

func Md5String(target string) string {
	return hex.EncodeToString(Md5(target))
}

func Sha256(target string) []byte {
	sha256Hash.Write([]byte(target))
	hash := sha256Hash.Sum(nil)
	sha256Hash.Reset()
	return hash
}

func Sha256String(target string) string {
	return hex.EncodeToString(Sha256(target))
}
