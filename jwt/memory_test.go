package jwt

import (
	"github.com/firmeve/firmeve/support/strings"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var memeory = NewMemoryStore()

func TestMemory_Put_Has_Forget(t *testing.T) {
	id1 := strings.Rand(10)
	id2 := strings.Rand(10)
	id3 := strings.Rand(10)
	id4 := strings.Rand(10)
	id5 := strings.Rand(10)
	err := memeory.Put(id1, newAudience("1"), time.Now().Add(time.Hour))
	assert.Nil(t, err)
	err = memeory.Put(id2, newAudience("1"), time.Now().Add(time.Hour))
	assert.Nil(t, err)
	err = memeory.Put(id3, newAudience("1"), time.Now().Add(time.Hour))
	assert.Nil(t, err)
	err = memeory.Put(id4, newAudience("2"), time.Now().Add(time.Hour))
	assert.Nil(t, err)
	err = memeory.Put(id5, newAudience("3"), time.Now().Add(0))
	assert.Nil(t, err)

	assert.Equal(t, true, memeory.Has(id1))
	assert.Equal(t, true, memeory.Has(id2))
	assert.Equal(t, true, memeory.Has(id3))
	assert.Equal(t, true, memeory.Has(id4))
	assert.Equal(t, false, memeory.Has(id5))
	assert.Equal(t, false, memeory.Has(strings.Rand(10)))

	// forget
	err = memeory.Forget(newAudience("1"))
	assert.Nil(t, err)
	assert.Equal(t, false, memeory.Has(id1))
	assert.Equal(t, false, memeory.Has(id2))
	assert.Equal(t, false, memeory.Has(id3))
	assert.Equal(t, true, memeory.Has(id4))
}
