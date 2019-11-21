package http

import (
	firmeve2 "github.com/firmeve/firmeve"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/magiconair/properties/assert"
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

func TestContext_Map(t *testing.T) {
	firmeve := firmeve2.New()
	req := testing2.NewMockRequest(http.MethodGet, "/", "").Request
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
	data := c.FormDecode(new(Data)).(*Data)
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
