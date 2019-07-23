package cache

import (
	"fmt"
	"github.com/firmeve/firmeve/cache/redis"
	config2 "github.com/firmeve/firmeve/config"
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

	redisStore := redis.NewRepository(goRedis.NewClient(&goRedis.Options{
		Addr: addr + ":6379",
		DB:   0,
	}), "redis")

	repository := NewRepository(redisStore)

	assert.IsType(t, repository, &Repository{})
	assert.Implements(t, new(Cache), repository)
	assert.Implements(t, new(Serialization), repository)

	//var repositoryN *Repository

	//wg.Add(1000)
	//for i := 0; i < 1000; i++ {
	//	go func(i int) {
	//		if i == 800 {
	//			repositoryN = NewRepository(redisStore)
	//		} else {
	//			NewRepository(redisStore)
	//		}
	//		wg.Done()
	//	}(i)
	//}
	//
	////assert.Equal(t, repository, repositoryN)
	//t.Logf("%p",repository)
	//t.Logf("%p\n",repositoryN)
	//t.Log("=================")
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
			key := randString(30) + strconv.Itoa(i)

			err := redisRepository.Put(key, "def", time.Now().Add(time.Second*50))
			assert.Nil(t, err)

			value, err := redisRepository.PullDefault(key, "abc")
			assert.Equal(t, "def", value.(string))

			value, err = redisRepository.PullDefault(key, "abc")
			assert.Equal(t, "abc", value.(string))

			value, err = redisRepository.PullDefault(randString(20), "abc")
			assert.Equal(t, "abc", value.(string))

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

	redisStore := redis.NewRepository(goRedis.NewClient(&goRedis.Options{
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

// -------------------- manager ---------------------------

func TestNewManager(t *testing.T) {
	config := config2.NewConfig("../testdata/conf")
	var managern *Manager
	manager := NewManager(config)

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(config *config2.Config, i int) {
			if i == 999 {
				fmt.Println("zzzzzzzzzzz")
				managern = NewManager(config)
			} else {
				NewManager(config)
			}

			wg.Done()

		}(config, i)
	}

	wg.Wait()

	if manager != managern {
		t.Fail()
	}
	//fmt.Printf("%#v", manager)

}

func TestManager_Driver(t *testing.T) {
	config := config2.NewConfig("../testdata/conf")

	var driver Cache

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			defer wg.Done()
			var err2 error
			manager := NewManager(config)
			if i == 800 {
				driver, err2 = manager.Driver(`redis`)
			} else {
				_, err2 = manager.Driver(`redis`)

			}
			if err2 != nil {
				fmt.Println("error:", err2.Error())
				t.Fail()
			}
		}(i)
	}

	wg.Wait()

	fmt.Println("==============")
	fmt.Printf("%#v", driver)
	fmt.Println("==============")

	cacheInterface := new(redis.Repository)
	assert.IsType(t, cacheInterface, driver)
}

func TestManager_Driver_Error(t *testing.T) {
	config := config2.NewConfig("../testdata/conf")

	manager := NewManager(config)
	_, err2 := manager.Driver(`redis2`)
	assert.NotNil(t, err2)
}
