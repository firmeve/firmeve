package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var (
	directory = "../testdata/"
	wg        sync.WaitGroup
)

func TestGetConfig(t *testing.T) {
	assert.Panics(t, func() {
		GetConfig()
	},"config instance")
}

func TestNewConfig(t *testing.T) {

	config := NewConfig(directory)

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(directory string) {
			NewConfig(directory)
			wg.Done()
		}(directory)
	}

	wg.Wait()

	// 单例测试
	config2 := NewConfig(directory)

	assert.Equal(t, config, config2)
}

func TestConfig_Set(t *testing.T) {

	config := NewConfig(directory)

	config.Item("app").Set("abc","def")

	tests := []struct {
		File  string
		Key   string
		Value string
	}{
		{"app", "x", "x"},
		{"app", "s1.x", "s1x"},
		{"app", "s1.z.y", "s1xy"},
		{"app", "x", "x"},
	}
	fmt.Printf("out,%p",config)

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(config *Config) {
			testn := make([]struct {
				File  string
				Key   string
				Value string
			},4)
			copy(testn,tests)

			config1 := NewConfig(directory)
			fmt.Printf("in,%p\n",config1)
			for _, test := range testn {
				config1.Item(test.File).Set(test.Key, test.Value)
			}

			wg.Done()
		}(config)
	}

	wg.Wait()
}

func TestConfig_Get(t *testing.T) {
	NewConfig(directory)

	value := GetConfig().Item("app").Get("def").(int)
	assert.Equal(t,123,value)

	GetConfig().Item("app").SetDefault("abcdef","def")

	value1 := GetConfig().Item("app").Get("abcdef").(string)
	assert.Equal(t,"def",value1)
}

func TestConfig_Load(t *testing.T) {
	NewConfig(directory)
	assert.Panics(t, func() {
		GetConfig().Load("abc.yaml")
	},"open yaml")
}

func TestConfig_Item(t *testing.T) {
	NewConfig(directory)
	assert.Panics(t, func() {
		GetConfig().Item("def").Load("abc.yaml")
	},"open yaml")
}

func TestConfig_GetBool(t *testing.T) {
	NewConfig(directory).Item("app").SetDefault(t.Name() + "bool",true)
	assert.Equal(t,true,GetConfig().Item(`app`).GetBool(t.Name() + `bool`))
}

func TestConfig_GetInt(t *testing.T) {
	NewConfig(directory).Item("app").SetDefault(t.Name() , 200)
	assert.Equal(t,200,GetConfig().Item(`app`).GetInt(t.Name() ))
}

func TestConfig_GetString(t *testing.T) {
	NewConfig(directory).Item("app").SetDefault(t.Name() , "abc")
	assert.Equal(t,"abc",GetConfig().Item(`app`).GetString(t.Name() ))
}

func TestConfig_GetFloat64(t *testing.T) {
	NewConfig(directory).Item("app").SetDefault(t.Name() , 20.02)
	assert.Equal(t,20.02,GetConfig().Item(`app`).GetFloat64(t.Name() ))
}

func TestConfig_GetIntSlice(t *testing.T) {
	value := []int{1,2,3}
	NewConfig(directory).Item("app").SetDefault(t.Name() , value)
	assert.Equal(t,value,GetConfig().Item(`app`).GetIntSlice(t.Name() ))
}

func TestConfig_GetStringSlice(t *testing.T) {
	value := []string{"a","b","c"}
	NewConfig(directory).Item("app").SetDefault(t.Name() , value)
	assert.Equal(t,value,GetConfig().Item(`app`).GetStringSlice(t.Name() ))
}

func TestConfig_GetStringMap(t *testing.T) {
	value := map[string]interface{}{"a":"a","b":1,"c":2.02}
	NewConfig(directory).Item("app").SetDefault(t.Name() , value)
	assert.Equal(t,value,GetConfig().Item(`app`).GetStringMap(t.Name() ))
}

func TestConfig_GetStringMapString(t *testing.T) {
	value := map[string]string{"a":"a","b":"b","c":"c"}
	NewConfig(directory).Item("app").SetDefault(t.Name() , value)
	assert.Equal(t,value,GetConfig().Item(`app`).GetStringMapString(t.Name() ))
}

func TestConfig_GetTime(t *testing.T) {
	value := time.Now()
	NewConfig(directory).Item("app").SetDefault(t.Name() , value)
	assert.Equal(t,value,GetConfig().Item(`app`).GetTime(t.Name() ))
}

func TestConfig_GetDuration(t *testing.T) {
	value := time.Second
	NewConfig(directory).Item("app").SetDefault(t.Name() , value)
	assert.Equal(t,value,GetConfig().Item(`app`).GetDuration(t.Name() ))
}

func TestConfig_Exists(t *testing.T) {
	NewConfig(directory).Item("app").SetDefault(t.Name() , "abc")
	assert.Equal(t,true,GetConfig().Item(`app`).Exists(t.Name() ))
}