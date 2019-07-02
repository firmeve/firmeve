package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

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

func client() *redis.Client {
	addr := os.Getenv(`REDIS_HOST`)
	if addr == "" {
		addr = "192.168.1.107"
	}

	return redis.NewClient(&redis.Options{
		Addr: addr + ":6379",
		DB:   0,
	})
}

func TestNewRepository(t *testing.T) {
	cache1 := NewRepository(client(), "redis_")

	var cacheN *Repository
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			if i == 800 {
				cacheN = NewRepository(client(), "redis_")
			} else {
				NewRepository(client(), "redis_")
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	if cache1 != cacheN {
		t.Failed()
	}
	fmt.Println(cache1)

}

func TestRepository_Put_String(t *testing.T) {
	cache := NewRepository(client(), "redis_put_string")

	tests := []struct {
		key   string
		value string
	}{
		{key: randString(5), value: randString(10)},
		{key: randString(5), value: randString(10)},
		{key: randString(5), value: randString(10)},
		{key: randString(5), value: randString(10)},
		{key: randString(5), value: randString(10)},
	}

	for _, v := range tests {
		err := cache.Put(v.key, v.value, time.Now().Add(time.Hour))
		if err != nil {
			t.Fail()
		}

		cur, err := cache.Get(v.key)
		if err != nil {
			t.Fail()
		}

		assert.Equal(t, v.value, cur.(string))
	}
}

func TestRepository_Put_Int(t *testing.T) {

	for i := 0; i < 1000; i++ {
		cache := NewRepository(client(), "redis_put_int")

		tests := []struct {
			key   string
			value int
		}{
			{key: randString(5), value: rand.Int()},
			{key: randString(5), value: rand.Int()},
			{key: randString(5), value: rand.Int()},
			{key: randString(5), value: rand.Int()},
			{key: randString(5), value: rand.Int()},
		}

		for _, v := range tests {
			err := cache.Put(v.key, v.value, time.Now().Add(time.Hour))
			if err != nil {
				t.Fail()
			}

			cur, err := cache.Get(v.key)
			if err != nil {
				t.Fail()
			}

			cur2, err := strconv.Atoi(cur.(string))
			assert.Equal(t, v.value, cur2)
		}
	}

}

func TestRepository_Has(t *testing.T) {
	cache := NewRepository(client(), "redis_put_has")

	key := randString(20)

	if cache.Has(key) {
		t.Fail()
	}

	err := cache.Put(key, "1", time.Now().Add(10))
	if err != nil {
		t.Fail()
	}

	if !cache.Has(key) {
		t.Fail()
	}
}

func TestRepository_Add(t *testing.T) {
	cache := NewRepository(client(), "redis_")
	expire := time.Now().Add(time.Second * 10)

	key := randString(50)

	err := cache.Add(key, "a", expire)
	if err != nil {
		t.Fail()
	}

	err = cache.Add(key, "b", expire)
	//if err == nil {
	//	t.Fail()
	//}

	value, err := cache.Get(key)
	if err != nil {
		t.Fail()
	}

	fmt.Println("=======================")
	fmt.Println(key, value.(string))
	fmt.Println("=======================")

	assert.Equal(t, "a", value.(string))
}

func TestRepository_Flush(t *testing.T) {
	client := client()
	cache := NewRepository(client, "redis_")
	err := cache.Flush()
	if err != nil {
		t.Fail()
	}

	strings, err := client.Keys(`*`).Result()
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 0, len(strings))
}

func TestRepository_Forget(t *testing.T) {
	cache := NewRepository(client(), "redis_")
	expire := time.Now().Add(time.Second * 10)

	key := randString(50)
	err := cache.Add(key, "a", expire)
	if err != nil {
		t.Fail()
	}

	err = cache.Forget(key)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, false, cache.Has(key))
}

func TestRepository_Increment(t *testing.T) {
	cache := NewRepository(client(), "redis_")
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
	cache := NewRepository(client(), "redis_")
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
	fmt.Println("=============")
	fmt.Println(num)
	fmt.Println("=============")
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

func TestRepository_Forever(t *testing.T) {
	cache := NewRepository(client(), "redis_")
	key := randString(50)

	err := cache.Forever(key, 100)
	if err != nil {
		t.Fail()
	}

	value, err := cache.Get(key)
	if err != nil {
		t.Fail()
	}

	num, err := strconv.Atoi(value.(string))
	assert.Equal(t, 100, num)
}
