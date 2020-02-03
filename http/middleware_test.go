package http

import (
	logging "github.com/firmeve/firmeve/logger"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/kataras/iris/core/errors"
	"net/http"
	"testing"
)

func TestRecovery(t *testing.T) {
	firmeve := testing2.TestingModeFirmeve()
	firmeve.Register(new(logging.Provider),true)
	req := testing2.NewMockRequest(http.MethodPost, "/?query=queryValue", "").Request
	req.Header.Set(`Content-Type`,MIMEPOSTForm)
	req.ParseMultipartForm(32 << 20)

	c := newContext(firmeve, testing2.NewMockResponseWriter(), req, func(c *Context) {
		panic(errors.New(`testing error`))
	})

	Recovery(c)
}