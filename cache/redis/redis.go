package redis

import (
	"github.com/go-redis/redis"
	"strings"
	"sync"
	"time"
)

var repository *Repository
var once sync.Once

type Repository struct {
	prefix string
	redis  *redis.Client
}

// 实例化Repository
func NewRepository(client *redis.Client, prefix string) *Repository {

	if repository != nil {
		return repository
	}

	once.Do(func() {
		repository = &Repository{
			prefix: prefix,
			redis:  client,
		}
	})

	return repository
}

func (this *Repository) Get(key string) (interface{}, error) {
	return this.redis.Get(this.getPrefix(key)).Result()
}

func (this *Repository) Add(key string, value interface{}, expire time.Time) error {

	//key = this.prefix(key)

	//if this.Has(key) {
	//	return &cache.Error{RepositoryError: errors.New("exists error")}
	//}

	//var err error
	//
	//if !this.Has(key) {
	//	 err = this.Put(key, value, expire)
	//}

	_, err := this.redis.SetNX(this.getPrefix(key), value, expire.Sub(time.Now())).Result()

	return err
}

func (this *Repository) Put(key string, value interface{}, expire time.Time) error {

	//var err error
	//key = this.prefix(key)

	//expireDuration :=
	//// 永久
	//if expireDuration < 1 {
	//	expireDuration = 0
	//}

	_, err := this.redis.Set(this.getPrefix(key), value, expire.Sub(time.Now())).Result()

	return err
	//
	//v := reflect.ValueOf(value)
	//
	//switch v.Kind() {
	//case reflect.String:
	//	this.redis.Set(key, value, expireDuration)
	//case reflect.Struct, reflect.Map:
	//	js, err := json.Marshal(value)
	//	if err != nil {
	//		return err
	//	}
	//	_, err = this.redis.Set(key, string(js), expireDuration).Result()
	//case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	//	_, err = this.redis.Set(key, strconv.Itoa(value.(int)), expireDuration).Result()
	//case reflect.Float64:
	//	_, err = this.redis.Set(key, strconv.FormatFloat(value.(float64), 'G', 5, 64), expireDuration).Result()
	//case reflect.Float32:
	//	_, err = this.redis.Set(key, strconv.FormatFloat(value.(float64), 'G', 5, 32), expireDuration).Result()
	//default:
	//	return &cache.Error{RepositoryError: errors.New("type error")}
	//}
	//
	//return err
}

func (this *Repository) Forget(key string) error {
	_, err := this.redis.Del(this.getPrefix(key)).Result()
	return err
}

func (this *Repository) Increment(key string, steps ...int64) error {
	_, err := this.redis.IncrBy(this.getPrefix(key), this.getStep(steps)).Result()
	return err
}

func (this *Repository) Decrement(key string, steps ...int64) error {
	_, err := this.redis.DecrBy(this.getPrefix(key), this.getStep(steps)).Result()
	return err
}

func (this *Repository) Forever(key string, value interface{}) error {
	_, err := this.redis.Set(this.getPrefix(key), value, 0).Result()

	return err
}

// has
func (this *Repository) Has(key string) bool {
	return this.redis.Exists(this.getPrefix(key)).Val() > 0
}

// flush all
func (this *Repository) Flush() error {
	_, err := this.redis.FlushDB().Result()
	return err
}

// 数值步长设置
func (this *Repository) getStep(steps []int64) int64 {
	if len(steps) == 0 {
		return 1
	} else {
		return int64(steps[0])
	}
}

// prefix
func (this *Repository) getPrefix(key string) string {
	//return this.prefix + ':' + key
	return strings.Join([]string{this.prefix, key}, `:`)
}
