package slices

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewWith(t *testing.T) {
	values := []interface{}{1, 2, 2, 3, 4, 4}
	s := NewWith(values)
	s.Add(5)
	assert.Equal(t, s.Len(), 7)
	assert.Equal(t, s.Exists(2), true)
	assert.Equal(t, s.Exists(20), false)
	assert.Equal(t, s.Index(2), 1)
	assert.Equal(t, s.LastIndex(2), 2)
	s.DeleteWithValue(2)
	fmt.Println(s.Values())

	//[1, 2, 2, 3, 4, 4];
}
