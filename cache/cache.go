package cache

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/kataras/iris/core/errors"
	"time"
)

type (
	Cache struct {
		current      string
		config       *Configuration
		repositories map[string]contract.CacheSerializable
	}

	Configuration struct {
		Prefix       string
		Default      string `json:"default" yaml:"default"`
		Repositories map[string]struct {
			Connection string `json:"connection" yaml:"connection"`
		} `json:"repositories" yaml:"repositories"`
	}
)

var (
	ErrDriverNotFound = errors.New(`driver not found`)
)

// Create a cache manager
func New(config *Configuration) contract.Cache {
	return &Cache{
		repositories: make(map[string]contract.CacheSerializable, 0),
		config:       config,
		current:      config.Default,
	}
}

// Get the cache driver of the finger
func (c *Cache) Driver(driver string) contract.CacheSerializable {
	var (
		current contract.CacheSerializable
		ok      bool
	)

	if current, ok = c.repositories[driver]; ok {
		return current
	}

	panic(ErrDriverNotFound)
}

// Register driver
func (c *Cache) Register(driver string, store contract.CacheStore) {
	c.repositories[driver] = NewRepository(store)
}

func (c *Cache) Store() contract.CacheStore {
	return c.Driver(c.current).Store()
}

func (c *Cache) Get(key string) (interface{}, error) {
	return c.Driver(c.current).Store().Get(key)
}

func (c *Cache) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	return c.Driver(c.current).GetDefault(key, defaultValue)
}

func (c *Cache) Pull(key string) (interface{}, error) {
	return c.Driver(c.current).Pull(key)
}

func (c *Cache) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	return c.Driver(c.current).PullDefault(key, defaultValue)
}

func (c *Cache) Add(key string, value interface{}, expire time.Time) error {
	return c.Driver(c.current).Store().Add(key, value, expire)
}

func (c *Cache) Put(key string, value interface{}, expire time.Time) error {
	return c.Driver(c.current).Store().Put(key, value, expire)
}

func (c *Cache) Forever(key string, value interface{}) error {
	return c.Driver(c.current).Store().Forever(key, value)
}

func (c *Cache) Forget(key string) error {
	return c.Driver(c.current).Store().Forget(key)
}

func (c *Cache) Increment(key string, steps ...int64) error {
	return c.Driver(c.current).Store().Increment(key, steps...)
}

func (c *Cache) Decrement(key string, steps ...int64) error {
	return c.Driver(c.current).Store().Decrement(key, steps...)
}

func (c *Cache) Has(key string) bool {
	return c.Driver(c.current).Store().Has(key)
}

func (c *Cache) Flush() error {
	return c.Driver(c.current).Store().Flush()
}

func (c *Cache) GetDecode(key string, to interface{}) (interface{}, error) {
	return c.Driver(c.current).GetDecode(key, to)
}

func (c *Cache) AddEncode(key string, value interface{}, expire time.Time) error {
	return c.Driver(c.current).AddEncode(key, value, expire)
}

func (c *Cache) ForeverEncode(key string, value interface{}) error {
	return c.Driver(c.current).ForeverEncode(key, value)
}

func (c *Cache) PutEncode(key string, value interface{}, expire time.Time) error {
	return c.Driver(c.current).PutEncode(key, value, expire)
}
