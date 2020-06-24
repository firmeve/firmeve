package cache

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	redis2 "github.com/firmeve/firmeve/redis"
	testing2 "github.com/firmeve/firmeve/testing"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Implement(t *testing.T) {
	assert.Implements(t, new(contract.Cache), Default())
}

//Create a cache manager
func Default() contract.Cache {
	app := testing2.ApplicationDefault(new(redis2.Provider), new(Provider))

	return app.Resolve(`cache`).(contract.Cache)
}

func TestRepository_Get(t *testing.T) {

	cache := Default()

	key := t.Name() + randString(30)

	value, err := cache.GetDefault(key, "abc")
	assert.Nil(t, err)

	assert.Equal(t, "abc", value.(string))

	err = cache.Put(key, "def", time.Now().Add(time.Second*50))
	assert.Nil(t, err)

	value, err = cache.GetDefault(key, "abc")
	assert.Nil(t, err)

	assert.Equal(t, "def", value.(string))
}

func TestRepository_Pull_Default(t *testing.T) {

	cache := Default()

	key := t.Name() + randString(30)

	err := cache.Put(key, "def", time.Now().Add(time.Second*150))
	assert.Nil(t, err)

	value1, err := cache.Pull(key)
	assert.Equal(t, "def", value1.(string))

	_, err = cache.Pull(key)
	assert.NotNil(t, err)

	value2, err := cache.PullDefault(key, "abc")
	assert.Equal(t, "abc", value2.(string))
	_, err = cache.PullDefault(key, "abc")
	assert.Nil(t, err)

	value3, err := cache.PullDefault(t.Name()+randString(20), "abcd")
	assert.Equal(t, "abcd", value3.(string))
}

func TestRepository_Forget(t *testing.T) {
	cache := Default()
	expire := time.Now().Add(time.Second * 10)

	key := randString(30)

	err := cache.Add(key, "a", expire)
	assert.Nil(t, err)

	err = cache.Forget(key)
	assert.Nil(t, err)

	assert.Equal(t, false, cache.Has(key))
}

// RandString 生成随机字符串
func randString(len int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

type EncodeTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRepository_AddEncode(t *testing.T) {
	test := &EncodeTest{"James", 10}
	cache := Default()
	key := randString(30)
	err := cache.AddEncode(key, test, time.Now().Add(time.Hour))
	assert.Nil(t, err)

	err1 := cache.AddEncode(``, nil, time.Now().Add(time.Hour))
	assert.NotNil(t, err1)
}

func TestRepository_PutEncode(t *testing.T) {
	test := &EncodeTest{"James", 10}
	cache := Default()
	key := randString(30)
	err := cache.PutEncode(key, test, time.Now().Add(time.Hour))
	assert.Nil(t, err)

	err1 := cache.PutEncode(``, nil, time.Now().Add(time.Hour))
	assert.NotNil(t, err1)
}

func TestRepository_ForeverEncode(t *testing.T) {
	test := &EncodeTest{"James", 10}
	cache := Default()
	key := randString(30)
	err := cache.ForeverEncode(key, test)
	assert.Nil(t, err)

	err1 := cache.ForeverEncode(``, nil)
	assert.NotNil(t, err1)
}

func TestRepository_GetDecode(t *testing.T) {

	test := &EncodeTest{"James", 10}
	cache := Default()
	key := randString(30)
	err := cache.AddEncode(key, test, time.Now().Add(time.Hour))
	assert.Nil(t, err)

	value, err := cache.GetDecode(key, &EncodeTest{})
	assert.Nil(t, err)

	assert.Equal(t, "James", value.(*EncodeTest).Name)
	assert.Equal(t, 10, value.(*EncodeTest).Age)

	_, err = cache.GetDecode(``, &EncodeTest{})
	assert.NotNil(t, err)
}

func TestRepository_Increment(t *testing.T) {
	cache := Default()
	key := randString(50)
	err := cache.Increment(key)
	if err != nil {
		t.Fail()
	}

	value, err := cache.Get(key)
	if err != nil {
		t.Fail()
	}

	num, err := strconv.Atoi(value.(string))
	assert.Equal(t, 1, num)

	err = cache.Increment(key)
	if err != nil {
		t.Fail()
	}

	value, err = cache.Get(key)
	if err != nil {
		t.Fail()
	}

	num, err = strconv.Atoi(value.(string))
	assert.Equal(t, 2, num)

	err = cache.Increment(key, 2)
	if err != nil {
		t.Fail()
	}

	value, err = cache.Get(key)
	if err != nil {
		t.Fail()
	}

	num, err = strconv.Atoi(value.(string))
	assert.Equal(t, 4, num)
}

func TestRepository_Decrement(t *testing.T) {
	cache := Default()
	key := randString(50)

	err := cache.Put(key, 100, time.Now().Add(time.Second*1000))
	if err != nil {
		t.Fail()
	}

	value, err := cache.Get(key)
	if err != nil {
		t.Fail()
	}

	fmt.Println(strconv.Atoi(value.(string)))

	err = cache.Decrement(key)
	if err != nil {
		t.Fail()
	}

	value, err = cache.Get(key)
	if err != nil {
		t.Fail()
	}

	num, err := strconv.Atoi(value.(string))
	assert.Equal(t, 99, num)

	err = cache.Decrement(key)
	if err != nil {
		t.Fail()
	}

	value, err = cache.Get(key)
	if err != nil {
		t.Fail()
	}

	num, err = strconv.Atoi(value.(string))
	assert.Equal(t, 98, num)

	err = cache.Decrement(key, 2)
	if err != nil {
		t.Fail()
	}

	value, err = cache.Get(key)
	if err != nil {
		t.Fail()
	}

	num, err = strconv.Atoi(value.(string))
	assert.Equal(t, 96, num)
}

func TestRepository_Pull(t *testing.T) {
	cache := Default()

	key := randString(30)

	err := cache.Put(key, "def", time.Now().Add(time.Second*50))
	if err != nil {
		t.Fail()
	}

	value, err := cache.PullDefault(key, "")
	assert.Nil(t, err)
	assert.Equal(t, "def", value.(string))
}

// -------------------- cache ---------------------------

func TestManager_Driver_Error(t *testing.T) {
	assert.Panics(t, func() {
		Default().Driver(`redis2`)
	}, "driver not found")
}

func TestCache_Forever(t *testing.T) {
	Default().Forever("a", "b")
	v, _ := Default().Get("a")
	assert.Equal(t, v.(string), "b")
}

func TestCache_Driver(t *testing.T) {
	assert.Implements(t, (*contract.CacheSerializable)(nil), Default().Driver("redis"))
}

func TestCache_Store(t *testing.T) {
	assert.Implements(t, (*contract.CacheStore)(nil), Default().Store())
}

func TestRepository_Flush(t *testing.T) {
	cache := Default()

	_ = cache.Put("abc", "1", time.Now().Add(time.Hour))

	err := cache.Flush()
	if err != nil {
		t.Fail()
	}

	result := cache.Has("abc")
	if result {
		t.Fail()
	}
}
