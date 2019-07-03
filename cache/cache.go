package cache

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/firmeve/firmeve/cache/redis"
	"github.com/firmeve/firmeve/config"
	"github.com/go-ini/ini"
	goRedis "github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

var (
	manager *Manager
	once    sync.Once
	mutex   sync.RWMutex
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

func NewRepository(store Cache) *Repository {
	return &Repository{
		store: store,
	}
}

func (this *Repository) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	if !this.store.Has(key) {
		return defaultValue, nil
	}

	return this.Get(key)
}

func (this *Repository) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	value, err := this.GetDefault(key, defaultValue)
	if err != nil {
		return nil, err
	}

	return value, this.Forget(key)
}

func (this *Repository) Get(key string) (interface{}, error) {
	return this.store.Get(key)
}

func (this *Repository) Add(key string, value interface{}, expire time.Time) error {
	return this.store.Add(key, value, expire)
}

func (this *Repository) Put(key string, value interface{}, expire time.Time) error {
	return this.store.Put(key, value, expire)
}

func (this *Repository) Forever(key string, value interface{}) error {
	return this.store.Forever(key, value)
}

func (this *Repository) Forget(key string) error {
	return this.store.Forget(key)
}

func (this *Repository) Increment(key string, steps ...int64) error {
	return this.store.Increment(key, steps...)
}

func (this *Repository) Decrement(key string, steps ...int64) error {
	return this.store.Decrement(key, steps...)
}

func (this *Repository) Has(key string) bool {
	return this.store.Has(key)
}

func (this *Repository) Flush() error {
	return this.store.Flush()
}

func (this *Repository) Pull(key string) (interface{}, error) {
	value, err := this.Get(key)
	if err != nil {
		return nil, err
	}

	return value, this.Forget(key)
}

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

func (this *Repository) AddEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return this.Add(key, valueBytes, expire)
}

func (this *Repository) ForeverEncode(key string, value interface{}) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return this.Forever(key, valueBytes)
}

func (this *Repository) PutEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return this.Put(key, valueBytes, expire)
}

func gobEncode(value interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buffer).Encode(value)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func gobDecode(data []byte, value interface{}) error {
	buffer := bytes.NewReader(data)
	return gob.NewDecoder(buffer).Decode(value)
}

// -------------------------- manager -----------------------

type Manager struct {
	config       *config.Config
	repositories map[string]Cache
}

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

func (this *Manager) createRedisDriver() Cache {

	addr := []string{
		this.config.GetDefault(`cache.redis.host`, `localhost`).(*ini.Key).String(),
		`:`,
		this.config.GetDefault(`cache.redis.port`, `6379`).(*ini.Key).String(),
	}

	db, _ := this.config.GetDefault(`cache.redis.host`, 0).(*ini.Key).Int()

	return redis.NewRepository(goRedis.NewClient(&goRedis.Options{
		Addr: strings.Join(addr, ``),
		DB:   db,
	}), this.config.GetDefault(`cache.prefix`, `firmeve`).(*ini.Key).String())
}

type Error struct {
	RepositoryError error
}

func (this *Error) Error() string {
	return this.RepositoryError.Error()
}
