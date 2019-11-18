package database

import (
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type dbConnection map[string]*gorm.DB

type ConfigDriver map[string]interface{}

type Config struct {
	Current string
	Drivers ConfigDriver
}

type DB struct {
	config      *Config
	db          *gorm.DB
	connections dbConnection
}

func New(config *Config) *DB {
	return &DB{
		config:      config,
		connections: make(dbConnection, 0),
	}
}

func Default() *DB {
	return &DB{
		config: &Config{
			Current: `mysql`,
		},
		connections: make(dbConnection, 0),
	}
}

func (d *DB) ConnectionDefault() *gorm.DB {
	return d.Connection(d.config.Current)
}

func (d *DB) Connection(driver string) *gorm.DB {
	if connection, ok := d.connections[driver]; ok {
		return connection
	}

	config := d.config.GetString(strings.Join([]string{driver, `addr`}, `.`))
	db, err := gorm.Open(driver, config)

	if err != nil {
		panic(err)
	}

	d.connections[driver] = db

	return db
}

func (d *DB) CloseDefault() {
	d.Close(d.config.GetString(`default`))
}

func (d *DB) Close(driver string) {
	if connection, ok := d.connections[driver]; ok {
		err := connection.Close()
		if err != nil {
			panic(err)
		}
		delete(d.connections, driver)
	}
}
