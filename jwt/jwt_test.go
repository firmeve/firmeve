package jwt

import (
	"fmt"
	config2 "github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/testing"
	"github.com/stretchr/testify/assert"
	testing2 "testing"
	"time"
)

func newJwt() contract.Jwt {
	app := testing.ApplicationDefault()
	frameConfig := app.Resolve(`config`).(*config2.Config).Item("framework")
	jwtConfig := app.Resolve(`config`).(*config2.Config).Item("jwt")
	return New(frameConfig.GetString("key"), jwtConfig, NewMemoryStore())
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
	frameConfig := app.Resolve(`config`).(*config2.Config).Item("framework")
	jwtConfig := app.Resolve(`config`).(*config2.Config).Item("jwt")
	jwtConfig.Set("lifetime", 1)
	jwt := New(frameConfig.GetString("key"), jwtConfig, NewMemoryStore())

	v, _ := jwt.Create(newAudience("2"))

	time.Sleep(time.Second * 3)
	b, err := jwt.Valid(v.Token)
	assert.Equal(t, false, b)
	assert.Equal(t, ErrorExpired, err)
}
