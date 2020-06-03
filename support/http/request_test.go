package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClientIP2(t *testing.T) {
	body := bytes.NewBuffer(make([]byte, 0))

	req, err := http.NewRequest("GET", "/", body)
	assert.Nil(t, err)
	req.Header.Set(`X-Real-IP`, `127.0.0.2:112111`)
	assert.Equal(t, `127.0.0.2`, ClientIP(req))
	req.Header.Set(`X-Real-IP`, ``)
	req.Header.Set(`X-Forwarded-For`, `127.0.0.1:112111`)
	assert.Equal(t, `127.0.0.1`, ClientIP(req))
	req.Header.Set(`X-Forwarded-For`, ``)
	req.RemoteAddr = `127.0.0.3`
	assert.Equal(t, `127.0.0.3`, ClientIP(req))
}
