package http

import (
	"bytes"
	"errors"
	"github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/kernel"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(500, `error`)
	assert.Equal(t, 500, err.code)
	assert.Equal(t, `error`, err.message)
	assert.Equal(t, `error`, err.Error())
}

func TestNewErrorResponse(t *testing.T) {
	err := NewError(500, `error`)
	firmeve := firmeve.New(kernel.ModeProduction, configPath)
	buffer := bytes.NewBuffer([]byte("content"))
	req, reqError := http.NewRequest(http.MethodGet, "/", buffer)
	if reqError != nil {
		t.Failed()
	}
	writer := &MockResponseWriter{
		Headers: http.Header{
			"Content-Type": []string{"application/json"},
		},
	}
	ctx := newContext(firmeve, writer, req)
	err.Response(ctx)
	assert.Equal(t, "error\n", string(writer.Bytes))

	req2, reqError2 := http.NewRequest(http.MethodGet, "/", buffer)
	if reqError2 != nil {
		t.Failed()
	}
	writer2 := &MockResponseWriter{
		Headers: http.Header{},
	}
	req2.Header.Set(`Accept`, "application/json")
	ctx2 := newContext(firmeve, writer2, req2)
	err.Response(ctx2)
	assert.Equal(t, `{"message":"error"}`, string(writer2.Bytes))

}

func TestNewErrorWithError(t *testing.T) {
	err1 := errors.New(`first_error`)
	err2 := NewErrorWithError(500, `http_error`, err1)
	assert.Equal(t, 500, err2.code)
	assert.Equal(t, `http_error`, err2.message)
	assert.Equal(t, `http_error`, err2.Error())
}

func TestError400(t *testing.T) {
	err := Error400(`error400`)
	assert.Equal(t, 400, err.code)
	assert.Equal(t, `error400`, err.message)
	assert.Equal(t, `error400`, err.Error())
}

func TestError400WithError(t *testing.T) {
	err1 := errors.New(`first_error`)
	err := Error400WithError(`error400`, err1)
	assert.Equal(t, 400, err.code)
	assert.Equal(t, `error400`, err.message)
	assert.Equal(t, `error400`, err.Error())
}

func TestError403(t *testing.T) {
	err := Error403(`error403`)
	assert.Equal(t, 403, err.code)
	assert.Equal(t, `error403`, err.message)
	assert.Equal(t, `error403`, err.Error())
}

func TestError403WithError(t *testing.T) {
	err1 := errors.New(`first_error`)
	err := Error403WithError(`error403`, err1)
	assert.Equal(t, 403, err.code)
	assert.Equal(t, `error403`, err.message)
	assert.Equal(t, `error403`, err.Error())
}

func TestError422(t *testing.T) {
	err := Error422(`error422`)
	assert.Equal(t, 422, err.code)
	assert.Equal(t, `error422`, err.message)
	assert.Equal(t, `error422`, err.Error())
}

func TestError422WithError(t *testing.T) {
	err1 := errors.New(`first_error`)
	err := Error422WithError(`error422`, err1)
	assert.Equal(t, 422, err.code)
	assert.Equal(t, `error422`, err.message)
	assert.Equal(t, `error422`, err.Error())
}

func TestError500(t *testing.T) {
	err := Error500(`error500`)
	assert.Equal(t, 500, err.code)
	assert.Equal(t, `error500`, err.message)
	assert.Equal(t, `error500`, err.Error())
}

func TestError500WithError(t *testing.T) {
	err1 := errors.New(`first_error`)
	err := Error500WithError(`error500`, err1)
	assert.Equal(t, 500, err.code)
	assert.Equal(t, `error500`, err.message)
	assert.Equal(t, `error500`, err.Error())
}
