package cache

import "time"

type Error struct {
	RepositoryError error
}

func (this *Error) Error() string {
	return this.RepositoryError.Error()
}


type Cache interface {
	Get(key string) error

	Add(key string, value interface{}, expire time.Time)

	Put(key string, value interface{}, expire time.Time)

	Forget(key string)

	Increment(key string)

	Decrement(key string)

	Forever(key string, value interface{}, expire time.Time)

	Has(key string) (bool, error)

	Flush() error
}

type Repository struct {
}

func NewRepository() {

}

func (Repository) Get(key string) error {
	panic("implement me")
}

func (Repository) Add(key string, value interface{}, expire time.Time) {
	panic("implement me")
}

func (Repository) Put(key string, value interface{}, expire time.Time) {
	panic("implement me")
}

func (Repository) Forget(key string) {
	panic("implement me")
}

func (Repository) Increment(key string) {
	panic("implement me")
}

func (Repository) Decrement(key string) {
	panic("implement me")
}

func (Repository) Forever(key string, value interface{}, expire time.Time) {
	panic("implement me")
}

func (Repository) Flush() error {
	panic("implement me")
}
