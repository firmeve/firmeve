package jwt

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/testing"
	"github.com/stretchr/testify/assert"
	testing2 "testing"
	"time"
)

func newJwt() contract.Jwt {
	app := testing.ApplicationDefault()
	config := new(Configuration)
	app.Resolve(`config`).(contract.Configuration).Bind(`jwt`, config)
	config.Secret = app.Resolve(`config`).(contract.Configuration).GetString(`framework.key`)
	return New(config, NewMemoryStore())
}

func TestJwt_Create(t *testing2.T) {
	jwt := newJwt()
	v, err := jwt.Create(newAudience("1"))
	assert.Nil(t, err)
	assert.Equal(t, true, v.Lifetime > 0)
	assert.Equal(t, true, len(v.Token) > 0)

	token, err := jwt.Parse(v.Token)
	assert.Nil(t, err)
	assert.Equal(t, true, token.Audience == "1")

	newV, err := jwt.Refresh(v.Token)
	assert.Nil(t, err)
	valid, err := jwt.Valid(v.Token)
	assert.NotNil(t, err)
	assert.Equal(t, false, valid)
	fmt.Println(newV)
}

func TestJwt_Parse_Error(t *testing2.T) {
	jwt := newJwt()
	v, _ := jwt.Create(newAudience("1"))
	_, err := jwt.Parse(v.Token + "!")
	assert.NotNil(t, err)
	//
	b, err := jwt.Valid(v.Token + "!")
	assert.NotNil(t, err)
	assert.Equal(t, false, b)
	//
	err = jwt.Invalidate(v.Token + "!")
	assert.Nil(t, err)
}

func TestJwt_Parse_Valid_Expired(t *testing2.T) {
	app := testing.ApplicationDefault()
	config := new(Configuration)
	app.Resolve(`config`).(contract.Configuration).Bind(`jwt`, config)
	config.Secret = app.Resolve(`config`).(contract.Configuration).GetString(`framework.key`)
	config.Lifetime = 1
	jwt := New(config, NewMemoryStore())

	v, _ := jwt.Create(newAudience("2"))

	time.Sleep(time.Second * 3)
	b, err := jwt.Valid(v.Token)
	assert.Equal(t, false, b)
	assert.Equal(t, ErrorExpired, err)
}

func TestJwt_Parse(t *testing2.T) {
	app := testing.ApplicationDefault()
	config := new(Configuration)
	config.Lifetime = 1
	app.Resolve(`config`).(contract.Configuration).Bind(`jwt`, config)
	config.Secret = `i6KIOXmYKMfKKgQ3Cr3bF2AhGv5hcY5i`

	jwt := New(config, NewMemoryStore())
	v, err := jwt.Parse(`eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOi8veWltZWkudGVzdC9hcGkvdjIvYXV0aC9sb2dpbiIsImlhdCI6MTU5NzMwNzM1NiwiZXhwIjoxNTk4NjAzMzU2LCJuYmYiOjE1OTczMDczNTYsImp0aSI6IlQzb3lKcUtmN0xKM1JQUFAiLCJzdWIiOiIyODEiLCJwcnYiOiI0OGU0NTM4MzFjZWJhNWU1N2E0NzVlNjg2NDljZmRlZTZlOTdkOGQyIn0.JYpKbA7pb06jAJE9B92J-U2INNHgcHTLnOWeU9OH_Z4`)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	fmt.Printf("%v", v)
}
