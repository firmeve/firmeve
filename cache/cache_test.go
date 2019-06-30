package cache

import (
	"firmeve/cache/redis"
	config2 "firmeve/config"
	"fmt"
	goRedis "github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestNewRepository(t *testing.T) {

	redisStore := redis.NewRepository(goRedis.NewClient(&goRedis.Options{
		Addr: "192.168.1.107:6379",
		DB:   0,
	}), "redis")

	repository := NewRepository(redisStore)
	fmt.Println(repository)
	//assert.IsType(t, Cache, repository)
}

func TestRepository_Get(t *testing.T) {
	redisRepository := redisRepository()

	key := randString(30)

	value, err := redisRepository.Get(key, "abc")
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

	value, err = redisRepository.Get(key, "abc")
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
	redisStore := redis.NewRepository(goRedis.NewClient(&goRedis.Options{
		Addr: "192.168.1.107:6379",
		DB:   0,
	}), "redis")
	return NewRepository(redisStore)
}

// -------------------- manager ---------------------------

func TestNewManager(t *testing.T) {
	config, err := config2.NewConfig("./testdata/conf")
	if err != nil {
		t.Fail()
	}
	var managern *Manager
	manager := NewManager(config)

	for i := 0; i < 1000; i++ {
		go func(config *config2.Config) {
			managern = NewManager(config)
		}(config)
	}

	if manager != managern{
		t.Fail()
	}
	//fmt.Printf("%#v", manager)

}
