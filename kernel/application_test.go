package kernel

import (
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var app contract.Application

type MockProvider struct {
	BaseProvider
}

func (m MockProvider) Name() string {
	return `mock`
}

func (m *MockProvider) Register() {
	m.Application.Bind(`mock.rand`, func() int {
		seed := rand.NewSource(time.Now().UnixNano())
		rand.New(seed)
		return rand.Int()
	}, container.WithCover(true))

	app.GetProvider(`mock`)
}

func (m MockProvider) Boot() {

}

func TestMain(m *testing.M) {
	app = New()
	app.Bind(`application`, app)
	m.Run()
}

func TestApplication_Boot_Register(t *testing.T) {
	app.Register(new(MockProvider), true)
	app.Register(new(MockProvider), false)
	app.RegisterMultiple([]contract.Provider{new(MockProvider)}, true)
	app.Boot()
}

func TestApplication_Mode(t *testing.T) {
	app.SetMode(contract.ModeDevelopment)
	assert.Equal(t, contract.ModeDevelopment, app.Mode())
	assert.Equal(t, true, app.IsDevelopment())

	app.SetMode(contract.ModeTesting)
	assert.Equal(t, contract.ModeTesting, app.Mode())
	assert.Equal(t, true, app.IsTesting())

	app.SetMode(contract.ModeProduction)
	assert.Equal(t, contract.ModeProduction, app.Mode())
	assert.Equal(t, true, app.IsProduction())
}

func TestApplication_Version(t *testing.T) {
	assert.Equal(t, version, app.Version())
}

func TestApplication_GetProvider(t *testing.T) {
	assert.Panics(t, func() {
		app.GetProvider(`something`)
	})
}

func TestApplication_HasProvider(t *testing.T) {
	assert.Equal(t, false, app.HasProvider(`something`))
}

func TestApplication_Resolve(t *testing.T) {
	app.Register(new(MockProvider), true)
	app.Boot()
	v := app.Resolve(`mock.rand`).(int)
	assert.Equal(t, true, v > 0)
}

func TestApplication_Reset(t *testing.T) {
	app.Register(new(MockProvider), true)
	app.Boot()
	app.Reset()
	assert.Equal(t, false, app.HasProvider(`mock`))
	assert.Equal(t, false, app.Has(`mock.rand`))
}
