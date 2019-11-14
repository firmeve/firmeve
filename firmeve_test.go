package firmeve

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockProvider struct {
	register bool
	boot     bool
}

func (m *MockProvider) Register() {
	m.register = true
}
func (m *MockProvider) Boot() {
	m.boot = true
}

func TestInstance(t *testing.T) {
	instance1 := Instance()
	instance2 := Instance()
	assert.Equal(t, instance1, instance2)
}

func TestF(t *testing.T) {
	assert.Equal(t, Instance(), F())

	Instance().Bind("testing", func() string {
		return "testing"
	})

	assert.Equal(t, "testing", F(`testing`))
}

func TestFirmeve_Register(t *testing.T) {
	firmeve := New()
	m1 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register("mock", m1)
	assert.Equal(t, true, m1.register)
	assert.Equal(t, false, m1.boot)
	firmeve.Boot()
	assert.Equal(t, true, m1.boot)
	firmeve.Boot()

	m2 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register("mock", m2, WithRegisterForce())
	assert.Equal(t, true, m2.register)
	assert.Equal(t, true, m2.boot)

	m3 := &MockProvider{
		register: false,
		boot:     false,
	}
	firmeve.Register("mock", m3)
	assert.Equal(t, false, m3.register)
	assert.Equal(t, false, m3.boot)
}

func TestFirmeve_GetProvider(t *testing.T) {
	firmeve := New()
	m3 := &MockProvider{
		register: false,
		boot:     true,
	}
	firmeve.Register("mock", m3)
	assert.Implements(t, (*Provider)(nil), firmeve.GetProvider("mock"))
	assert.Equal(t, firmeve.GetProvider("mock").(*MockProvider).boot, true)

	assert.Panics(t, func() {
		firmeve.GetProvider("nothing")
	}, "service provider not exists")
}
