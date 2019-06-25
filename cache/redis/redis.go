package redis

import (
	"encoding/json"
	"firmeve/cache"
	"fmt"
	"github.com/go-redis/redis"
	"reflect"
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

type Test struct {
	Z string
	X int
}

func (this *Repository) Add(key string, value interface{}, expire time.Time) {
	//z := 1.03
	//fmt.Println(strconv.FormatFloat(z,'E', -1, 32))
	fmt.Println(string(reflect.Int32))
	fmt.Printf("%s\n",reflect.Int32)
	fmt.Printf("%d\n",reflect.Int32)
	//v := reflect.Int16//ValueOf(value)
	//
	//switch v.Kind() {
	//	case reflect.String:
	//	case reflect.Struct:
	//case reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64,reflect.Float32,reflect.Float64:
	//	value := string(value)
	//}
	//
	//s,err := json.Marshal(&Test{Z:"fdsaz",X:1})
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	//
	//fmt.Printf("%s",s)
	v,_ := json.Marshal(value)

	test := &Test{}
	json.Unmarshal(v,test)
	fmt.Println(test)

	//this.redis.Do("SETEX", key, int(expire.Sub(time.Now()).Seconds()),string(v))
	//fmt.Println(z.String())
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