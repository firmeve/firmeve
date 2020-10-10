package kernel

import (
	"github.com/firmeve/firmeve/support/path"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

var (
	configPath = path.RunRelative("../testdata/config/config.testing.yaml")
	wg         sync.WaitGroup
)

func TestConfig_Get(t *testing.T) {
	instance := NewConfig(configPath)

	value := instance.Get("framework.lang").(string)
	assert.Equal(t, `zh-CN`, value)

	instance.SetDefault("app.abcdef", "def")

	value1 := instance.Get("app.abcdef").(string)
	assert.Equal(t, "def", value1)
}

//
//func TestConfig_GetBool(t *testing.T) {
//	instance := NewConfig(configPath)
//	instance.Item("app").SetDefault(t.Name()+"bool", true)
//	assert.Equal(t, true, instance.Item(`app`).GetBool(t.Name()+`bool`))
//}
//
//func TestConfig_GetInt(t *testing.T) {
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), 200)
//	assert.Equal(t, 200, instance.Item(`app`).GetInt(t.Name()))
//}
//
//func TestConfig_GetString(t *testing.T) {
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), "abc")
//	assert.Equal(t, "abc", instance.Item(`app`).GetString(t.Name()))
//}
//
//func TestConfig_GetFloat64(t *testing.T) {
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), 20.02)
//	assert.Equal(t, 20.02, instance.Item(`app`).GetFloat(t.Name()))
//}
//
//func TestConfig_GetIntSlice(t *testing.T) {
//	value := []int{1, 2, 3}
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), value)
//	assert.Equal(t, value, instance.Item(`app`).GetIntSlice(t.Name()))
//}
//
//func TestConfig_GetStringSlice(t *testing.T) {
//	value := []string{"a", "b", "c"}
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), value)
//	assert.Equal(t, value, instance.Item(`app`).GetStringSlice(t.Name()))
//}
//
//func TestConfig_GetStringMap(t *testing.T) {
//	value := map[string]interface{}{"a": "a", "b": 1, "c": 2.02}
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), value)
//	assert.Equal(t, value, instance.Item(`app`).GetStringMap(t.Name()))
//}
//
//func TestConfig_GetStringMapString(t *testing.T) {
//	value := map[string]string{"a": "a", "b": "b", "c": "c"}
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), value)
//	assert.Equal(t, value, instance.Item(`app`).GetStringMapString(t.Name()))
//}
//
//func TestConfig_GetTime(t *testing.T) {
//	value := time.Now()
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), value)
//	assert.Equal(t, value, instance.Item(`app`).GetTime(t.Name()))
//}
//
//func TestConfig_GetDuration(t *testing.T) {
//	value := time.Second
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), value)
//	assert.Equal(t, value, instance.Item(`app`).GetDuration(t.Name()))
//}
//
//func TestConfig_Exists(t *testing.T) {
//	instance := New(directory)
//	instance.Item("app").SetDefault(t.Name(), "abc")
//	assert.Equal(t, true, instance.Item(`app`).Exists(t.Name()))
//}

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
