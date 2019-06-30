package cache

import (
	"sync"
	"time"
)

var (
	//repository *Repository
	once sync.Once
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

//func (this *Repository) drive(driver string) *Repository {
//
//}

type Error struct {
	RepositoryError error
}

func (this *Error) Error() string {
	return this.RepositoryError.Error()
}
