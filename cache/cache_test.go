package cache

import (
	"fmt"
	"github.com/firmeve/firmeve/cache/redis"
	config2 "github.com/firmeve/firmeve/config"
	goRedis "github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
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
	fmt.Println(repository)
	//assert.IsType(t, Cache, repository)
}

func TestRepository_Get(t *testing.T) {
	redisRepository := redisRepository()

	key := randString(30)

	value, err := redisRepository.GetDefault(key, "abc")
	if err != nil {
		t.Fail()
	}

	if value.(string) != "abc" {
		t.Fail()
	}

	err = redisRepository.Put(key, "def", time.Now().Add(time.Second*50))
	if err != nil {
		t.Fail()
	}

	value, err = redisRepository.GetDefault(key, "abc")
	if err != nil {
		t.Fail()
	}

	if value.(string) != "def" {
		t.Fail()
	}
}

func TestRepository_Pull(t *testing.T) {
	redisRepository := redisRepository()

	key := randString(30)

	err := redisRepository.Put(key, "def", time.Now().Add(time.Second*50))
	if err != nil {
		t.Fail()
	}

	value, err := redisRepository.Pull(key, "abc")
	if value.(string) != "def" {
		t.Fail()
	}

	value, err = redisRepository.Pull(key, "abc")
	fmt.Println(value.(string))
	if value.(string) != "abc" {
		t.Fail()
	}

	value, err = redisRepository.Pull(randString(20), "abc")
	if value.(string) != "abc" {
		t.Fail()
	}
}

func TestRepository_Forget(t *testing.T) {
	redisRepository := redisRepository()
	expire := time.Now().Add(time.Second * 10)

	key := randString(50)
	err := redisRepository.Add(key, "a", expire)
	if err != nil {
		t.Fail()
	}

	err = redisRepository.Forget(key)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, false, redisRepository.Has(key))
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

// -------------------- manager ---------------------------

func TestNewManager(t *testing.T) {
	config, err := config2.NewConfig("../testdata/conf")
	if err != nil {
		t.Fail()
	}
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
	config, err := config2.NewConfig("../testdata/conf")
	if err != nil {
		t.Fail()
	}

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
