package firmeve

import (
	"fmt"
	"github.com/firmeve/firmeve/container"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/testdata/structs"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockProvider struct {
	register bool
	boot     bool
}

var (
	configPath = "./testdata/config"
)

func (m *MockProvider) Name() string {
	return `mock`
}
func (m *MockProvider) Register() {
	m.register = true
}
func (m *MockProvider) Boot() {
	m.boot = true
}

func TestFirmeve(t *testing.T) {
	assert.NotEqual(t, New(kernel.ModeProduction, configPath), New(kernel.ModeProduction, configPath))
	//assert.NotEqual(t, Default(kernel.ModeProduction, configPath), Default(kernel.ModeProduction, configPath))
}

func TestFirmeve_DefaultMode(t *testing.T) {
	f := New(kernel.ModeProduction, configPath)
	assert.Equal(t, f.Mode(), kernel.ModeProduction)
}

func TestFirmeve_SetMode_WithMode(t *testing.T) {
	f := New(kernel.ModeDevelopment, configPath)
	assert.Equal(t, f.Mode(), kernel.ModeDevelopment)
	f.SetMode(kernel.ModeTesting)
	assert.Equal(t, f.Mode(), kernel.ModeTesting)
}

func TestFirmeve_IsProduction(t *testing.T) {
	f := New(kernel.ModeProduction, configPath)
	assert.Equal(t, true, f.IsProduction())
}

func TestFirmeve_IsTesting(t *testing.T) {
	f := New(kernel.ModeTesting, configPath)
	assert.Equal(t, true, f.IsTesting())
}

func TestFirmeve_IsDevelopment(t *testing.T) {
	f := New(kernel.ModeDevelopment, configPath)
	assert.Equal(t, true, f.IsDevelopment())
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
	firmeve := New(kernel.ModeProduction, configPath)
	m1 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register(m1, false)
	assert.Equal(t, true, m1.register)
	assert.Equal(t, true, m1.boot)

	m2 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register(m2, true)
	assert.Equal(t, true, m2.register)
	assert.Equal(t, true, m2.boot)

	m3 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register(m3, true)
	assert.Equal(t, true, m3.register)
	assert.Equal(t, true, m3.boot)
}

func TestFirmeve_Resolve(t *testing.T) {
	c := New(kernel.ModeProduction, configPath)
	c.Bind("dynamic", func() *structs.Nesting {
		return &structs.Nesting{
			NId: 15,
		}
	})
	assert.IsType(t, &structs.Nesting{}, c.Resolve("dynamic"))
}

func TestFirmeve_GetProvider(t *testing.T) {
	//firmeve := New(kernel.ModeProduction, configPath)
	m3 := &MockProvider{
		register: false,
		boot:     true,
	}
	//@todo container会重新new一个对象，原来的值并不保存
	c := container.New()
	fmt.Printf("%#v\n", new(MockProvider))
	fmt.Println(c.Make(m3))
	fmt.Println("##########")
	//firmeve.Register(m3, false)
	//
	//assert.Implements(t, (*kernel.IProvider)(nil), firmeve.GetProvider("mock"))
	//assert.Equal(t, firmeve.GetProvider("mock").(*MockProvider).boot, true)
	//assert.Equal(t, firmeve.GetProvider("mock").(*MockProvider).register, true)
	//fmt.Printf("%#v\n", firmeve.GetProvider("mock").(*MockProvider))
	//assert.Panics(t, func() {
	//	firmeve.GetProvider("nothing")
	//}, "service provider not exists")
}
