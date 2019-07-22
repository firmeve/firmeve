package config

import (
	"github.com/go-ini/ini"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// Package global variable
var (
	config *Config
	mutex  sync.Mutex
	once   sync.Once
)

// Configurator interface
type Configurator interface {
	Get(keys string) (interface{}, error)
	Set(key string, value string) error
	All() map[string]*viper.Viper
}

// Config struct
type Config struct {
	directory string
	configs   map[string]*viper.Viper
	delimiter string
	extension string
}

// Create a new config instance
func NewConfig(directory string) (*Config, error) {
	// singleton
	if config != nil {
		return config, nil
	}

	var err error
	// 单例加锁
	//mu.Lock()
	//defer mu.Unlock()

	//直接使用once，其实就是调用mu.Lock和Unlock
	once.Do(func() {
		directory, err := absPath(directory)
		if err != nil {
			return
		}

		config = &Config{
			directory: directory,
			delimiter: `.`, extension: `.yaml`,
			configs: make(map[string]*viper.Viper),
		}
		// loadAll
		err = config.loadAll()
	})

	return config, err
}

// Get the absolute path of the specified directory
func absPath(directory string) (string, error) {
	return filepath.Abs(directory)
}

// Get the specified key configuration
func (this *Config) Get(keys string) (interface{}, error) {

	keySlices := parseKey(keys, this.delimiter)

	length := len(keySlices)

	var cfg *ini.File
	if _, ok := this.configs[keySlices[0]]; !ok {
		return nil, &FormatError{message: `index error`}
	} else {
		cfg = this.configs[keySlices[0]]
	}

	if length == 1 {
		return cfg, nil
	} else if length == 2 {
		if !cfg.Section(ini.DefaultSection).HasKey(keySlices[1]) {
			return nil, &FormatError{message: `value not found`}
		}
		return cfg.Section(ini.DefaultSection).GetKey(keySlices[1])
	} else {
		// 取得可能是section的slice拼接
		sectionName := strings.Join(keySlices[1:length-1], this.delimiter)
		section, err := cfg.GetSection(sectionName)
		// 如果section不存在拼接全部的
		if err != nil {
			sectionName = sectionName + `.` + keySlices[length-1:][0]
			section, err = cfg.GetSection(sectionName)
			if err != nil {
				return nil, err
			}

			return section, nil
		}

		// 否则最后一位肯定是key
		key := keySlices[length-1:][0]
		if section.HasKey(key) {
			return section.GetKey(key)
		}

		return nil, &FormatError{message: `value not found`}
	}
}

// Set configuration value
func (this *Config) Set(keys string, value string) error {
	// map加锁
	mutex.Lock()
	defer mutex.Unlock()

	keySlices := parseKey(keys, this.delimiter)

	length := len(keySlices)

	if length == 1 {
		return &FormatError{message: "incorrect parameter format"}
	}

	var err error

	if _, ok := this.configs[keySlices[0]]; !ok {
		this.configs[keySlices[0]], err = loadConf(this.fullPath(keySlices[0]))
		if err != nil {
			return err
		}
	}

	if length == 2 {
		_, err = this.configs[keySlices[0]].Section(ini.DefaultSection).NewKey(keySlices[1], value)
	} else {
		_, err = this.configs[keySlices[0]].Section(strings.Join(keySlices[1:length-1], this.delimiter)).NewKey(keySlices[length-1], value)
	}
	if err != nil {
		return err
	}

	file, err := os.OpenFile(this.fullPath(keySlices[0]), os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = this.configs[keySlices[0]].WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}

// Get all configurations
func (this *Config) All() map[string]*viper.Viper {
	return this.configs
}

// Load all configuration files at once
func (this *Config) loadAll() error {

	err := filepath.Walk(this.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		filename := filepath.Base(path)

		if strings.Contains(filename, this.extension) && !info.IsDir() {
			fileKey := strings.Replace(filename, this.extension, "", 1)

			cfg, err := loadConf(path)
			if err != nil {
				return err
			}

			this.configs[fileKey] = cfg
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Get the current full path by filename (without extension)
func (this *Config) fullPath(filename string) string {
	return path.Clean(this.directory + "/" + filename + this.extension)
}

// Load configuration file
func loadConf(filename string) (*viper.Viper, error) {

	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
			defer file.Close()
		}
	}

	conf := viper.New()
	conf.SetConfigFile(filename)
	err = conf.ReadInConfig()
	return conf, err
}

// Parsing data key
func parseKey(key string, delimiter string) []string {
	return strings.Split(key, delimiter)
}

// ------------------------- config error --------------------------------

//config format error
type FormatError struct {
	message string
	//err error
}

func (this *FormatError) Error() string {
	return this.message
}
