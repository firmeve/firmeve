package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/firmeve/firmeve/cache/redis"
	"github.com/firmeve/firmeve/config"
	"github.com/firmeve/firmeve/container"
	goRedis "github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

type Cache interface {
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
	store Cache
}

type Manager struct {
	config       *config.Config
	repositories map[string]Cache
}

type Error struct {
	RepositoryError error
}

type ServiceProvider struct {
	Firmeve *container.Firmeve `inject:"firmeve"`
}


var (
	manager *Manager
	once    sync.Once
	mutex   sync.Mutex
)

func init()  {
	firmeve := container.GetFirmeve()
	firmeve.Register(`cache`, firmeve.Resolve(new(ServiceProvider)).(*ServiceProvider))
}

func (sp *ServiceProvider) Register() {
	sp.Firmeve.Bind(`cache`, NewManager, container.WithBindShare(true))
}

func (sp *ServiceProvider) Boot() {

}

// Create a new cache repository
func NewRepository(store Cache) *Repository {
	return &Repository{
		store: store,
	}
}

// Get the value of the specified key, or return the default value if the value does not exist
func (this *Repository) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	if !this.store.Has(key) {
		return defaultValue, nil
	}

	return this.Get(key)
}

// Get the value of the specified key, or return the default value if the value does not exist
// Delete this key if it exists
func (this *Repository) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	value, err := this.GetDefault(key, defaultValue)
	if err != nil {
		return nil, err
	}

	return value, this.Forget(key)
}

// Get the value of the specified key
func (this *Repository) Get(key string) (interface{}, error) {
	return this.store.Get(key)
}

// Add a cache value to the specified key
// If the key already exists, it will not be updated
func (this *Repository) Add(key string, value interface{}, expire time.Time) error {
	return this.store.Add(key, value, expire)
}

// Put a cache value to the specified key
func (this *Repository) Put(key string, value interface{}, expire time.Time) error {
	return this.store.Put(key, value, expire)
}

// Permanent storage a cache value to the specified key
func (this *Repository) Forever(key string, value interface{}) error {
	return this.store.Forever(key, value)
}

// Delete the specified key
func (this *Repository) Forget(key string) error {
	return this.store.Forget(key)
}

// Auto increment by the specified step size
func (this *Repository) Increment(key string, steps ...int64) error {
	return this.store.Increment(key, steps...)
}

// Auto decrement by the specified step size
func (this *Repository) Decrement(key string, steps ...int64) error {
	return this.store.Decrement(key, steps...)
}

// Determine if the specified key value exists
func (this *Repository) Has(key string) bool {
	return this.store.Has(key)
}

// Empty the entire cache
func (this *Repository) Flush() error {
	return this.store.Flush()
}

// Get the value of the specified key
// Delete this key if it exists
func (this *Repository) Pull(key string) (interface{}, error) {
	value, err := this.Get(key)
	if err != nil {
		return nil, err
	}

	return value, this.Forget(key)
}

// Get a value that needs to be decoded
// Often used for map, struct
func (this *Repository) GetDecode(key string, to interface{}) (interface{}, error) {
	value, err := this.Get(key)
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
func (this *Repository) AddEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return this.Add(key, valueBytes, expire)
}

// Permanent storage a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (this *Repository) ForeverEncode(key string, value interface{}) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return this.Forever(key, valueBytes)
}

// Put a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (this *Repository) PutEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return this.Put(key, valueBytes, expire)
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

// -------------------------- manager -----------------------

// Create a cache manager
func NewManager(config *config.Config) *Manager {
	if manager != nil {
		return manager
	}

	once.Do(func() {
		manager = &Manager{
			config:       config,
			repositories: make(map[string]Cache),
		}
	})

	return manager
}

// Get the cache driver of the finger
func (this *Manager) Driver(driver string) (Cache, error) {
	var repository Cache
	var err error
	var ok bool

	mutex.Lock()
	defer mutex.Unlock()

	if repository, ok = this.repositories[driver]; ok {
		return repository, err
	}

	switch driver {
	case `redis`:
		repository = this.createRedisDriver()
	default:
		err = &Error{RepositoryError: errors.New("driver not found")}
		return nil, err
	}

	this.repositories[driver] = repository

	return this.repositories[driver], err
}

// Create a redis cache driver
func (this *Manager) createRedisDriver() Cache {
	var (
		host   = this.config.Item("cache").GetString(`redis.host`)
		port   = this.config.Item("cache").GetString(`redis.port`)
		db     = this.config.Item("cache").GetInt(`redis.db`)
		prefix = this.config.Item("cache").GetString(`prefix`)
	)

	addr := []string{
		host,
		`:`,
		port,
	}

	return redis.NewRepository(goRedis.NewClient(&goRedis.Options{
		Addr: strings.Join(addr, ``),
		DB:   db,
	}), prefix)
}

// -------------------------- error -----------------------

// error message
func (this *Error) Error() string {
	return this.RepositoryError.Error()
}
