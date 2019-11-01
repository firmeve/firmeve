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
	mutex  sync.Mutex
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
	Item(item string) Configurator
	Load(file string)
}

// Config struct
type config struct {
	directory string
	items     map[string]*viper.Viper
	current   *viper.Viper
	delimiter string
	extension string
}

// Create a new config instance
func New(directory string) Configurator {
	directory, err := filepath.Abs(directory)
	if err != nil {
		panic(err)
	}

	config := &config{
		directory: directory,
		current:   nil,
		delimiter: `.`,
		extension: `.yaml`,
		items:     make(map[string]*viper.Viper),
	}
	config.loadAll()

	return config
}

//---------------------- config ------------------------

// Get the current file node
func (c *config) Item(item string) Configurator {
	mutex.Lock()
	defer mutex.Unlock()
	if itemConfig, ok := c.items[item]; ok {
		c.current = itemConfig
		return c
	}

	panic(fmt.Errorf(`the config %s not exists`, item))
}

// Get the value of the specified key
func (c *config) Get(key string) interface{} {
	return c.current.Get(key)
}

// Get the bool value of the specified key
func (c *config) GetBool(key string) bool {
	return c.current.GetBool(key)
}

// Get the float value of the specified key
func (c *config) GetFloat(key string) float64 {
	return c.current.GetFloat64(key)
}

// Get the int value of the specified key
func (c *config) GetInt(key string) int {
	return c.current.GetInt(key)
}

// Get the int slice value of the specified key
func (c *config) GetIntSlice(key string) []int {
	return c.current.GetIntSlice(key)
}

// Get the string value of the specified key
func (c *config) GetString(key string) string {
	return c.current.GetString(key)
}

// Get the string map value of the specified key
func (c *config) GetStringMap(key string) map[string]interface{} {
	return c.current.GetStringMap(key)
}

// Get the string map string value of the specified key
func (c *config) GetStringMapString(key string) map[string]string {
	return c.current.GetStringMapString(key)
}

// Get the string slice value of the specified key
func (c *config) GetStringSlice(key string) []string {
	return c.current.GetStringSlice(key)
}

// Get the time type value of the specified key
func (c *config) GetTime(key string) time.Time {
	return c.current.GetTime(key)
}

// Get the time duration value of the specified key
func (c *config) GetDuration(key string) time.Duration {
	return c.current.GetDuration(key)
}

// Determine if the specified key exists
func (c *config) Exists(key string) bool {
	return c.current.IsSet(key)
}

// Set configuration value
func (c *config) Set(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	c.current.Set(key, value)
}

// Set item default value if not exists
func (c *config) SetDefault(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	c.current.SetDefault(key, value)
}

// Load all configuration files at once
func (c *config) loadAll() {
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

func (c *config) Load(file string) {
	mutex.Lock()
	defer mutex.Unlock()

	conf := viper.New()
	conf.SetConfigFile(file)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}

	c.items[strings.Replace(filepath.Base(file), c.extension, "", 1)] = conf
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