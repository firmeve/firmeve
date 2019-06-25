package redis

import (
	"firmeve/cache"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

type Repository struct {
	redis *redis.Client
}

var repository *Repository
var once sync.Once

func NewRepository() *Repository {

	if repository != nil {
		return repository
	}

	once.Do(func() {
		repository = new(Repository)

		repository.redis = redis.NewClient(&redis.Options{
			Addr:     "192.168.1.107:6379", // use default Addr
			Password: "",                   // no password set
			DB:       0,                    // use default DB
		})

	})

	return repository
}

func (this *Repository) Get(key string) error {
	panic("implement me")
}

func (this *Repository) Add(key string, value interface{}, expire time.Time) {

	z := this.redis.Do("SETEX", key, int(expire.Sub(time.Now()).Seconds()),value)
	fmt.Println(z.String())
	//v := reflect.ValueOf(value)
	//switch v.Kind() {
	//case reflect.String:
	//	this.redis.String().
	//}
	//fmt.Println("==================")
	//fmt.Println(v)
	//fmt.Println("==================")
	//switch v:= value.(type) {
	//case string:
	//	fmt.Println("=========================")
	//	fmt.Println(v)
	//	fmt.Println("=========================")
	//	this.redis.Do("SETEX", key, expire.Sub(expire).Seconds(),expire.Unix())
	//case :
	//
	//}

}

func (this *Repository) Put(key string, value interface{}, expire time.Time) {
	panic("implement me")
}

func (this *Repository) Forget(key string) {
	panic("implement me")
}

func (this *Repository) Increment(key string) {
	panic("implement me")
}

func (this *Repository) Decrement(key string) {
	panic("implement me")
}

func (this *Repository) Forever(key string, value interface{}, expire time.Time) {
	panic("implement me")
}

func (this *Repository) Flush() error {
	panic("implement me")
}

func (this *Repository) Has(key string) (bool, error) {
	exists, err := this.redis.Do("EXISTS", key).Bool()
	if err != nil {
		return exists, &cache.Error{RepositoryError:err}
	}

	return exists, nil
}