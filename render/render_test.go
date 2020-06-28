package render

import (
	"github.com/firmeve/firmeve/testing/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func TestErrorRender_Render(t *testing.T) {
//
//	err := kernel.Error(`error message`)
//
//	Error.Render(err)
//}

func TestPlain_Render(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	response := mock.NewMockHttpWrapResponseWriter(ctrl)
	response.EXPECT().WriteHeader(200)
	protocol := mock.NewMockHttpProtocol(ctrl)
	protocol.EXPECT().SetHeader(`Content-Type`, `text/plain`)
	protocol.EXPECT().ResponseWriter().Return(response)
	protocol.EXPECT().Write([]byte(`text`)).Return(len([]byte(`text`)), nil)

	err := Plain.Render(protocol, 200, `text`)

	response.EXPECT().WriteHeader(200)
	protocol.EXPECT().SetHeader(`Content-Type`, `text/plain`)
	protocol.EXPECT().ResponseWriter().Return(response)
	protocol.EXPECT().Write([]byte(`text`)).Return(len([]byte(`text`)), nil)
	err1 := Plain.Render(protocol, 200, []byte(`text`))
	assert.Nil(t, err)
	assert.Nil(t, err1)
}

func TestRenderNull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	response := mock.NewMockHttpWrapResponseWriter(ctrl)
	response.EXPECT().WriteHeader(200)
	protocol := mock.NewMockHttpProtocol(ctrl)
	protocol.EXPECT().ResponseWriter().Return(response)
	protocol.EXPECT().Write(nil)

	err1 := Null.Render(protocol, 200, nil)
	assert.Nil(t, err1)
}

func TestRenderHtml(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	response := mock.NewMockHttpWrapResponseWriter(ctrl)
	response.EXPECT().WriteHeader(200)
	protocol := mock.NewMockHttpProtocol(ctrl)
	protocol.EXPECT().SetHeader(`Content-Type`, `text/html`)
	protocol.EXPECT().ResponseWriter().Return(response)
	protocol.EXPECT().Write([]byte(`<h1>hello</h1>`))

	err1 := Html.Render(protocol, 200, `<h1>hello</h1>`)
	assert.Nil(t, err1)
}

func TestRenderJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	response := mock.NewMockHttpWrapResponseWriter(ctrl)
	response.EXPECT().WriteHeader(200)
	protocol := mock.NewMockHttpProtocol(ctrl)
	protocol.EXPECT().SetHeader(`Content-Type`, `application/json`)
	protocol.EXPECT().ResponseWriter().Return(response)
	protocol.EXPECT().Write([]byte(`{"key":"value"}`))

	err1 := JSON.Render(protocol, 200, map[string]string{
		"key": "value",
	})
	assert.Nil(t, err1)
}
func TestRenderData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	response := mock.NewMockHttpWrapResponseWriter(ctrl)
	response.EXPECT().WriteHeader(200)
	protocol := mock.NewMockHttpProtocol(ctrl)
	protocol.EXPECT().SetHeader(`Content-Type`, `application/json`)
	protocol.EXPECT().ResponseWriter().Return(response)
	protocol.EXPECT().Write([]byte(`{"data":{"key":"value"}}`))

	err1 := Data.Render(protocol, 200, map[string]string{
		"key": "value",
	})
	assert.Nil(t, err1)
}

func TestRenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	response := mock.NewMockHttpWrapResponseWriter(ctrl)
	response.EXPECT().WriteHeader(401)
	protocol := mock.NewMockHttpProtocol(ctrl)
	protocol.EXPECT().SetHeader(`Content-Type`, `application/json`)
	protocol.EXPECT().Accept().Return([]string{`application/json`})
	protocol.EXPECT().ResponseWriter().Return(response)
	protocol.EXPECT().Write([]byte(`{"message":"未认证"}`))

	err1 := Error.Render(protocol, 401, "未认证")
	assert.Nil(t, err1)
}
