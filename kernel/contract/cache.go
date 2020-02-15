package contract

import "time"

type (
	CacheStore interface {
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

	CacheSerializable interface {
		Store() CacheStore

		GetDefault(key string, defaultValue interface{}) (interface{}, error)

		Pull(key string) (interface{}, error)

		PullDefault(key string, defaultValue interface{}) (interface{}, error)

		GetDecode(key string, to interface{}) (interface{}, error)

		AddEncode(key string, value interface{}, expire time.Time) error

		ForeverEncode(key string, value interface{}) error

		PutEncode(key string, value interface{}, expire time.Time) error
	}

	Cache interface {
		Driver(driver string) CacheSerializable
		Register(driver string, store CacheStore)
		CacheStore
		CacheSerializable
	}
)
