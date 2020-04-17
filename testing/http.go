package testing

import (
	"bytes"
	"net/http"
)

type MockResponseWriter struct {
	//mock.Mock
	Bytes      []byte
	StatusCode int
	Headers    http.Header
}

func (m *MockResponseWriter) Header() http.Header {
	return m.Headers
}

func (m *MockResponseWriter) Write(p []byte) (int, error) {
	m.Bytes = p
	return len(p), nil
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.StatusCode = statusCode
}

func NewMockResponseWriter() http.ResponseWriter {
	return &MockResponseWriter{
		Headers: make(http.Header),
	}
}

type MockRequest struct {
	Request *http.Request
}

func NewMockRequest(method, url, body string) *MockRequest {
	buffer := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		panic(err)
	}
	return &MockRequest{
		Request: req,
	}
}
