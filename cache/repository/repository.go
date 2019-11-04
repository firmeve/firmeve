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

	Serializable
}

type Serializable interface {
	Store() Cacheable

	GetDefault(key string, defaultValue interface{}) (interface{}, error)

	Pull(key string) (interface{}, error)

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
func New(store Cacheable) Serializable {
	return &repository{
		store: store,
	}
}

// Get the value of the specified key, or return the default value if the value does not exist
func (r *repository) GetDefault(key string, defaultValue interface{}) (interface{}, error) {
	if !r.store.Has(key) {
		return defaultValue, nil
	}

	return r.store.Get(key)
}

// Get the value of the specified key
// Delete this key if it exists
func (r *repository) Pull(key string) (interface{}, error) {
	value, err := r.store.Get(key)
	if err != nil {
		return nil, err
	}

	return value, r.store.Forget(key)
}

// Get the value of the specified key, or return the default value if the value does not exist
// Delete this key if it exists
func (r *repository) PullDefault(key string, defaultValue interface{}) (interface{}, error) {
	value, err := r.GetDefault(key, defaultValue)
	if err != nil {
		return nil, err
	}

	return value, r.store.Forget(key)
}

//
func (r *repository) Store() Cacheable {
	return r.store
}

// Get a value that needs to be decoded
// Often used for map, struct
func (r *repository) GetDecode(key string, to interface{}) (interface{}, error) {
	value, err := r.store.Get(key)
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

	return r.store.Add(key, valueBytes, expire)
}

// Permanent storage a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (r *repository) ForeverEncode(key string, value interface{}) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return r.store.Forever(key, valueBytes)
}

// Put a value that needs to be encode
// Often used for map, struct
// If the key already exists, it will not be updated
func (r *repository) PutEncode(key string, value interface{}, expire time.Time) error {
	valueBytes, err := gobEncode(value)
	if err != nil {
		return err
	}

	return r.store.Put(key, valueBytes, expire)
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
