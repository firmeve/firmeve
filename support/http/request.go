package http

import (
	"net/http"
	"strings"
)

func ClientIP(req *http.Request) string {
	addr := req.Header.Get(`X-Real-IP`)
	if addr == `` {
		addr = req.Header.Get(`X-Forwarded-For`)
		if addr == `` {
			addr = req.RemoteAddr
		}
	}
	return strings.Split(addr, `:`)[0]
}
