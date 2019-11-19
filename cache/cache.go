package cache

import (
	"strings"
	"sync"
	"time"

	"github.com/firmeve/firmeve/cache/redis"
	"github.com/firmeve/firmeve/cache/repository"
	"github.com/firmeve/firmeve/config"
	goRedis "github.com/go-redis/redis"
)

type Cache struct {
	config       config.Configurator
	current      repository.Serializable
	repositories map[string]repository.Serializable
}

var (
	mutex sync.Mutex
)

// Create a cache manager
func New(config config.Configurator) *Cache {
	cache := &Cache{
		config:       config,
		repositories: make(map[string]repository.Serializable, 0),
	}
	cache.current = cache.Driver(cache.config.GetString(`default`))

	return cache
}

// Get the cache driver of the finger
func (c *Cache) Driver(driver string) repository.Serializable {
	var current repository.Serializable
	var ok bool

	mutex.Lock()
	defer mutex.Unlock()

	if current, ok = c.repositories[driver]; ok {
		return current
	}

	switch driver {
	case `redis`:
		current = repository.New(c.createRedisDriver())
	default:
		panic(`driver not found`)
	}

	c.repositories[driver] = current

	return c.repositories[driver]
}

// Create a redis cache driver
func (c *Cache) createRedisDriver() repository.Cacheable {
	var (
		host   = c.config.GetString(`repositories.redis.host`)
		port   = c.config.GetString(`repositories.redis.port`)
		db     = c.config.GetInt(`repositories.redis.db`)
		prefix = c.config.GetString(`prefix`)
	)

	addr := []string{host, `:`, port}

	return redis.New(goRedis.NewClient(&goRedis.Options{
		Addr: strings.Join(addr, ``),
		DB:   db,
	}), prefix)
}

func (c *Cache) Store() repository.Cacheable {
	return c.current.Store()
}

func (c *Cache) Get(key string) (interface{}, error) {
	return c.current.Store().Get(key)
}

func (c *Cache) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	return c.current.GetDefault(key, defaultValue)
}

func (c *Cache) Pull(key string) (interface{}, error) {
	return c.current.Pull(key)
}

func (c *Cache) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	return c.current.PullDefault(key, defaultValue)
}

func (c *Cache) Add(key string, value interface{}, expire time.Time) error {
	return c.current.Store().Add(key, value, expire)
}

func (c *Cache) Put(key string, value interface{}, expire time.Time) error {
	return c.current.Store().Put(key, value, expire)
}

func (c *Cache) Forever(key string, value interface{}) error {
	return c.current.Store().Forever(key, value)
}

func (c *Cache) Forget(key string) error {
	return c.current.Store().Forget(key)
}

func (c *Cache) Increment(key string, steps ...int64) error {
	return c.current.Store().Increment(key, steps...)
}

func (c *Cache) Decrement(key string, steps ...int64) error {
	return c.current.Store().Decrement(key, steps...)
}

func (c *Cache) Has(key string) bool {
	return c.current.Store().Has(key)
}

func (c *Cache) Flush() error {
	return c.current.Store().Flush()
}

func (c *Cache) GetDecode(key string, to interface{}) (interface{}, error) {
	return c.current.GetDecode(key, to)
}

func (c *Cache) AddEncode(key string, value interface{}, expire time.Time) error {
	return c.current.AddEncode(key, value, expire)
}

func (c *Cache) ForeverEncode(key string, value interface{}) error {
	return c.current.ForeverEncode(key, value)
}

func (c *Cache) PutEncode(key string, value interface{}, expire time.Time) error {
	return c.current.PutEncode(key, value, expire)
}
