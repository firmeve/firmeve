package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	assert.Equal(t, `a`, Join(`.`, `a`))
	assert.Equal(t, `a.b`, Join(`.`, `a`, `b`))
}

func TestUcWords(t *testing.T) {
	str1 := `HelloWorld`
	assert.Equal(t, str1, UcWords(`hello`, `world`))
}

func TestSnakeCase(t *testing.T) {
	str := `hello_world`
	assert.Equal(t, str, SnakeCase(`helloWorld`))
}

func TestHTMLEntity(t *testing.T) {
	assert.Equal(t, `&amp;`, HTMLEntity(`&`))
}

func TestUcFirst(t *testing.T) {
	str1 := `HelloWorld`
	assert.Equal(t, str1, UcFirst(`helloWorld`))
}

func TestRand(t *testing.T) {
	assert.NotEqual(t,Rand(10),Rand(10))
}