package rand

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRangeInt(t *testing.T) {
	for i := 0; i < 100; i++ {
		assert.Equal(t, true, RangeInt(1, 10) <= 10)
		assert.Equal(t, true, RangeInt(1, 10) >= 1)
		fmt.Println(RangeInt(1, 10))
	}
}
