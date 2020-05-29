package binding

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	net_http "net/http"
	"testing"
)

const jsonData = `{
	"name":"abc",
	"mobile":"18715155678",
	"password":"123456",
	"sms":"11111"
}`

type User struct {
	Name     string
	Password string
}

type Http struct {
	mock.Mock
}

func (h Http) Application() contract.Application {
	panic("implement me")
}

func (h Http) SetSession(session contract.Session) {
	panic("implement me")
}

func (h Http) Session() contract.Session {
	panic("implement me")
}

func (h Http) SessionValue(key string) interface{} {
	panic("implement me")
}

func (h Http) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (h Http) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (h Http) Name() string {
	return `http`
}

func (h Http) Metadata() map[string][]string {
	panic("implement me")
}

func (h Http) Message() ([]byte, error) {
	return []byte(jsonData), nil
}

func (h Http) Values() map[string][]string {
	panic("implement me")
}

func (h Http) Request() *net_http.Request {
	panic("implement me")
}

func (h Http) ResponseWriter() net_http.ResponseWriter {
	panic("implement me")
}

func (h Http) SetHeader(key, value string) {
	panic("implement me")
}

func (h Http) SetParams(params []httprouter.Param) {
	panic("implement me")
}

func (h Http) Params() []httprouter.Param {
	panic("implement me")
}

func (h Http) Param(key string) httprouter.Param {
	panic("implement me")
}

func (h Http) SetRoute(route contract.HttpRoute) {
	panic("implement me")
}

func (h Http) Route() contract.HttpRoute {
	panic("implement me")
}

func (h Http) Header(key string) string {
	panic("implement me")
}

func (h Http) IsContentType(key string) bool {
	panic("implement me")
}

func (h Http) IsAccept(key string) bool {
	panic("implement me")
}

func (h Http) IsMethod(key string) bool {
	if key == net_http.MethodGet {
		return false
	}
	return true
}

func (h Http) ContentType() string {
	return `application/json`
}

func (h Http) Accept() []string {
	return []string{
		`application/json`,
	}
}

func (h Http) SetStatus(status int) {
	panic("implement me")
}

func (h Http) SetCookie(cookie *net_http.Cookie) {
	panic("implement me")
}

func (h Http) Cookie(name string) (string, error) {
	panic("implement me")
}

func (h Http) Redirect(status int, location string) {
	panic("implement me")
}

func (h Http) Clone() contract.Protocol {
	return h
}

func TestBindJSON(t *testing.T) {
	//	req := testing2.NewMockRequest("post", "/", `{
	//name:abcdef
	//password:123445
	//}`)
	//	req.Request.Header.Set("Content-Type", "application/json")
	//	protocol := http.NewHttp(req.Request, &testing2.MockResponseWriter{})

	user := new(User)
	err := Bind(&Http{}, user)
	assert.Equal(t, user.Name, "abc")
	assert.Equal(t, user.Password, "123456")
	assert.Nil(t, err)
	fmt.Printf("%v", user)
}
