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
	"testing"
)

var (
	httpGet contract.HttpProtocol
)

func TestMain(t *testing.M) {
	testing2.TestingApplication.Register(new(Provider), true)
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
	http := NewHttp(testing2.TestingApplication, req, nil)
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
	http := NewHttp(testing2.TestingApplication, samplePostRequest(), nil)
	assert.Equal(t, contract.HttpMimeMultipartForm, http.ContentType())
}

//func TestHttp_Values(t *testing.T) {
//
//}
