package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Package global variable
var (
	config *Config
	mutex  sync.Mutex
	once   sync.Once
)

// Configurator interface
type Configurator interface {
	Get(key string) interface{}
	Set(key string, value interface{})
	Exists(key string) bool
}

// Config struct
type Config struct {
	directory string
	items     map[string]*viper.Viper
	current   *viper.Viper
	delimiter string
	extension string
}

// Create a new config instance
func NewConfig(directory string) *Config {
	// singleton
	if config != nil {
		return config
	}

	//直接使用once，其实就是调用mu.Lock和Unlock
	once.Do(func() {
		directory, err := filepath.Abs(directory)
		if err != nil {
			panic(err.Error())
		}

		config = &Config{
			directory: directory,
			current:   nil,
			delimiter: `.`,
			extension: `.yaml`,
			items:     make(map[string]*viper.Viper),
		}

		// loadAll
		err = config.loadAll()
		if err != nil {
			panic(err.Error())
		}
	})

	return config
}

func GetConfig() *Config {
	// singleton
	if config != nil {
		return config
	}

	panic(`config instance not exists`)
}

//---------------------- config ------------------------

// Get the current file node
func (c *Config) Item(item string) *Config {
	mutex.Lock()
	defer mutex.Unlock()
	if itemConfig, ok := c.items[item]; ok {
		c.current = itemConfig
		return c
	}

	panic(fmt.Sprintf(`the config %s not exists`, item))
}

// Get the value of the specified key
func (c *Config) Get(key string) interface{} {
	return c.current.Get(key)
}

// Get the bool value of the specified key
func (c *Config) GetBool(key string) bool {
	return c.current.GetBool(key)
}

// Get the float value of the specified key
func (c *Config) GetFloat64(key string) float64 {
	return c.current.GetFloat64(key)
}

// Get the int value of the specified key
func (c *Config) GetInt(key string) int {
	return c.current.GetInt(key)
}

// Get the int slice value of the specified key
func (c *Config) GetIntSlice(key string) []int {
	return c.current.GetIntSlice(key)
}

// Get the string value of the specified key
func (c *Config) GetString(key string) string {
	return c.current.GetString(key)
}

// Get the string map value of the specified key
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.current.GetStringMap(key)
}

// Get the string map string value of the specified key
func (c *Config) GetStringMapString(key string) map[string]string {
	return c.current.GetStringMapString(key)
}

// Get the string slice value of the specified key
func (c *Config) GetStringSlice(key string) []string {
	return c.current.GetStringSlice(key)
}

// Get the time type value of the specified key
func (c *Config) GetTime(key string) time.Time {
	return c.current.GetTime(key)
}

// Get the time duration value of the specified key
func (c *Config) GetDuration(key string) time.Duration {
	return c.current.GetDuration(key)
}

// Determine if the specified key exists
func (c *Config) Exists(key string) bool {
	return c.current.IsSet(key)
}

// Set configuration value
func (c *Config) Set(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	c.current.Set(key, value)
}

// Set item default value if not exists
func (c *Config) SetDefault(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	c.current.SetDefault(key, value)
}

// Load all configuration files at once
func (c *Config) loadAll() error {
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
		return err
	}

	return nil
}

func (c *Config) Load(file string) {
	mutex.Lock()
	defer mutex.Unlock()

	conf := viper.New()
	conf.SetConfigFile(file)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}

	c.items[strings.Replace(filepath.Base(file), c.extension, "", 1)] = conf
}
