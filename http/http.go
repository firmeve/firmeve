package http

import (
	"encoding/json"
	"github.com/firmeve/firmeve/context"
	"github.com/firmeve/firmeve/input/parser"
	"github.com/firmeve/firmeve/kernel/contract"
	"io/ioutil"
	"net/http"
)

type (
	Http struct {
		request        *http.Request
		responseWriter http.ResponseWriter
		message        []byte
		status int
	}
)

func NewHttp(request *http.Request, responseWriter http.ResponseWriter) context.Protocol {
	return &Http{
		request:        request,
		responseWriter: responseWriter,
		//message:make([]byte,0),
	}
}

func (*Http) Name() string {
	return `http`
}

func (h *Http) Read(p []byte) (n int, err error) {
	return h.request.Body.Read(p)
}

//func (h *Http) Seek(offset int64, whence int) (int64, error) {
//	panic("implement me")
//}

func (h *Http) Metadata() map[string][]string {
	return h.request.Header
}

func (h *Http) Message() ([]byte, error) {
	if h.message != nil {
		return h.message, nil
	}

	var err error
	h.message, err = ioutil.ReadAll(h.request.Body)

	return h.message, err
}

func (h *Http) Values() map[string][]string {
	if c.Request.Method == http.MethodGet {
		return parser.NewForm(c.Request.URL.Query())
	}

	switch c.ContentType() {
	case MIMEMultipartPOSTForm:
		c.Request.ParseMultipartForm(32 << 20)
		return parser.NewMultipartForm(c.Request.MultipartForm)
	case MIMEJSON:
		return parser.NewJSON(h.Message())
	case MIMEPOSTForm:
		c.Request.ParseForm()
		fmt.Println(c.Request.Form)
		return parser.NewForm(c.Request.Form)
	}


	if h.request.Method == http.MethodGet {
		return h.request.Form
	}
	h.responseWriter
	return h.request.Form[key]
}

//func (h *Http) Write(v interface{}) (int, error) {
//	panic("implement me")
//}

func (h *Http) Request() *http.Request {
	panic("implement me")
}

func (h *Http) ResponseWriter() http.ResponseWriter {
	panic("implement me")
}

func (h *Http) Write(bytes []byte) (int, error) {
	return h.responseWriter.Write(bytes)
}

func (h *Http) Input(protocol contract.Protocol) map[string][]string {
	panic("implement me")
}

func (h *Http) Output(protocol contract.Protocol,) []byte {
	if protocol.Metadata()[`Accept`][0] == `application/json` {
		json.Marshal()
	}
}

//func (h *Http)  Status(code int) {
//
//}