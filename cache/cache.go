package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/firmeve/firmeve/cache/redis"
	goRedis "github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

type Cacheable interface {
	Get(key string) (interface{}, error)

	Add(key string, value interface{}, expire time.Time) error

	Put(key string, value interface{}, expire time.Time) error

	Forever(key string, value interface{}) error

	Forget(key string) error

	Increment(key string, steps ...int64) error

	Decrement(key string, steps ...int64) error

	Has(key string) bool

	Flush() error
}

type Serialization interface {
	GetDecode(key string, to interface{}) (interface{}, error)

	AddEncode(key string, value interface{}, expire time.Time) error

	ForeverEncode(key string, value interface{}) error

	PutEncode(key string, value interface{}, expire time.Time) error
}

type Repository struct {
	store Cacheable
}

type ConfigRepositoryType map[string]interface{}

type Config struct {
	Prefix       string
	Current      string
	Repositories ConfigRepositoryType
}

type Cache struct {
	config       *Config
	repositories map[string]Cacheable
}

type Error struct {
	RepositoryError error
}

var (
	instance *Cache
	once  sync.Once
	mutex sync.Mutex
)

// Create a new cache repository
func NewRepository(store Cacheable) *Repository {
	return &Repository{
		store: store,
	}
}

// Get the value of the specified key, or return the default value if the value does not exist
func (r *Repository) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	if !r.store.Has(key) {
		return defaultValue, nil
	}

	return r.Get(key)
}

// Get the value of the specified key, or return the default value if the value does not exist
// Delete this key if it exists
func (r *Repository) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	value, err := r.GetDefault(key, defaultValue)
	if err != nil {
		return nil, err
	}

	return value, r.Forget(key)
}

// Get the value of the specified key
func (r *Repository) Get(key string) (interface{}, error) {
	return r.store.Get(key)
}

// Add a cache value to the specified key
// If the key already exists, it will not be updated
func (r *Repository) Add(key string, value interface{}, expire time.Time) error {
	return r.store.Add(key, value, expire)
}

// Put a cache value to the specified key
func (r *Repository) Put(key string, value interface{}, expire time.Time) error {
	return r.store.Put(key, value, expire)
}

// Permanent storage a cache value to the specified key
func (r *Repository) Forever(key string, value interface{}) error {
	return r.store.Forever(key, value)
}

// Delete the specified key
func (r *Repository) Forget(key string) error {
	return r.store.Forget(key)
}

// Auto increment by the specified step size
func (r *Repository) Increment(key string, steps ...int64) error {
	return r.store.Increment(key, steps...)
}

// Auto decrement by the specified step size
func (r *Repository) Decrement(key string, steps ...int64) error {
	return r.store.Decrement(key, steps...)
}

// Determine if the specified key value exists
func (r *Repository) Has(key string) bool {
	return r.store.Has(key)
}

// Empty the entire cache
func (r *Repository) Flush() error {
	return r.store.Flush()
}

// Get the value of the specified key
// Delete this key if it exists
func (r *Repository) Pull(key string) (interface{}, error) {
	value, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	return value, r.Forget(key)
}

// Get a value that needs to be decoded
// Often used for map, struct
func (r *Repository) GetDecode(key string, to interface{}) (interface{}, error) {
	value, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	err = gobDecode([]byte(value.(string)), to)
	if err != nil {
		return nil, err
	}

	return to, nil
}

// Add a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (r *Repository) AddEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return r.Add(key, valueBytes, expire)
}

// Permanent storage a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (r *Repository) ForeverEncode(key string, value interface{}) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return r.Forever(key, valueBytes)
}

// Put a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (r *Repository) PutEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return r.Put(key, valueBytes, expire)
}

// Use gob mode encode
func gobEncode(value interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buffer).Encode(value)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Use gob mode decode
func gobDecode(data []byte, value interface{}) error {
	buffer := bytes.NewReader(data)
	return gob.NewDecoder(buffer).Decode(value)
}

// -------------------------- Cache -----------------------

// Create a cache manager
func New(config *Config) *Cache {
	return &Cache{
		config:       config,
		repositories: make(map[string]Cacheable, 0),
	}
}

func Instance(params ...interface{}) *Cache {
	if instance != nil {
		return instance
	}

	once.Do(func() {
		instance = New(params[0].(*Config))
	})

	return instance
}

// Create a cache manager
func Default() *Cache {
	return New(&Config{
		Prefix:  `firmeve_cache`,
		Current: `redis`,
		Repositories: ConfigRepositoryType{
			`redis`: ConfigRepositoryType{
				`host`: `localhost`,
				`port`: `6379`,
				`db`:   0,
			},
		},
	})
}

// Get the cache driver of the finger
func (c *Cache) Driver(driver string) (Cacheable, error) {
	var repository Cacheable
	var err error
	var ok bool

	mutex.Lock()
	defer mutex.Unlock()

	if repository, ok = c.repositories[driver]; ok {
		return repository, err
	}

	switch driver {
	case `redis`:
		repository = c.createRedisDriver()
	default:
		err = &Error{RepositoryError: errors.New("driver not found")}
		return nil, err
	}

	c.repositories[driver] = repository

	return c.repositories[driver], err
}

// Create a redis cache driver
func (c *Cache) createRedisDriver() Cacheable {
	var (
		host   = c.config.Repositories[`redis`].(ConfigRepositoryType)[`host`].(string)
		port   = c.config.Repositories[`redis`].(ConfigRepositoryType)[`port`].(string)
		db     = c.config.Repositories[`redis`].(ConfigRepositoryType)[`db`].(int)
		prefix = c.config.Prefix
	)

	addr := []string{host, `:`, port,}

	return redis.New(goRedis.NewClient(&goRedis.Options{
		Addr: strings.Join(addr, ``),
		DB:   db,
	}), prefix)
}

// -------------------------- error -----------------------

// error message
func (err *Error) Error() string {
	return err.RepositoryError.Error()
}
