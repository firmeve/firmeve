package cache

import "time"

type Error struct {
	RepositoryError error
}

func (this *Error) Error() string {
	return this.RepositoryError.Error()
}

type Cache interface {
	Get(key string) (interface{}, error)

	Add(key string, value interface{}, expire time.Time) error

	Put(key string, value interface{}, expire time.Time) error

	Forget(key string) error

	Increment(key string, steps ...int64) error

	Decrement(key string, steps ...int64) error

	Forever(key string, value interface{}) error

	Has(key string) bool

	Flush() error
}

type Repository struct {
}

func NewRepository() {

}
