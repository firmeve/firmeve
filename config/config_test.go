package config

import (
	"fmt"
	firmeve2 "github.com/firmeve/firmeve"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	directory = "../testdata/config"
	wg        sync.WaitGroup
)

func TestNew(t *testing.T) {

	config := New(directory)

	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(directory string) {
			New(directory)
			wg.Done()
		}(directory)
	}

	wg.Wait()

	// 单例测试
	config2 := New(directory)

	assert.Equal(t, config, config2)
	config = nil
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
		go func(config *Config) {
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
		}(config)
	}

	wg.Wait()
}

func TestConfig_Get(t *testing.T) {
	New(directory)

	value := Get().Item("app").Get("def").(int)
	assert.Equal(t, 123, value)

	Get().Item("app").SetDefault("abcdef", "def")

	value1 := Get().Item("app").Get("abcdef").(string)
	assert.Equal(t, "def", value1)
}

func TestConfig_Load(t *testing.T) {
	New(directory)
	assert.Panics(t, func() {
		Get().Load("abc.yaml")
	}, "open yaml")
}

func TestConfig_Item(t *testing.T) {
	New(directory)
	assert.Panics(t, func() {
		Get().Item("def").Load("abc.yaml")
	}, "open yaml")
}

func TestConfig_GetBool(t *testing.T) {
	New(directory).Item("app").SetDefault(t.Name()+"bool", true)
	assert.Equal(t, true, Get().Item(`app`).GetBool(t.Name()+`bool`))
}

func TestConfig_GetInt(t *testing.T) {
	New(directory).Item("app").SetDefault(t.Name(), 200)
	assert.Equal(t, 200, Get().Item(`app`).GetInt(t.Name()))
}

func TestConfig_GetString(t *testing.T) {
	New(directory).Item("app").SetDefault(t.Name(), "abc")
	assert.Equal(t, "abc", Get().Item(`app`).GetString(t.Name()))
}

func TestConfig_GetFloat64(t *testing.T) {
	New(directory).Item("app").SetDefault(t.Name(), 20.02)
	assert.Equal(t, 20.02, Get().Item(`app`).GetFloat64(t.Name()))
}

func TestConfig_GetIntSlice(t *testing.T) {
	value := []int{1, 2, 3}
	New(directory).Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, Get().Item(`app`).GetIntSlice(t.Name()))
}

func TestConfig_GetStringSlice(t *testing.T) {
	value := []string{"a", "b", "c"}
	New(directory).Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, Get().Item(`app`).GetStringSlice(t.Name()))
}

func TestConfig_GetStringMap(t *testing.T) {
	value := map[string]interface{}{"a": "a", "b": 1, "c": 2.02}
	New(directory).Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, Get().Item(`app`).GetStringMap(t.Name()))
}

func TestConfig_GetStringMapString(t *testing.T) {
	value := map[string]string{"a": "a", "b": "b", "c": "c"}
	New(directory).Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, Get().Item(`app`).GetStringMapString(t.Name()))
}

func TestConfig_GetTime(t *testing.T) {
	value := time.Now()
	New(directory).Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, Get().Item(`app`).GetTime(t.Name()))
}

func TestConfig_GetDuration(t *testing.T) {
	value := time.Second
	New(directory).Item("app").SetDefault(t.Name(), value)
	assert.Equal(t, value, Get().Item(`app`).GetDuration(t.Name()))
}

func TestConfig_Exists(t *testing.T) {
	New(directory).Item("app").SetDefault(t.Name(), "abc")
	assert.Equal(t, true, Get().Item(`app`).Exists(t.Name()))
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

	os.Setenv("lower","lower")
	assert.Equal(t,"",GetEnv("lower"))
}

func TestProvider_Register(t *testing.T) {
	firmeve := firmeve2.NewFirmeve()
	firmeve.Boot()
	assert.Equal(t, true, firmeve.HasProvider("config"))
	assert.Equal(t,true,firmeve.Has(`config`))
}