package cache

import (
	"github.com/firmeve/firmeve/cache/redis"
	"github.com/firmeve/firmeve/cache/repository"
	goRedis "github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	config       *Config
	current      repository.Serializable
	repositories map[string]repository.Serializable
}

type Config struct {
	Prefix       string
	Current      string
	Repositories configRepositoryType
}

type configRepositoryType map[string]interface{}

var (
	mutex sync.Mutex
)

// Create a cache manager
func New(config *Config) *Cache {
	cache := &Cache{
		config:       config,
		repositories: make(map[string]repository.Serializable, 0),
	}
	cache.current = cache.Driver(cache.config.Current)

	return cache
}

// Create a cache manager
func Default() *Cache {
	return New(&Config{
		Prefix:  `firmeve_cache`,
		Current: `redis`,
		Repositories: configRepositoryType{
			`redis`: configRepositoryType{
				`host`: `localhost`,
				`port`: `6379`,
				`db`:   0,
			},
		},
	})
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
		host   = c.config.Repositories[`redis`].(configRepositoryType)[`host`].(string)
		port   = c.config.Repositories[`redis`].(configRepositoryType)[`port`].(string)
		db     = c.config.Repositories[`redis`].(configRepositoryType)[`db`].(int)
		prefix = c.config.Prefix
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
