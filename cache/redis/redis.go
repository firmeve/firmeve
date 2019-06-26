package redis

import (
	"encoding/json"
	"errors"
	"firmeve/cache"
	"github.com/go-redis/redis"
	"reflect"
	"strconv"
	"sync"
	"time"
)

var repository *Repository
var once sync.Once

type Repository struct {
	Prefix string
	redis  *redis.Client
}

func NewRepository(client *redis.Client, prefix string) *Repository {

	if repository != nil {
		return repository
	}

	once.Do(func() {
		repository = &Repository{
			Prefix: prefix,
			redis:  client,
		}
	})

	return repository
}

func (this *Repository) Get(key string) (interface{}, error) {
	return this.redis.Get(key).Result()
}

func (this *Repository) Add(key string, value interface{}, expire time.Time) error {

	if this.Has(key) {
		return &cache.Error{RepositoryError: errors.New("exists error")}
	}

	return this.Put(key, value, expire)
}

func (this *Repository) Put(key string, value interface{}, expire time.Time) error {
	var err error

	expireDuration := expire.Sub(time.Now())
	// 永久
	if expireDuration < 1 {
		expireDuration = 0
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String:
		this.redis.Set(key, value, expireDuration)
	case reflect.Struct, reflect.Map:
		js, err := json.Marshal(value)
		if err != nil {
			return err
		}
		_, err = this.redis.Set(key, string(js), expireDuration).Result()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err = this.redis.Set(key, strconv.Itoa(value.(int)), expireDuration).Result()
	case reflect.Float64:
		_, err = this.redis.Set(key, strconv.FormatFloat(value.(float64), 'G', 5, 64), expireDuration).Result()
	case reflect.Float32:
		_, err = this.redis.Set(key, strconv.FormatFloat(value.(float64), 'G', 5, 32), expireDuration).Result()
	default:
		return &cache.Error{RepositoryError: errors.New("type error")}
	}

	return err
}

func (this *Repository) Forget(key string) error {
	_, err := this.redis.Del(key).Result()
	return err
}

func (this *Repository) Increment(key string, steps ...int64) error {
	_, err := this.redis.IncrBy(key, this.step(steps)).Result()
	return err
}

func (this *Repository) Decrement(key string, steps ...int64) error {
	_, err := this.redis.DecrBy(key, this.step(steps)).Result()
	return err
}

func (this *Repository) Forever(key string, value interface{}) error {
	return this.Put(key, value, time.Now())
}

func (this *Repository) Flush() error {
	_, err := this.redis.FlushDB().Result()
	return err
}

func (this *Repository) Has(key string) bool {
	return this.redis.Exists(key).Val() > 0
}

func (this *Repository) step(steps []int64) int64 {
	if len(steps) == 0 {
		return 1
	} else {
		return int64(steps[0])
	}
}
