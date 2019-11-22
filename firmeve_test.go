package firmeve

import (
	"github.com/firmeve/firmeve/testdata/structs"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockProvider struct {
	register bool
	boot     bool
}

func (m *MockProvider) Name() string {
	return `mock`
}
func (m *MockProvider) Register() {
	m.register = true
}
func (m *MockProvider) Boot() {
	m.boot = true
}

func TestFirmeve_DefaultMode(t *testing.T) {
	f := New()
	assert.Equal(t, f.mode, ModeProduction)
}

func TestFirmeve_SetMode_WithMode(t *testing.T) {
	f := New(WithMode(ModeDevelopment))
	assert.Equal(t, f.mode, ModeDevelopment)
	f.SetMode(ModeTesting)
	assert.Equal(t, f.Mode(), ModeTesting)
}

func TestFirmeve_IsProduction(t *testing.T) {
	f := New(WithMode(ModeProduction))
	assert.Equal(t,true,f.IsProduction())
}

func TestFirmeve_IsTesting(t *testing.T) {
	f := New(WithMode(ModeTesting))
	assert.Equal(t,true,f.IsTesting())
}

func TestFirmeve_IsDevelopment(t *testing.T) {
	f := New(WithMode(ModeDevelopment))
	assert.Equal(t,true,f.IsDevelopment())
}

//func TestInstance(t *testing.T) {
//	firmeve := New()
//	BindingInstance(firmeve)
//	instance1 := Instance()
//	instance2 := Instance()
//	assert.Equal(t, instance1, instance2)
//	assert.Equal(t, firmeve, instance2)
//	assert.Equal(t, fmt.Sprintf("%p", firmeve), fmt.Sprintf("%p", instance2), fmt.Sprintf("%p", instance1))
//}
//
//func TestF(t *testing.T) {
//	assert.Equal(t, Instance(), F())
//
//	Instance().Bind("testing", func() string {
//		return "testing"
//	})
//
//	assert.Equal(t, "testing", F(`testing`))
//}

func TestFirmeve_Register(t *testing.T) {
	firmeve := New()
	m1 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register(m1)
	assert.Equal(t, true, m1.register)
	assert.Equal(t, false, m1.boot)
	firmeve.Boot()
	assert.Equal(t, true, m1.boot)
	firmeve.Boot()

	m2 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register(m2, WithRegisterForce())
	assert.Equal(t, true, m2.register)
	assert.Equal(t, true, m2.boot)

	m3 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register(m3)
	assert.Equal(t, false, m3.register)
	assert.Equal(t, false, m3.boot)
}

func TestFirmeve_Resolve(t *testing.T) {
	c := New()
	c.Bind("dynamic", func() *structs.Nesting {
		return &structs.Nesting{
			NId: 15,
		}
	})
	assert.IsType(t, &structs.Nesting{}, c.Resolve("dynamic"))
}

func TestFirmeve_GetProvider(t *testing.T) {
	firmeve := New()
	m3 := &MockProvider{
		register: false,
		boot:     true,
	}
	firmeve.Register(m3)
	assert.Implements(t, (*Provider)(nil), firmeve.GetProvider("mock"))
	assert.Equal(t, firmeve.GetProvider("mock").(*MockProvider).boot, true)

	assert.Panics(t, func() {
		firmeve.GetProvider("nothing")
	}, "service provider not exists")
}
