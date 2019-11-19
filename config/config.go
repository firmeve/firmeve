package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Package global variable
var (
	mutex sync.Mutex
)

// Configurator interface
type Configurator interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetFloat(key string) float64
	GetInt(key string) int
	GetIntSlice(key string) []int
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	Exists(key string) bool
	Set(key string, value interface{})
	SetDefault(key string, value interface{})
	//Item(item string) Configurator
	//Load(file string)
}

// Config struct
type Config struct {
	directory string
	items     map[string]*item
	delimiter string
	extension string
}

type item struct {
	config *viper.Viper
}

//
//var (
//	configVar *config
//	once      sync.Once
//)

//func Config() *config {
//	if configVar != nil{
//		return configVar
//	}
//
//	once.Do(func() {
//		config := New()
//	})
//}

// Create a new config instance
func New(directory string) *Config {
	directory, err := filepath.Abs(directory)
	if err != nil {
		panic(err)
	}

	config := &Config{
		directory: directory,
		delimiter: `.`,
		extension: `.yaml`,
		items:     make(map[string]*item),
	}
	config.loadAll()

	return config
}

//---------------------- config ------------------------

// Get the current file node
func (c *Config) Item(item string) Configurator {
	mutex.Lock()
	defer mutex.Unlock()
	if itemConfig, ok := c.items[item]; ok {
		return itemConfig
	}

	panic(fmt.Errorf(`the config %s not exists`, item))
}

// Load all configuration files at once
func (c *Config) loadAll() {
	err := filepath.Walk(c.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(filepath.Base(path), c.extension) && !info.IsDir() {
			c.Load(path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (c *Config) Load(file string) {
	mutex.Lock()
	defer mutex.Unlock()

	conf := viper.New()
	conf.SetConfigFile(file)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	c.items[strings.Replace(filepath.Base(file), c.extension, "", 1)] = &item{
		config: conf,
	}
}

// Get the value of the specified key
func (i *item) Get(key string) interface{} {
	return i.config.Get(key)
}

// Get the bool value of the specified key
func (i *item) GetBool(key string) bool {
	return i.config.GetBool(key)
}

// Get the float value of the specified key
func (i *item) GetFloat(key string) float64 {
	return i.config.GetFloat64(key)
}

// Get the int value of the specified key
func (i *item) GetInt(key string) int {
	return i.config.GetInt(key)
}

// Get the int slice value of the specified key
func (i *item) GetIntSlice(key string) []int {
	return i.config.GetIntSlice(key)
}

// Get the string value of the specified key
func (i *item) GetString(key string) string {
	return i.config.GetString(key)
}

// Get the string map value of the specified key
func (i *item) GetStringMap(key string) map[string]interface{} {
	return i.config.GetStringMap(key)
}

// Get the string map string value of the specified key
func (i *item) GetStringMapString(key string) map[string]string {
	return i.config.GetStringMapString(key)
}

// Get the string slice value of the specified key
func (i *item) GetStringSlice(key string) []string {
	return i.config.GetStringSlice(key)
}

// Get the time type value of the specified key
func (i *item) GetTime(key string) time.Time {
	return i.config.GetTime(key)
}

// Get the time duration value of the specified key
func (i *item) GetDuration(key string) time.Duration {
	return i.config.GetDuration(key)
}

// Determine if the specified key exists
func (i *item) Exists(key string) bool {
	return i.config.IsSet(key)
}

// Set configuration value
func (i *item) Set(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	i.config.Set(key, value)
}

// Set item default value if not exists
func (i *item) SetDefault(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	i.config.SetDefault(key, value)
}

// -------------------------------- Env ----------------------------------------

// Set env
func SetEnv(key, value string) {
	err := os.Setenv(strings.ToUpper(key), value)
	if err != nil {
		panic(err)
	}
}

// Get all env
func GetEnv(key string) string {
	// 自动查找所有环境变量
	viper.AutomaticEnv()

	return viper.GetString(strings.ToUpper(key))
}
