package redis

import (
	"fmt"
	"testing"
	"time"
)

func TestNewRepository(t *testing.T) {
	cache1 := NewRepository()

	var cacheN *Repository
	for i := 0; i < 10; i++ {
		cacheN = NewRepository()
	}

	if cache1 != cacheN {
		t.Failed()
	}
}

func TestCache_Add(t *testing.T) {
	cache := NewRepository()
	//z := make([]int,2)
	//z := []int{1,2,3}
	z := struct {
		Z string
		X int
	}{
		Z:"z",
		X:10,
	}
	cache.Add("ADD1",z,time.Now().Add(time.Second * 15))
	//now := time.Now()
	//m, _ := time.ParseDuration("-1m")
	//_ := time.Now().Add(time.Second)
	//fmt.Println(now.Add(time.ParseDuration("1s")))
	//time.Sleep(time.Second)
	//cache.Add("x","1",)


}

func TestCache_Has(t *testing.T) {
	cache := NewRepository()
	cache.Add("x","1",time.Time{})
	exists,err := cache.Has("x")
	fmt.Println(exists,err)
}