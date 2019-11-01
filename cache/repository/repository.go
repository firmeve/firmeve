package repository

import (
	"bytes"
	"encoding/gob"
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

type CacheSerializable interface {
	Cacheable

	GetDefault(key string, defaultValue interface{}) (interface{}, error)

	PullDefault(key string, defaultValue interface{}) (interface{}, error)

	GetDecode(key string, to interface{}) (interface{}, error)

	AddEncode(key string, value interface{}, expire time.Time) error

	ForeverEncode(key string, value interface{}) error

	PutEncode(key string, value interface{}, expire time.Time) error
}

type repository struct {
	store Cacheable
}

// Create a new cache repository
func New(store Cacheable) CacheSerializable {
	return &repository{
		store: store,
	}
}

// Get the value of the specified key, or return the default value if the value does not exist
func (r *repository) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	if !r.store.Has(key) {
		return defaultValue, nil
	}

	return r.Get(key)
}

// Get the value of the specified key, or return the default value if the value does not exist
// Delete this key if it exists
func (r *repository) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	value, err := r.GetDefault(key, defaultValue)
	if err != nil {
		return nil, err
	}

	return value, r.Forget(key)
}

// Get the value of the specified key
func (r *repository) Get(key string) (interface{}, error) {
	return r.store.Get(key)
}

// Add a cache value to the specified key
// If the key already exists, it will not be updated
func (r *repository) Add(key string, value interface{}, expire time.Time) error {
	return r.store.Add(key, value, expire)
}

// Put a cache value to the specified key
func (r *repository) Put(key string, value interface{}, expire time.Time) error {
	return r.store.Put(key, value, expire)
}

// Permanent storage a cache value to the specified key
func (r *repository) Forever(key string, value interface{}) error {
	return r.store.Forever(key, value)
}

// Delete the specified key
func (r *repository) Forget(key string) error {
	return r.store.Forget(key)
}

// Auto increment by the specified step size
func (r *repository) Increment(key string, steps ...int64) error {
	return r.store.Increment(key, steps...)
}

// Auto decrement by the specified step size
func (r *repository) Decrement(key string, steps ...int64) error {
	return r.store.Decrement(key, steps...)
}

// Determine if the specified key value exists
func (r *repository) Has(key string) bool {
	return r.store.Has(key)
}

// Empty the entire cache
func (r *repository) Flush() error {
	return r.store.Flush()
}

// Get the value of the specified key
// Delete this key if it exists
func (r *repository) Pull(key string) (interface{}, error) {
	value, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	return value, r.Forget(key)
}

// Get a value that needs to be decoded
// Often used for map, struct
func (r *repository) GetDecode(key string, to interface{}) (interface{}, error) {
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
func (r *repository) AddEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return r.Add(key, valueBytes, expire)
}

// Permanent storage a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (r *repository) ForeverEncode(key string, value interface{}) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return r.Forever(key, valueBytes)
}

// Put a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (r *repository) PutEncode(key string, value interface{}, expire time.Time) error {
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
