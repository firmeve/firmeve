package context

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	"github.com/firmeve/firmeve/testing/mock"
	"github.com/firmeve/firmeve/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestNewContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	appMock := mock.NewMockApplication(ctrl)
	protocolMock := mock.NewMockProtocol(ctrl)
	ctx := NewContext(appMock, protocolMock, func(c contract.Context) {
		c.Next()
	})
	assert.Implements(t, new(contract.Context), ctx)
}

type ValidateData struct {
	Name string `form:"name" validate:"required"`
	Age  string `form:"age" validate:"required"`
}

func TestContext_BindValidate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)
	protocolMock.EXPECT().IsMethod(`GET`).Return(false)
	protocolMock.EXPECT().ContentType().Return(contract.HttpMimeForm)
	protocolMock.EXPECT().Values().Return(url.Values{
		`name`: []string{`name`},
	})
	app := testing2.ApplicationDefault(new(validator.Provider))
	ctx := NewContext(app, protocolMock, func(c contract.Context) {
		c.Next()
	})

	validate := new(ValidateData)
	err := ctx.BindValidate(validate)
	assert.NotNil(t, err)
	fmt.Printf("%v", err)
}

func TestContext_BindValidate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)

	gomock.InOrder(
		protocolMock.EXPECT().IsMethod(`GET`).Return(false),
		protocolMock.EXPECT().IsMethod(`GET`).Return(false),
	)
	gomock.InOrder(
		protocolMock.EXPECT().ContentType().Return(contract.HttpMimeForm),
		protocolMock.EXPECT().ContentType().Return(contract.HttpMimeForm),
	)

	gomock.InOrder(
		protocolMock.EXPECT().Values().Return(url.Values{
			`name`: []string{`name`},
			`age`:  []string{`10`},
		}),
		protocolMock.EXPECT().Values().Return(url.Values{
			`s`: []string{`s`},
		}),
	)

	app := testing2.ApplicationDefault(new(validator.Provider))
	ctx := NewContext(app, protocolMock, func(c contract.Context) {
		c.Next()
	})

	validate2 := new(ValidateData)
	err := ctx.BindValidate(validate2)
	assert.Nil(t, err)

	validate3 := new(ValidateData)
	err1 := ctx.BindValidate(validate3)
	assert.NotNil(t, err1)
	//fmt.Printf("%v", err1)
}

func TestContext_Abort(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)
	app := testing2.ApplicationDefault()
	ctx := NewContext(app, protocolMock, func(c contract.Context) {
		c.Next()
	}, func(c contract.Context) {
		c.Next()
	})
	ctx.Abort()
	assert.Equal(t, -1, ctx.Current())
}

func TestContext_AddEntity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)
	app := testing2.ApplicationDefault()
	ctx := NewContext(app, protocolMock, func(c contract.Context) {
		c.Next()
	}, func(c contract.Context) {
		c.Next()
	})
	ctx.AddEntity(`a`, `a`)

	entry := ctx.Entity(`a`)
	assert.Equal(t, `a`, entry.Value.(string))

	notExistEntry := ctx.Entity(`b`)
	assert.Nil(t, notExistEntry)
}

func TestContext_Handlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)
	app := testing2.ApplicationDefault()
	ctx := NewContext(app, protocolMock, func(c contract.Context) {
		c.Next()
	}, func(c contract.Context) {
		c.Next()
	})
	assert.Equal(t, 2, len(ctx.Handlers()))
}

func TestContext_Firmeve_Application(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)
	app := testing2.ApplicationDefault()
	ctx := NewContext(app, protocolMock, func(c contract.Context) {
		c.Next()
	}, func(c contract.Context) {
		c.Next()
	})
	assert.Implements(t, new(contract.Application), ctx.Application())
	assert.Implements(t, new(contract.Application), ctx.Firmeve())
}

func TestContext_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)
	gomock.InOrder(
		protocolMock.EXPECT().Values().Return(url.Values{
			`a`: []string{`a`},
		}),
		protocolMock.EXPECT().Values().Return(url.Values{
			`a`: []string{`a`},
		}),
		protocolMock.EXPECT().Values().Return(url.Values{
			`a`: []string{`a`},
			`s`: []string{`s`, `s2`},
		}),
	)

	app := testing2.ApplicationDefault()
	ctx := NewContext(app, protocolMock, nil)
	v := ctx.Get(`a`)
	assert.Equal(t, `a`, v.(string))
	v2 := ctx.Get(`b`)
	assert.Nil(t, v2)
	v3 := ctx.Get(`s`)
	assert.Equal(t, []string{`s`, `s2`}, v3)
}

func TestContext_Protocol(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	protocolMock := mock.NewMockHttpProtocol(ctrl)
	app := testing2.ApplicationDefault()
	ctx := NewContext(app, protocolMock, nil)
	assert.Implements(t, new(contract.HttpProtocol), ctx.Protocol())
}

func TestContext_Bind(t *testing.T) {

}
