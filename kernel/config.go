package kernel

import (
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/filesystem"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Package global variable
var (
	mutex sync.Mutex
)

type (
	Config struct {
		driver *viper.Viper
		path   string
	}
)

// Create a new config instance
func NewConfig(path string) *Config {
	viper := viper.New()
	viper.SetConfigType(`yaml`)

	if filesystem.IsFile(path) {
		viper.SetConfigFile(path)
	} else {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{
		path:   path,
		driver: viper,
	}
}

func (c *Config) Bind(name string, object interface{}) error {
	if err := c.driver.Sub(name).Unmarshal(object); err != nil {
		return Error(err)
	}

	return nil
}

// Get the value of the specified key
func (c *Config) Get(key string) interface{} {
	return c.driver.Get(key)
}

// Get the bool value of the specified key
func (c *Config) GetBool(key string) bool {
	return c.driver.GetBool(key)
}

// Get the float value of the specified key
func (c *Config) GetFloat(key string) float64 {
	return c.driver.GetFloat64(key)
}

// Get the int value of the specified key
func (c *Config) GetInt(key string) int {
	return c.driver.GetInt(key)
}

// Get the int slice value of the specified key
func (c *Config) GetIntSlice(key string) []int {
	return c.driver.GetIntSlice(key)
}

// Get the string value of the specified key
func (c *Config) GetString(key string) string {
	return c.driver.GetString(key)
}

// Get the string map value of the specified key
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.driver.GetStringMap(key)
}

// Get the string map string value of the specified key
func (c *Config) GetStringMapString(key string) map[string]string {
	return c.driver.GetStringMapString(key)
}

// Get the string slice value of the specified key
func (c *Config) GetStringSlice(key string) []string {
	return c.driver.GetStringSlice(key)
}

// Get the time type value of the specified key
func (c *Config) GetTime(key string) time.Time {
	return c.driver.GetTime(key)
}

// Get the time duration value of the specified key
func (c *Config) GetDuration(key string) time.Duration {
	return c.driver.GetDuration(key)
}

// Determine if the specified key exists
func (c *Config) Exists(key string) bool {
	return c.driver.IsSet(key)
}

// Set configuration value
func (c *Config) Set(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	c.driver.Set(key, value)
}

// Set item default value if not exists
func (c *Config) SetDefault(key string, value interface{}) {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	c.driver.SetDefault(key, value)
}

func (c *Config) Clone() contract.Configuration {
	config := &Config{}
	*config = *c
	*config.driver = *c.driver
	return config
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
