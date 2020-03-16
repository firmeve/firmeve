package config

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
)

var (
	directory  = path.RunRelative("../testdata/config")
	configPath = path.RunRelative("../testdata/config/config.yaml")
	wg         sync.WaitGroup
)

func TestConfig_File(t *testing.T) {
	config := New(configPath)

	assert.NotEqual(t, fmt.Sprintf("%p", config.Item(`database`)), fmt.Sprintf("%p", config.Item(`cache`)))
	assert.Equal(t, `mysql`, config.Item(`database`).Get(`default`).(string))
	assert.Equal(t, `firmeve_cache`, config.Item(`cache`).Get(`prefix`).(string))
}

func TestConfig_Set(t *testing.T) {
	config := New(directory)

	config.Item("app").Set("abc", "def")
	fmt.Printf("%#v", config)
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
	fmt.Printf("out,%p", config)

	wg.Add(2)
	for i := 0; i < 2; i++ {
		go func(config contract.Configuration) {
			testn := make([]struct {
				File  string
				Key   string
				Value string
			}, 4)
			copy(testn, tests)

			config1 := New(directory)
			fmt.Printf("in,%p\n", config1)
			for _, test := range testn {
				config1.Item(test.File).Set(test.Key, test.Value)
			}

			wg.Done()
		}(config.Item("app"))
	}

	wg.Wait()
}

func TestConfig_Get(t *testing.T) {
	instance := New(directory)

	value := instance.Item("app").Get("def").(int)
	assert.Equal(t, 123, value)

	instance.Item("app").SetDefault("abcdef", "def")

	value1 := instance.Item("app").Get("abcdef").(string)
	assert.Equal(t, "def", value1)
}

func TestConfig_Load(t *testing.T) {
	assert.Panics(t, func() {
		New(directory).loadFile("abc.yaml")
	}, "open yaml")
}

func TestConfig_Item(t *testing.T) {
	assert.Panics(t, func() {
		New(directory).loadFile("abc.yaml")
	}, "open yaml")
}

func TestConfig_GetBool(t *testing.T) {
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name()+"bool", true)
	assert.Equal(t, true, instance.Item(`app`).GetBool(t.Name()+`bool`))
}

func TestConfig_GetInt(t *testing.T) {
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), 200)
	assert.Equal(t, 200, instance.Item(`app`).GetInt(t.Name()))
}

func TestConfig_GetString(t *testing.T) {
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), "abc")
	assert.Equal(t, "abc", instance.Item(`app`).GetString(t.Name()))
}

func TestConfig_GetFloat64(t *testing.T) {
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), 20.02)
	assert.Equal(t, 20.02, instance.Item(`app`).GetFloat(t.Name()))
}

func TestConfig_GetIntSlice(t *testing.T) {
	value := []int{1, 2, 3}
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, instance.Item(`app`).GetIntSlice(t.Name()))
}

func TestConfig_GetStringSlice(t *testing.T) {
	value := []string{"a", "b", "c"}
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, instance.Item(`app`).GetStringSlice(t.Name()))
}

func TestConfig_GetStringMap(t *testing.T) {
	value := map[string]interface{}{"a": "a", "b": 1, "c": 2.02}
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, instance.Item(`app`).GetStringMap(t.Name()))
}

func TestConfig_GetStringMapString(t *testing.T) {
	value := map[string]string{"a": "a", "b": "b", "c": "c"}
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, instance.Item(`app`).GetStringMapString(t.Name()))
}

func TestConfig_GetTime(t *testing.T) {
	value := time.Now()
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, instance.Item(`app`).GetTime(t.Name()))
}

func TestConfig_GetDuration(t *testing.T) {
	value := time.Second
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, instance.Item(`app`).GetDuration(t.Name()))
}

func TestConfig_Exists(t *testing.T) {
	instance := New(directory)
	instance.Item("app").SetDefault(t.Name(), "abc")
	assert.Equal(t, true, instance.Item(`app`).Exists(t.Name()))
}

//func TestSetEnv(t *testing.T) {
//	SetEnv("s", "s")
//	assert.Equal(t, "s", os.Getenv("S"))
//}

func TestGetEnv(t *testing.T) {
	// key must capital
	os.Setenv("AAA", "a")
	// 当使用SetEnv函数时会自动转换为大写
	SetEnv("abc", "abc")

	assert.Equal(t, "a", GetEnv("AAA"))
	// 当使用GetEnv函数时会自动转换为大写
	assert.Equal(t, "abc", GetEnv("abc"))

	os.Setenv("lower", "lower")
	assert.Equal(t, "", GetEnv("lower"))
}
