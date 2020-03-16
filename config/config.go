package config

import (
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/filesystem"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		path      string
		items     map[string]*item
		delimiter string
		extension string
	}

	item struct {
		config *viper.Viper
	}
)

// Package global variable
var (
	mutex sync.Mutex
)

// Create a new config instance
func New(path string) *Config {
	config := &Config{
		path:      path,
		delimiter: `.`,
		extension: `.yaml`,
		items:     make(map[string]*item),
	}

	if filesystem.IsFile(path) {
		config.load(path)
	} else {
		directory, err := filepath.Abs(path)
		if err != nil {
			panic(err)
		}
		config.loadFiles(directory)
	}

	return config
}

//---------------------- config ------------------------

// Get the current file node
func (c *Config) Item(item string) contract.Configuration {
	mutex.Lock()
	defer mutex.Unlock()
	if itemConfig, ok := c.items[item]; ok {
		return itemConfig
	}

	panic(fmt.Errorf(`the config %s not exists`, item))
}

// Load all configuration files at once
func (c *Config) loadFiles(directory string) {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(filepath.Base(path), c.extension) && !info.IsDir() {
			c.loadFile(path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (c *Config) loadFile(file string) {
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

func (c *Config) load(path string) {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	for key := range conf.AllSettings() {
		c.items[key] = &item{
			config: conf.Sub(key),
		}
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
