package cache

import (
	"fmt"
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/firmeve/firmeve/cache/redis"
	goRedis "github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestNewRepository(t *testing.T) {
	addr := os.Getenv(`REDIS_HOST`)
	if addr == "" {
		addr = "192.168.1.107"
	}

	redisStore := redis.New(goRedis.NewClient(&goRedis.Options{
		Addr: addr + ":6379",
		DB:   0,
	}), "redis")

	repository := NewRepository(redisStore)

	assert.IsType(t, repository, &Repository{})
	assert.Implements(t, new(CacheInterface), repository)
	assert.Implements(t, new(Serialization), repository)
}

func TestRepository_Get(t *testing.T) {

	redisRepository := redisRepository()
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {

			key := randString(30) + strconv.Itoa(i)

			value, err := redisRepository.GetDefault(key, "abc")
			assert.Nil(t, err)

			assert.Equal(t, "abc", value.(string))

			err = redisRepository.Put(key, "def", time.Now().Add(time.Second*50))
			assert.Nil(t, err)

			value, err = redisRepository.GetDefault(key, "abc")
			assert.Nil(t, err)

			assert.Equal(t, "def", value.(string))

			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestRepository_Pull_Default(t *testing.T) {

	redisRepository := redisRepository()

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			key := t.Name() + randString(30) + strconv.Itoa(i)

			err := redisRepository.Put(key, "def", time.Now().Add(time.Second*150))
			assert.Nil(t, err)

			value1, err := redisRepository.PullDefault(key, "def1")
			assert.Equal(t, "def", value1.(string))

			value2, err := redisRepository.PullDefault(key, "abc")
			assert.Equal(t, "abc", value2.(string))

			value3, err := redisRepository.PullDefault(t.Name() + randString(20), "abcd")
			assert.Equal(t, "abcd", value3.(string))

			wg.Done()
		}(i)
	}
	wg.Wait()

}

func TestRepository_Forget(t *testing.T) {
	redisRepository := redisRepository()
	expire := time.Now().Add(time.Second * 10)

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			key := randString(30) + strconv.Itoa(i)

			err := redisRepository.Add(key, "a", expire)
			assert.Nil(t, err)

			err = redisRepository.Forget(key)
			assert.Nil(t, err)

			assert.Equal(t, false, redisRepository.Has(key))

			wg.Done()
		}(i)
	}
	wg.Wait()
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

func redisRepository() *Repository {
	addr := os.Getenv(`REDIS_HOST`)
	if addr == "" {
		addr = "192.168.1.107"
	}

	redisStore := redis.New(goRedis.NewClient(&goRedis.Options{
		Addr: addr + ":6379",
		DB:   0,
	}), "redis")
	return NewRepository(redisStore)
}

type EncodeTest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestRepository_AddEncode(t *testing.T) {
	test := &EncodeTest{"James", 10}
	redisRepository := redisRepository()
	key := randString(30)
	err := redisRepository.AddEncode(key, test, time.Now().Add(time.Hour))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}
}

func TestRepository_PutEncode(t *testing.T) {
	test := &EncodeTest{"James", 10}
	redisRepository := redisRepository()
	key := randString(30)
	err := redisRepository.PutEncode(key, test, time.Now().Add(time.Hour))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}
}

func TestRepository_ForeverEncode(t *testing.T) {
	test := &EncodeTest{"James", 10}
	redisRepository := redisRepository()
	key := randString(30)
	err := redisRepository.ForeverEncode(key, test)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}
}

func TestRepository_GetDecode(t *testing.T) {

	test := &EncodeTest{"James", 10}
	redisRepository := redisRepository()
	key := randString(30)
	err := redisRepository.AddEncode(key, test, time.Now().Add(time.Hour))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		t.Fail()
	}

	value, err := redisRepository.GetDecode(key, &EncodeTest{})
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, "James", value.(*EncodeTest).Name)
	assert.Equal(t, 10, value.(*EncodeTest).Age)
}

func TestRepository_Increment(t *testing.T) {
	cache := redisRepository()
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
	cache := redisRepository()
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

func TestRepository_Pull(t *testing.T) {
	redisRepository := redisRepository()

	key := randString(30)

	err := redisRepository.Put(key, "def", time.Now().Add(time.Second*50))
	if err != nil {
		t.Fail()
	}

	value, err := redisRepository.Pull(key)
	assert.Nil(t, err)
	assert.Equal(t, "def", value.(string))
}

func TestRepository_Pull_Error(t *testing.T) {
	redisRepository := redisRepository()

	//test := struct {
	//	A string
	//	B int
	//}{
	//	"A",10,
	//}

	_, err := redisRepository.Pull("fdasf!@#$%")
	assert.NotNil(t, err)
}

// -------------------- cache ---------------------------

func TestNewCache(t *testing.T) {
	c1 := Default()
	c2 := Default()
	assert.Equal(t,c1,c2)
}

func TestManager_Driver(t *testing.T) {
	var driver CacheInterface
	driver, err2 := Default().Driver(`redis`)
	if err2 != nil {
		fmt.Println("error:", err2.Error())
		t.Fail()
	}

	cacheInterface := new(redis.Repository)
	assert.IsType(t, cacheInterface, driver)
}

func TestManager_Driver_Error(t *testing.T) {
	_, err2 := Default().Driver(`redis2`)
	assert.NotNil(t, err2)
}


func TestRepository_Flush(t *testing.T) {
	redisRepository := redisRepository()

	_ = redisRepository.Put("abc", "1", time.Now().Add(time.Hour))

	err := redisRepository.Flush()
	if err != nil {
		t.Fail()
	}

	result := redisRepository.Has("abc")
	if result {
		t.Fail()
	}
}


func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.NewFirmeve()
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("cache"))
	assert.Equal(t,true,firmeve.Has(`cache`))
}