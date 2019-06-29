package cache

import (
	"time"
)

var repository *Repository

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
	AddEncode(key string, value interface{}, expire time.Time) error
	ForeverEncode(key string, value interface{}, expire time.Time) error
	PutEncode(key string, value interface{}, expire time.Time) error
	GetDecode(key string) (interface{}, error)
}

type Repository struct {
	repositories map[string]Repository
}

func (this *Repository) NewRepository() *Repository {
	repository = &Repository{}
	return repository
}

type Error struct {
	RepositoryError error
}

func (this *Error) Error() string {
	return this.RepositoryError.Error()
}
