package binding

import (
	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/julienschmidt/httprouter"
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

var (
	app contract.Application
)

/*
<input type="text" name="name" value="joeybloggs"/>
  <input type="text" name="age" value="3"/>
  <input type="text" name="gender" value="Male"/>
  <input type="text" name="address[0].name" value="26 Here Blvd."/>
  <input type="text" name="address[0].phone" value="9(999)999-9999"/>
  <input type="text" name="address[1].name" value="26 There Blvd."/>
  <input type="text" name="address[1].phone" value="1(111)111-1111"/>
  <input type="text" name="active" value="true"/>
  <input type="text" name="map_example[key]" value="value"/>
  <input type="text" name="nested_map[key][key]" value="value"/>
  <input type="text" name="nested_array[0][0]" value="value"/>
  <input type="submit"/>
*/

// Address contains address information
type Address struct {
	Name  string `form:"name"`
	Phone string `form:"phone"`
}

// User contains user information
type User struct {
	Name        string                       `form:"name"`
	Age         uint8                        `form:"age"`
	Gender      string                       `form:"gender"`
	Active      bool                         `form:"active"`
	Address     []Address                    `form:"address"`
	MapExample  map[string]string            `form:"map_example"`
	NestedMap   map[string]map[string]string `form:"nested_map"`
	NestedArray [][]string                   `form:"nested_array"`
}

//func TestBindForm(t *testing.T) {
//	buf := bytes.NewBuffer(nil)
//	buf.WriteString("age=3")
//	//buf.WriteString("&name=simon")
//	buf.WriteString("&active=true")
//	buf.WriteString("&gender=1")
//	buf.WriteString("&map_example[key1]=v1")
//	buf.WriteString("&map_example[key2]=v2")
//	req, _ := net_http.NewRequest(net_http.MethodPost, "/?name=simon&age=10", buf)
//	req.Header.Set(`Content-Type`, `application/x-www-form-urlencoded`)
//	//mu := multipart.NewWriter(buf)
//	http := http.NewHttp(app, req, http.NewWrapResponseWriter(&testing2.MockResponseWriter{
//		Bytes:   nil,
//		Headers: nil,
//	}))
//
//	user := new(User)
//	assert.Nil(t, Bind(http, user))
//	fmt.Printf("%v", user)
//}

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
	return false
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

func (h Http) ResponseWriter() contract.HttpWrapResponseWriter {
	return nil
}

func (h Http) ClientIP() string {
	return `0.0.0.0`
}

func TestMain(t *testing.M) {
	app = testing2.ApplicationDefault()
	t.Run()
}

//func TestBindJSON(t *testing.T) {
//	assert.Implements(t, new(contract.HttpProtocol), &Http{})
//	user := new(User)
//	err := Bind(&Http{}, user)
//	fmt.Printf("%v", user)
//	assert.Equal(t, user.Name, "abc")
//	assert.Equal(t, user.Password, "123456")
//	assert.Nil(t, err)
//}
