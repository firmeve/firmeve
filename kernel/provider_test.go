package kernel

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseProvider_Bind_Resolve(t *testing.T) {
	app := New()
	app.Bind(`application`, app)
	baseProvider := app.Resolve(new(BaseProvider)).(*BaseProvider)
	baseProvider.Bind(`string`, `string`)
	v := baseProvider.Resolve(`string`).(string)
	assert.Equal(t, `string`, v)
}
