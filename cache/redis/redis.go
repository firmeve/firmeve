package redis

import (
	"github.com/firmeve/firmeve/cache"
	"github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

var (
	once sync.Once
)

type repository struct {
	prefix string
	redis  *redis.Client
}

// Initialize a new Repository
func New(client *redis.Client, prefix string) cache.Cacheable {
	return &repository{
		prefix: prefix,
		redis:  client,
	}
}

// Get a warehouse value
func (this *repository) Get(key string) (interface{}, error) {
	return this.redis.Get(this.getPrefix(key)).Result()
}

// Add a new value to the repository
// If the key already exists, the addition will be invalid and will not overwrite the original value.
func (this *repository) Add(key string, value interface{}, expire time.Time) error {
	_, err := this.redis.SetNX(this.getPrefix(key), value, expire.Sub(time.Now())).Result()

	return err
}

// Update a value to the repository
// If the key already exists, it will automatically overwrite
func (this *repository) Put(key string, value interface{}, expire time.Time) error {
	_, err := this.redis.Set(this.getPrefix(key), value, expire.Sub(time.Now())).Result()

	return err
}

// Delete an existing key
func (this *repository) Forget(key string) error {
	_, err := this.redis.Del(this.getPrefix(key)).Result()

	return err
}

// Automatically increments a key by the specified step size
func (this *repository) Increment(key string, steps ...int64) error {
	_, err := this.redis.IncrBy(this.getPrefix(key), this.getStep(steps)).Result()

	return err
}

// Automatically decrements a key by the specified step size
func (this *repository) Decrement(key string, steps ...int64) error {
	_, err := this.redis.DecrBy(this.getPrefix(key), this.getStep(steps)).Result()

	return err
}

// Store a permanent valid key
func (this *repository) Forever(key string, value interface{}) error {
	_, err := this.redis.Set(this.getPrefix(key), value, 0).Result()

	return err
}

// Determine if a key exists
func (this *repository) Has(key string) bool {
	return this.redis.Exists(this.getPrefix(key)).Val() > 0
}

// Clear the current entire database
func (this *repository) Flush() error {
	_, err := this.redis.FlushDB().Result()
	return err
}

// Get a specified step size
func (this *repository) getStep(steps []int64) int64 {
	if len(steps) == 0 {
		return 1
	} else {
		return int64(steps[0])
	}
}

// Get prefix
func (this *repository) getPrefix(key string) string {
	return strings.Join([]string{this.prefix, key}, `:`)
}
