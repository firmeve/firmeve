package http

import (
	"context"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type (
	Http struct {
		application    contract.Application
		request        *http.Request
		responseWriter http.ResponseWriter
		message        []byte
		params         []httprouter.Param
		route          contract.HttpRoute
		session        contract.Session
	}
)

var (
	defaultMaxSize int64 = 32 << 20
)

func NewHttp(application contract.Application, request *http.Request, responseWriter http.ResponseWriter) contract.HttpProtocol {
	return &Http{
		application:    application,
		request:        request,
		responseWriter: responseWriter,
		params:         make([]httprouter.Param, 0),
		route:          nil,
		session:        nil,
	}
}

func (*Http) Name() string {
	return `http`
}

func (h *Http) Application() contract.Application {
	return h.application
}

func (h *Http) Read(p []byte) (n int, err error) {
	return h.request.Body.Read(p)
}

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

func (h *Http) Request() *http.Request {
	return h.request
}

func (h *Http) SetSession(session contract.Session) {
	h.session = session
}

func (h *Http) Session() contract.Session {
	return h.session
}

func (h *Http) SessionValue(key string) interface{} {
	return h.session.Get(key)
}

func (h *Http) ResponseWriter() http.ResponseWriter {
	return h.responseWriter
}

func (h *Http) Write(bytes []byte) (int, error) {
	return h.responseWriter.Write(bytes)
}

func (h *Http) SetHeader(key, value string) {
	h.request.Header.Set(key, value)
}

func (h *Http) Header(key string) string {
	return h.request.Header.Get(key)
}

func (h *Http) IsContentType(key string) bool {
	return h.ContentType() == key
}

func (h *Http) IsAccept(key string) bool {
	accept := h.Accept()
	for _, v := range accept {
		if v == key {
			return true
		}
	}

	return false
}

func (h *Http) IsMethod(key string) bool {
	return strings.ToUpper(h.request.Method) == strings.ToUpper(key)
}

func (h *Http) SetCookie(cookie *http.Cookie) {
	if cookie.Path == "" {
		cookie.Path = "/"
	}

	cookie.Value = url.QueryEscape(cookie.Value)

	http.SetCookie(h.responseWriter, cookie)
}

func (h *Http) Cookie(name string) (string, error) {
	cookie, err := h.request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (h *Http) Redirect(status int, location string) {
	http.Redirect(h.responseWriter, h.request, location, status)
}

func (h *Http) SetStatus(status int) {
	h.responseWriter.WriteHeader(status)
}

func (h *Http) ContentType() string {
	return strings.Split(h.Header(`Content-Type`), `;`)[0]
}

func (h *Http) Accept() []string {
	return strings.Split(h.Header(`Accept`), `,`)
}

func (h *Http) SetParams(params []httprouter.Param) {
	h.params = params
}

func (h *Http) Params() []httprouter.Param {
	return h.params
}

func (h *Http) Param(key string) httprouter.Param {
	for i := range h.params {
		if h.params[i].Key == key {
			return h.params[i]
		}
	}

	return httprouter.Param{Value: ``, Key: key}
}

func (h *Http) SetRoute(route contract.HttpRoute) {
	h.route = route
}

func (h *Http) Route() contract.HttpRoute {
	return h.route
}

func (h *Http) Values() map[string][]string {
	var params map[string][]string
	if h.IsMethod(http.MethodGet) {
		params = h.request.URL.Query()
	} else {
		switch h.ContentType() {
		case contract.HttpMimeForm:
			params = h.request.Form
		case contract.HttpMimeMultipartForm:
			err := h.request.ParseMultipartForm(defaultMaxSize)
			if err != nil {
				panic(err)
			}
			params = h.request.Form
		}
	}

	// append route params
	routerParams := h.Params()
	if len(routerParams) > 0 {
		for i := range routerParams {
			params[routerParams[i].Key] = []string{routerParams[i].Value}
		}
	}

	return params
}

func (h *Http) Clone() contract.Protocol {
	newHttp := new(Http)
	*newHttp = *h
	//@todo 待验证
	newHttp.request = h.request.Clone(context.Background()) //context.WithCancel(h.request.Context())
	newHttp.responseWriter = h.responseWriter
	return newHttp
}
