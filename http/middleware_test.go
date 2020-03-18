package http

import (
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
	logging "github.com/firmeve/firmeve/logger"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/kataras/iris/core/errors"
	"net/http"
	"testing"
)

func TestRecovery(t *testing.T) {
	firmeve := testing2.TestingModeFirmeve()
	firmeve.Register(new(logging.Provider), true)
	req := testing2.NewMockRequest(http.MethodPost, "/?query=queryValue", "").Request
	req.Header.Set(`Content-Type`, contract.HttpMimePlain)
	req.ParseMultipartForm(32 << 20)

	c := kernel.NewContext(
		firmeve,
		NewHttp(req, testing2.NewMockResponseWriter()),
		func(c contract.Context) {
			panic(errors.New(`testing error`))
		},
	)

	Recovery(c)
}
