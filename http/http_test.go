package http

import (
	"bytes"
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	http2 "net/http"
	"testing"
)

func TestHttp_Values(t *testing.T) {
	req, err := http2.NewRequest(http2.MethodPost, "/mock?p1=1&p2=2", bytes.NewBuffer([]byte(`p2=new2&array[]=1&array[]=2&object[key]=value`)))
	if err != nil {
		panic(err)
	}
	req.Header.Set(`Content-Type`, contract.HttpMimeForm)
	//m := multipart.NewWriter(bytes.NewBuffer([]byte(``)))
	//m.WriteField()

	h := NewHttp(req, &MockResponseWriter{})

	values := h.Values()
	assert.NotNil(t, values)
	assert.Equal(t, []string{"new2", "2"}, values[`p2`])
	fmt.Printf("%#v", h.Values())
}

func basePostFormRequest() *http2.Request {
	body := bytes.NewBuffer([]byte(``))
	m := multipart.NewWriter(body)
	defer m.Close()

	m.WriteField("t1", "v1")
	m.WriteField("t2", "v2")
	m.WriteField("t3[0]", "v30")
	m.WriteField("t3[1]", "v31")
	m.WriteField("t4[key1]", "t4value1")
	m.WriteField("t4[key2]", "t4value2")

	req, err := http2.NewRequest(http2.MethodPost, "/", body)
	if err == nil {
		panic(err)
	}

	req.Header.Set("Content-Type", contract.HttpMimeForm)
	return req
}

// -------------------- context ------------------------

func TestContextBind(t *testing.T) {
	//context := kernel.NewContext()
}

//func TestHttp_Values(t *testing.T) {
//	req := testing2.NewMockRequest("POST", "/mock?p1=1&p2=2", `{
//	"p2":"abc",
//	"name":"hello",
//	"password":"password"
//}`)
//	req.Request.Header.Set("Content-Type", `application/json`)
//	//req.Request.ParseForm()
//
//	//buff := bytes.NewBuffer([]byte())
//	params := map[string]interface{}{
//		"p2":       "abc",
//		"name":     "hello",
//		"password": "password",
//	}
//	jsonBytes, err := json.Marshal(params)
//	if err != nil {
//		panic(err)
//	}
//
//	req2, _ := http2.NewRequest(http2.MethodPost, "/?p1=1&p2=2", bytes.NewReader(jsonBytes))
//	req2.Header.Set("Content-Type", "application/json;charset=utf-8")
//
//	req2.ParseForm()
//	fmt.Printf("%s", req2.PostForm)
//	http := NewHttp(req2, &MockResponseWriter{})
//	fmt.Printf("%#v", http.Values())
//}
