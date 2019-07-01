package cache

import (
	"errors"
	"firmeve/cache/redis"
	"firmeve/config"
	"github.com/go-ini/ini"
	goRedis "github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

var (
	manager *Manager
	once    sync.Once
	drivers = map[string]string{"redis": "createRedisDriver"}
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
	GetDecode(key string) (interface{}, error)

	AddEncode(key string, value interface{}, expire time.Time) error

	ForeverEncode(key string, value interface{}, expire time.Time) error

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

func (this *Repository) Get(key string, defaultValue interface{}) (interface{}, error) {
	if !this.store.Has(key) {
		return defaultValue, nil
	}

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

func (this *Repository) Pull(key string, defaultValue interface{}) (interface{}, error) {
	value, err := this.Get(key, defaultValue)
	if err != nil {
		return nil, err
	}

	return value, this.Forget(key)
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
			config: config,
		}
	})

	return manager
}

func (this *Manager) driver(driver string) (Cache, error) {
	if repository, ok := this.repositories[driver]; ok {
		return repository, nil
	}

	//var repositoryName string
	//var ok bool
	//if repositoryName, ok = drivers[driver]; !ok {
	//	return nil, &Error{RepositoryError: errors.New("driver not found")}
	//}

	switch driver {
	case `redis`:
		return this.createRedisDriver(), nil
	default:
		return nil, &Error{RepositoryError: errors.New("driver not found")}
	}

	/*value := reflect.ValueOf(this).MethodByName(repositoryName).Call([]reflect.Value{})
	fmt.Println(value[0])
	//return (value[0]).(Cache), nil
	return value[0].(Cache),nil*/
}

func (this *Manager) createRedisDriver() Cache {

	addr := []string{
		this.config.GetDefault(`redis.host`, `localhost`).(*ini.Key).String(),
		`:`,
		this.config.GetDefault(`redis.port`, `6379`).(*ini.Key).String(),
	}

	db, _ := this.config.GetDefault(`redis.host`, 0).(*ini.Key).Int()

	return redis.NewRepository(goRedis.NewClient(&goRedis.Options{
		Addr: strings.Join(addr, ``),
		DB:   db,
	}), this.config.GetDefault(`prefix`, `firmeve`).(*ini.Key).String())
}

type Error struct {
	RepositoryError error
}

func (this *Error) Error() string {
	return this.RepositoryError.Error()
}
