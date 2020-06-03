package hash

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSha1ByteAndString(t *testing.T) {
	assert.Equal(t, hex.EncodeToString(Sha1("hello")), Sha1String("hello"))
}

func TestMd5AndString(t *testing.T) {
	assert.Equal(t, hex.EncodeToString(Md5("hello")), Md5String("hello"))
}

func TestSha256AndString(t *testing.T) {
	assert.Equal(t, hex.EncodeToString(Sha256("hello")), Sha256String("hello"))
}
