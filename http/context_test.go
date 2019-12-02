package http

import (
	"bytes"
	"fmt"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type Data struct {
	Username        string            `form:"username"`
	Password        string            `form:"password"`
	Status          int               `form:"status"`
	Active          bool              `form:"active"`
	ConfirmPassword string            `form:"confirm_password"`
	Multiple        map[string]string `form:"multiple"`
	Hobby           []string          `form:"hobby"`
	Nested          []Nesting         `form:"nested"`
}

type Nesting struct {
	Name  string   `form:"name"`
	Money float64  `form:"money"`
	Hobby []string `form:"hobby"`
}

func TestContext_SetParams_Param(t *testing.T) {
	req := testing2.NewMockRequest(http.MethodGet, "/", "").Request
	c := newContext(
		testing2.TestingModeFirmeve(),
		testing2.NewMockResponseWriter(),
		req,
	)
	c.SetParams(Params{"key": "value"})
	assert.Equal(t, "value", c.Param("key"))
	assert.Equal(t, "", c.Param("nothing"))
}

func TestContext_Form_Query(t *testing.T) {
	firmeve := testing2.TestingModeFirmeve()
	req := testing2.NewMockRequest(http.MethodPost, "/", "").Request
	req.Header.Set(`Content-Type`,MIMEPOSTForm)
	req.Form = make(url.Values, 0)
	req.Form.Set(`username`, `usernameValue`)
	req.Form.Set(`password`, `passwordValue`)
	req.URL.RawQuery = `?query=queryValue1`
	req.ParseMultipartForm(32 << 20)

	c := newContext(firmeve, testing2.NewMockResponseWriter(), req)
	fmt.Println(req.URL.Query().Get(`query1`))
	assert.Equal(t, `usernameValue`, c.Form(`username`))
	assert.Equal(t, `passwordValue`, c.Form(`password`))
	//assert.Equal(t, `queryValue`, c.Query(`query`))
	//assert.Equal(t, `queryValue1`, c.Query(`query1`))
}

func TestContext_FormDecode(t *testing.T) {
	firmeve := testing2.TestingModeFirmeve()
	req := testing2.NewMockRequest(http.MethodPost, "/", "").Request
	req.Header.Set(`Content-Type`,MIMEPOSTForm)
	req.Form = make(url.Values, 0)
	req.Form.Set(`username`, `usernameValue`)
	req.Form.Set(`password`, `passwordValue`)
	req.Form.Set(`confirm_password`, `confirm_passwordValue`)
	req.Form.Set(`status`, `1`)
	req.Form.Set(`active`, `0`)
	req.Form.Set(`hobby[0]`, `a`)
	req.Form.Set(`multiple[k1]`, `v1`)
	req.Form.Set(`multiple[k2]`, `v2`)
	req.Form.Set(`nested[0].name`, `name0`)
	req.Form.Set(`nested[0].money`, `15.0352`)
	req.Form.Set(`nested[0].hobby[0]`, `hobby00`)
	req.Form.Set(`nested[0].hobby[1]`, `hobby01`)
	c := newContext(firmeve, testing2.NewMockResponseWriter(), req)
	data := new(Data)
	c.Bind(data)
	fmt.Printf("%#v\n",data)
	assert.Equal(t, `usernameValue`, data.Username)
	assert.Equal(t, `passwordValue`, data.Password)
	assert.Equal(t, `confirm_passwordValue`, data.ConfirmPassword)
	assert.Equal(t, 1, data.Status)
	assert.Equal(t, false, data.Active)
	assert.Equal(t, map[string]string{`k1`: `v1`, `k2`: `v2`}, data.Multiple)
	assert.Equal(t, []string{`a`}, data.Hobby)
	assert.Equal(t, Nesting{
		Name:  `name0`,
		Money: 15.0352,
		Hobby: []string{`hobby00`, `hobby01`},
	}, data.Nested[0])
}

func TestContext_Entity(t *testing.T) {
	firmeve := testing2.TestingModeFirmeve()
	req := testing2.NewMockRequest(http.MethodGet, "/", "").Request
	c := newContext(firmeve, testing2.NewMockResponseWriter(), req)
	assert.Nil(t, c.Entity("nothing"))

	c.AddEntity("user", map[string]string{
		"username": "username",
		"password": "password",
	})
	u := c.Entity("user").Value.(map[string]string)
	assert.Equal(t,`username`,u[`username`])
	assert.Equal(t,`password`,u[`password`])

	u2 := c.EntityValue("user").(map[string]string)
	assert.Equal(t,`username`,u2[`username`])
	assert.Equal(t,`password`,u2[`password`])
}

func TestContext_ContentType_JSON(t *testing.T) {
	req := testing2.NewMockRequest(http.MethodPost, "/", "").Request
	req.Header.Set(`Content-Type`,MIMEJSON)
	b := bytes.NewBuffer([]byte(`{"str":"中文","numbering":1}`))
	req.Body = ioutil.NopCloser(b)
	firmeve := testing2.TestingModeFirmeve()
	c := newContext(firmeve, testing2.NewMockResponseWriter(), req)
	assert.Equal(t,"中文",c.Input.GetString("str"))
	assert.Equal(t,float64(1),c.Input.GetFloat("numbering"))
}