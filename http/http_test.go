package http

import (
	"bytes"
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	http2 "net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	httpGet contract.HttpProtocol
	app     contract.Application
)

func TestMain(t *testing.M) {
	app = testing2.ApplicationDefault(new(Provider))
	//req, _ := http2.NewRequest(
	//	http2.MethodGet,
	//	"/get",
	//	os.Stdin,
	//)
	//httpGet = NewHttp(req, &MockResponseWriter{
	//	Bytes:      nil,
	//	StatusCode: 0,
	//	Headers: map[string][]string{
	//		`Accept`: "application/json;text/html",
	//	},
	//})
	t.Run()
}
func pe(err error) {
	if err != nil {
		panic(err)
	}
}

func sampleGetRequest() *http2.Request {
	body := new(bytes.Buffer)
	mu := multipart.NewWriter(body)
	pe(mu.WriteField("p3", "p4"))

	req, err := http2.NewRequest(http2.MethodGet, "/get?p1=v1&p2=v2", body)
	pe(err)

	req.Header.Set("Accept", "application/json,text/html")

	return req
}

func samplePostRequest() *http2.Request {
	body := new(bytes.Buffer)
	mu := multipart.NewWriter(body)
	pe(mu.WriteField("p3", "v3"))

	req, err := http2.NewRequest(http2.MethodPost, "/get?p1=v1&p2=v2", body)
	pe(err)

	req.Header.Set("Accept", "application/json,text/html")
	req.Header.Set("Content-Type", contract.HttpMimeMultipartForm)

	return req
}

func TestHttp_Request_Sample(t *testing.T) {
	req := sampleGetRequest()
	http := NewHttp(app, req, nil)
	accepts := http.Accept()
	assert.Equal(t, []string{"application/json", "text/html"}, accepts)

	v := http.Values()
	assert.Equal(t, map[string][]string{
		"p1": []string{"v1"},
		"p2": []string{"v2"},
	}, v)

	message, err := http.Message()
	pe(err)
	assert.Contains(t, fmt.Sprintf("%s", message), "p3")

	http.SetParams([]httprouter.Param{
		httprouter.Param{
			Key:   "k1",
			Value: "v1",
		},
		httprouter.Param{
			Key:   "k2",
			Value: "v2",
		},
	})

	assert.Equal(t, 2, len(http.Params()))

	assert.Equal(t, "v1", http.Param("k1").Value)
	assert.Equal(t, "v2", http.Param("k2").Value)

	assert.Equal(t, true, http.IsMethod("GET"))
	assert.Equal(t, true, http.IsMethod("get"))

	assert.Equal(t, "http", http.Name())

	assert.Equal(t, true, http.IsAccept("application/json"))
}

func TestHttp_Request_Post_Sample(t *testing.T) {
	http := NewHttp(app, samplePostRequest(), nil)
	assert.Equal(t, contract.HttpMimeMultipartForm, http.ContentType())
}

func TestHttpPostForm(t *testing.T) {
	body := bytes.NewBuffer(nil)
	body.WriteString(`{"name":"simon","phone":"18715155678","password":"123456","sms":"11111"}`)
	req := httptest.NewRequest(http2.MethodPost, "/post", body)
	req.Header.Set(`Content-Type`, `application/json`)
	req.Header.Set(`X-Custom`, `value`)
	req.Header.Set(`Accept`, `application/json`)
	req.AddCookie(&http2.Cookie{
		Name:    "foo",
		Value:   "foo_value",
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	})
	h := NewHttp(app, req, NewWrapResponseWriter(httptest.NewRecorder()))
	assert.Equal(t, h.Header(`x-custom`), `value`)
	assert.Equal(t, h.Header(`X-Custom`), `value`)

	assert.Equal(t, app, h.Application())
	assert.Equal(t, `192.0.2.1`, h.ClientIP())
	assert.Equal(t, req, h.Request())
	assert.Equal(t, true, h.IsMethod(http2.MethodPost))
	assert.Equal(t, `http`, h.Name())
	assert.Equal(t, 0, len(h.Params()))
	h.SetParams([]httprouter.Param{
		{Key: "key", Value: "value"},
	})
	assert.Equal(t, ``, h.Param("a").Value)
	assert.Equal(t, `value`, h.Param("key").Value)
	assert.Equal(t, true, h.IsContentType(contract.HttpMimeJson))
	assert.Equal(t, true, h.IsAccept(contract.HttpMimeJson))
	cookieValue, _ := h.Cookie(`foo`)
	assert.Equal(t, "foo_value", cookieValue)
	message, _ := h.Message()
	assert.Equal(t, true, len(message) > 0)
	//raw := make([]byte, 0)
	//n, _ := h.Read(raw)
	//assert.Equal(t, true, n > 0)
	assert.Equal(t, true, len(h.Metadata()) > 0)

	// response
	h.SetCookie(&http2.Cookie{
		Name:    "foo",
		Value:   "foo_value",
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	})
	h.SetStatus(200)
	h.Write([]byte("string"))
	response := h.ResponseWriter()
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, true, strings.Contains(response.Header().Get(`Set-Cookie`), `foo=foo_value`))
	//assert.Equal(t,true,h.)
	//h.
}
