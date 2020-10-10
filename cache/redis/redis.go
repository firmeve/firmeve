package redis

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type repository struct {
	prefix string
	redis  *redis.Client
}

// Initialize a new Repository
func New(client *redis.Client, prefix string) contract.CacheStore {
	return &repository{
		prefix: prefix,
		redis:  client,
	}
}

// Get a warehouse value
func (r *repository) Get(key string) (interface{}, error) {
	return r.redis.Get(r.getPrefix(key)).Result()
}

// Add a new value to the repository
// If the key already exists, the addition will be invalid and will not overwrite the original value.
func (r *repository) Add(key string, value interface{}, expire time.Time) error {
	_, err := r.redis.SetNX(r.getPrefix(key), value, expire.Sub(time.Now())).Result()

	return err
}

// Update a value to the repository
// If the key already exists, it will automatically overwrite
func (r *repository) Put(key string, value interface{}, expire time.Time) error {
	_, err := r.redis.Set(r.getPrefix(key), value, expire.Sub(time.Now())).Result()

	return err
}

// Delete an existing key
func (r *repository) Forget(key string) error {
	_, err := r.redis.Del(r.getPrefix(key)).Result()

	return err
}

// Automatically increments a key by the specified step size
func (r *repository) Increment(key string, steps ...int64) error {
	_, err := r.redis.IncrBy(r.getPrefix(key), r.getStep(steps)).Result()

	return err
}

// Automatically decrements a key by the specified step size
func (r *repository) Decrement(key string, steps ...int64) error {
	_, err := r.redis.DecrBy(r.getPrefix(key), r.getStep(steps)).Result()

	return err
}

// Store a permanent valid key
func (r *repository) Forever(key string, value interface{}) error {
	_, err := r.redis.Set(r.getPrefix(key), value, 0).Result()

	return err
}

// Determine if a key exists
func (r *repository) Has(key string) bool {
	return r.redis.Exists(r.getPrefix(key)).Val() > 0
}

// Clear the current entire database
func (r *repository) Flush() error {
	_, err := r.redis.FlushDB().Result()
	return err
}

// Get a specified step size
func (r *repository) getStep(steps []int64) int64 {
	if len(steps) == 0 {
		return 1
	} else {
		return steps[0]
	}
}

// Get prefix
func (r *repository) getPrefix(key string) string {
	return strings.Join([]string{r.prefix, key}, `:`)
}
